package main

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/sethvargo/go-retry"

	"course-project/internal/config"
	"course-project/internal/database"
	"course-project/internal/handlers"
	"course-project/internal/middleware"
	"course-project/internal/repository"
	"course-project/internal/service"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var pool *database.PoolWrapper
	err = retry.Do(ctx, retry.WithMaxDuration(30*time.Second, retry.NewConstant(1*time.Second)), func(ctx context.Context) error {
		pool, err = database.NewPool(ctx, cfg)
		return err
	})
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("connected to database")

	if err := runMigrations(pool); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	log.Println("migrations applied")

	repo := repository.New(pool.DB())
	svc := service.New(repo, cfg)
	h := handlers.New(svc)

	app := fiber.New(fiber.Config{
		BodyLimit: 4 * 1024 * 1024,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(func(c fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization")
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	})

	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/login", h.Login)

	protected := api.Group("", middleware.AuthMiddleware(svc))

	schedule := protected.Group("/schedule")
	schedule.Get("", h.GetSchedule)
	schedule.Post("/entries", middleware.RoleMiddleware("teacher", "admin"), h.CreateScheduleEntry)
	schedule.Put("/entries/:id", middleware.RoleMiddleware("teacher", "admin"), h.UpdateScheduleEntry)
	schedule.Delete("/entries/:id", middleware.RoleMiddleware("teacher", "admin"), h.DeleteScheduleEntry)

	teachers := protected.Group("/teachers")
	teachers.Get("", h.GetTeachers)
	teachers.Post("", middleware.RoleMiddleware("admin"), h.CreateTeacher)
	teachers.Put("/:id", middleware.RoleMiddleware("admin"), h.UpdateTeacher)
	teachers.Delete("/:id", middleware.RoleMiddleware("admin"), h.DeleteTeacher)

	subjects := protected.Group("/subjects")
	subjects.Get("", h.GetSubjects)
	subjects.Post("", middleware.RoleMiddleware("admin"), h.CreateSubject)
	subjects.Put("/:id", middleware.RoleMiddleware("admin"), h.UpdateSubject)
	subjects.Delete("/:id", middleware.RoleMiddleware("admin"), h.DeleteSubject)

	classrooms := protected.Group("/classrooms")
	classrooms.Get("", h.GetClassrooms)
	classrooms.Post("", middleware.RoleMiddleware("admin"), h.CreateClassroom)
	classrooms.Put("/:id", middleware.RoleMiddleware("admin"), h.UpdateClassroom)
	classrooms.Delete("/:id", middleware.RoleMiddleware("admin"), h.DeleteClassroom)

	groups := protected.Group("/groups")
	groups.Get("", h.GetGroups)
	groups.Post("", middleware.RoleMiddleware("admin"), h.CreateGroup)
	groups.Put("/:id", middleware.RoleMiddleware("admin"), h.UpdateGroup)
	groups.Delete("/:id", middleware.RoleMiddleware("admin"), h.DeleteGroup)

	go func() {
		if err := app.Listen(":" + cfg.ServerPort); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
}

func runMigrations(pool *database.PoolWrapper) error {
	db, err := sql.Open("pgx", pool.DSN())
	if err != nil {
		return err
	}
	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.Up(db, "migrations")
}
