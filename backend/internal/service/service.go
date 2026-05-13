package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"course-project/internal/config"
	"course-project/internal/models"
	"course-project/internal/repository"
)

type Service struct {
	repo   *repository.Repository
	config *config.Config
}

func New(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{repo: repo, config: cfg}
}

type Claims struct {
	UserID int    `json:"user_id"`
	Login  string `json:"login"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (s *Service) Login(ctx context.Context, login, password string) (*models.LoginResponse, error) {
	teacher, err := s.repo.CheckTeacherPassword(ctx, login, password)
	if err != nil {
		return nil, err
	}
	if teacher == nil {
		return nil, nil
	}

	claims := &Claims{
		UserID: teacher.ID,
		Login:  teacher.Login,
		Role:   teacher.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{Token: tokenString, Role: teacher.Role}, nil
}

func (s *Service) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}

func (s *Service) GetTeachers(ctx context.Context) ([]models.Teacher, error) {
	return s.repo.GetTeachers(ctx)
}

func (s *Service) CreateTeacher(ctx context.Context, t *models.Teacher) error {
	return s.repo.CreateTeacher(ctx, t)
}

func (s *Service) GetTeacherByID(ctx context.Context, id int) (*models.Teacher, error) {
	return s.repo.GetTeacherByID(ctx, id)
}

func (s *Service) UpdateTeacher(ctx context.Context, t *models.Teacher) error {
	return s.repo.UpdateTeacher(ctx, t)
}

func (s *Service) DeleteTeacher(ctx context.Context, id int) error {
	return s.repo.DeleteTeacher(ctx, id)
}

func (s *Service) GetSubjects(ctx context.Context) ([]models.Subject, error) {
	return s.repo.GetSubjects(ctx)
}

func (s *Service) CreateSubject(ctx context.Context, sbj *models.Subject) error {
	return s.repo.CreateSubject(ctx, sbj)
}

func (s *Service) GetSubjectByID(ctx context.Context, id int) (*models.Subject, error) {
	return s.repo.GetSubjectByID(ctx, id)
}

func (s *Service) UpdateSubject(ctx context.Context, sbj *models.Subject) error {
	return s.repo.UpdateSubject(ctx, sbj)
}

func (s *Service) DeleteSubject(ctx context.Context, id int) error {
	return s.repo.DeleteSubject(ctx, id)
}

func (s *Service) GetClassrooms(ctx context.Context) ([]models.Classroom, error) {
	return s.repo.GetClassrooms(ctx)
}

func (s *Service) CreateClassroom(ctx context.Context, c *models.Classroom) error {
	return s.repo.CreateClassroom(ctx, c)
}

func (s *Service) GetClassroomByID(ctx context.Context, id int) (*models.Classroom, error) {
	return s.repo.GetClassroomByID(ctx, id)
}

func (s *Service) UpdateClassroom(ctx context.Context, c *models.Classroom) error {
	return s.repo.UpdateClassroom(ctx, c)
}

func (s *Service) DeleteClassroom(ctx context.Context, id int) error {
	return s.repo.DeleteClassroom(ctx, id)
}

func (s *Service) GetGroups(ctx context.Context) ([]models.Group, error) {
	return s.repo.GetGroups(ctx)
}

func (s *Service) CreateGroup(ctx context.Context, g *models.Group) error {
	return s.repo.CreateGroup(ctx, g)
}

func (s *Service) GetGroupByID(ctx context.Context, id int) (*models.Group, error) {
	return s.repo.GetGroupByID(ctx, id)
}

func (s *Service) UpdateGroup(ctx context.Context, g *models.Group) error {
	return s.repo.UpdateGroup(ctx, g)
}

func (s *Service) DeleteGroup(ctx context.Context, id int) error {
	return s.repo.DeleteGroup(ctx, id)
}

func (s *Service) GetSchedule(ctx context.Context, groupID, teacherID int, date string) ([]models.Schedule, error) {
	schedules, err := s.repo.GetScheduleByFilters(ctx, groupID, teacherID, date)
	if err != nil {
		return nil, err
	}

	for i := range schedules {
		entries, err := s.repo.GetScheduleEntries(ctx, schedules[i].ID)
		if err != nil {
			return nil, err
		}
		schedules[i].Pairs = entries
	}

	return schedules, nil
}

func (s *Service) CreateScheduleEntry(ctx context.Context, req *models.CreateEntryRequest) (*models.ScheduleEntry, error) {
	return s.repo.CreateScheduleEntry(ctx, req)
}

func (s *Service) UpdateScheduleEntry(ctx context.Context, id int, req *models.CreateEntryRequest) error {
	return s.repo.UpdateScheduleEntry(ctx, id, req)
}

func (s *Service) DeleteScheduleEntry(ctx context.Context, id int) error {
	return s.repo.DeleteScheduleEntry(ctx, id)
}

func (s *Service) GetScheduleEntryByID(ctx context.Context, id int) (*models.ScheduleEntry, error) {
	return s.repo.GetScheduleEntryByID(ctx, id)
}

func (s *Service) Repo() *repository.Repository {
	return s.repo
}
