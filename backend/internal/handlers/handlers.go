package handlers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"

	"course-project/internal/models"
	"course-project/internal/service"
	"course-project/internal/validators"
)

type Handler struct {
	svc       *service.Service
	validator *validators.ConflictValidator
}

func New(svc *service.Service) *Handler {
	return &Handler{
		svc:       svc,
		validator: validators.NewConflictValidator(svc.Repo()),
	}
}

func (h *Handler) Login(c fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	resp, err := h.svc.Login(c.Context(), req.Login, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if resp == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	return c.JSON(resp)
}

func (h *Handler) GetTeachers(c fiber.Ctx) error {
	teachers, err := h.svc.GetTeachers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(teachers)
}

func (h *Handler) CreateTeacher(c fiber.Ctx) error {
	var t models.Teacher
	if err := c.Bind().JSON(&t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if err := h.svc.CreateTeacher(c.Context(), &t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(t)
}

func (h *Handler) UpdateTeacher(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var t models.Teacher
	if err := c.Bind().JSON(&t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	t.ID = id

	if err := h.svc.UpdateTeacher(c.Context(), &t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(t)
}

func (h *Handler) DeleteTeacher(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if err := h.svc.DeleteTeacher(c.Context(), id); err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "foreign key constraint") || strings.Contains(errStr, "23503") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "cannot delete teacher: currently assigned to schedule entries"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errStr})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GetSubjects(c fiber.Ctx) error {
	subjects, err := h.svc.GetSubjects(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(subjects)
}

func (h *Handler) CreateSubject(c fiber.Ctx) error {
	var s models.Subject
	if err := c.Bind().JSON(&s); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if err := h.svc.CreateSubject(c.Context(), &s); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(s)
}

func (h *Handler) UpdateSubject(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var s models.Subject
	if err := c.Bind().JSON(&s); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	s.ID = id

	if err := h.svc.UpdateSubject(c.Context(), &s); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(s)
}

func (h *Handler) DeleteSubject(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if err := h.svc.DeleteSubject(c.Context(), id); err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "foreign key constraint") || strings.Contains(errStr, "23503") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "cannot delete subject: currently assigned to schedule entries"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errStr})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GetClassrooms(c fiber.Ctx) error {
	classrooms, err := h.svc.GetClassrooms(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(classrooms)
}

func (h *Handler) CreateClassroom(c fiber.Ctx) error {
	var cl models.Classroom
	if err := c.Bind().JSON(&cl); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if err := h.svc.CreateClassroom(c.Context(), &cl); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(cl)
}

func (h *Handler) UpdateClassroom(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var cl models.Classroom
	if err := c.Bind().JSON(&cl); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	cl.ID = id

	if err := h.svc.UpdateClassroom(c.Context(), &cl); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(cl)
}

func (h *Handler) DeleteClassroom(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if err := h.svc.DeleteClassroom(c.Context(), id); err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "foreign key constraint") || strings.Contains(errStr, "23503") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "cannot delete classroom: currently assigned to schedule entries"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errStr})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GetGroups(c fiber.Ctx) error {
	groups, err := h.svc.GetGroups(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(groups)
}

func (h *Handler) CreateGroup(c fiber.Ctx) error {
	var g models.Group
	if err := c.Bind().JSON(&g); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if err := h.svc.CreateGroup(c.Context(), &g); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(g)
}

func (h *Handler) UpdateGroup(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var g models.Group
	if err := c.Bind().JSON(&g); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	g.ID = id

	if err := h.svc.UpdateGroup(c.Context(), &g); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(g)
}

func (h *Handler) DeleteGroup(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if err := h.svc.DeleteGroup(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GetSchedule(c fiber.Ctx) error {
	groupID, _ := strconv.Atoi(c.Query("group_id", "0"))
	teacherID, _ := strconv.Atoi(c.Query("teacher_id", "0"))
	date := c.Query("date")

	schedules, err := h.svc.GetSchedule(c.Context(), groupID, teacherID, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(schedules)
}

func (h *Handler) CreateScheduleEntry(c fiber.Ctx) error {
	var req models.CreateEntryRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	role := c.Locals("role").(string)

	conflictResult, err := h.validator.CheckConflicts(c.Context(), &req, 0)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if conflictResult.HasConflict && role == "admin" && !req.ForceSave {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   "conflict detected",
			"details": conflictResult.Warnings,
		})
	}

	entry, err := h.svc.CreateScheduleEntry(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	result := fiber.Map{"entry": entry}
	if conflictResult.HasConflict {
		result["warnings"] = conflictResult.Warnings
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

func (h *Handler) UpdateScheduleEntry(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var req models.CreateEntryRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	role := c.Locals("role").(string)

	conflictResult, err := h.validator.CheckConflicts(c.Context(), &req, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if conflictResult.HasConflict && role == "admin" && !req.ForceSave {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   "conflict detected",
			"details": conflictResult.Warnings,
		})
	}

	if err := h.svc.UpdateScheduleEntry(c.Context(), id, &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	entry, _ := h.svc.GetScheduleEntryByID(c.Context(), id)
	result := fiber.Map{"entry": entry}
	if conflictResult.HasConflict {
		result["warnings"] = conflictResult.Warnings
	}

	return c.JSON(result)
}

func (h *Handler) DeleteScheduleEntry(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if err := h.svc.DeleteScheduleEntry(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
