package validators

import (
	"context"

	"course-project/internal/models"
	"course-project/internal/repository"
)

type ConflictValidator struct {
	repo *repository.Repository
}

func NewConflictValidator(repo *repository.Repository) *ConflictValidator {
	return &ConflictValidator{repo: repo}
}

func (v *ConflictValidator) CheckConflicts(ctx context.Context, req *models.CreateEntryRequest, excludeEntryID int) (*models.ConflictCheckResult, error) {
	result := &models.ConflictCheckResult{HasConflict: false, Warnings: []models.ConflictWarning{}}

	schedule, err := v.repo.GetOrCreateSchedule(ctx, req.GroupID, req.Date)
	if err != nil {
		return nil, err
	}

	classroomConflict, entryID1, err := v.repo.CheckClassroomConflict(ctx, schedule.ID, req.ClassroomID, req.PairNumber, excludeEntryID)
	if err != nil {
		return nil, err
	}
	if classroomConflict {
		result.HasConflict = true
		classroom, _ := v.repo.GetClassroomByID(ctx, req.ClassroomID)
		classroomStr := "unknown"
		if classroom != nil {
			classroomStr = classroom.Number
		}
		result.Warnings = append(result.Warnings, models.ConflictWarning{
			Type:    "classroom",
			Message: "Аудитория " + classroomStr + " уже занята в эту пару",
			Field:   "classroom_id",
			EntryID: entryID1,
		})
	}

	teacherConflict, entryID2, err := v.repo.CheckTeacherConflict(ctx, schedule.ID, req.TeacherID, req.PairNumber, excludeEntryID)
	if err != nil {
		return nil, err
	}
	if teacherConflict {
		result.HasConflict = true
		teacher, _ := v.repo.GetTeacherByID(ctx, req.TeacherID)
		teacherStr := "unknown"
		if teacher != nil {
			teacherStr = teacher.Name
		}
		result.Warnings = append(result.Warnings, models.ConflictWarning{
			Type:    "teacher",
			Message: "Преподаватель " + teacherStr + " уже занят в эту пару",
			Field:   "teacher_id",
			EntryID: entryID2,
		})
	}

	return result, nil
}
