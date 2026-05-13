package validators

import (
	"context"
	"testing"

	"course-project/internal/models"
)

type ConflictChecker interface {
	GetOrCreateSchedule(ctx context.Context, groupID int, date string) (*models.Schedule, error)
	CheckClassroomConflict(ctx context.Context, scheduleID int, classroomID int, pairNumber int, excludeEntryID int) (bool, int, error)
	CheckTeacherConflict(ctx context.Context, scheduleID int, teacherID int, pairNumber int, excludeEntryID int) (bool, int, error)
	GetClassroomByID(ctx context.Context, id int) (*models.Classroom, error)
	GetTeacherByID(ctx context.Context, id int) (*models.Teacher, error)
}

type mockConflictRepo struct {
	classroomConflicts map[string]int
	teacherConflicts   map[string]int
	classrooms         map[int]*models.Classroom
	teachers           map[int]*models.Teacher
}

func newMockConflictRepo() *mockConflictRepo {
	return &mockConflictRepo{
		classroomConflicts: make(map[string]int),
		teacherConflicts:   make(map[string]int),
		classrooms:         make(map[int]*models.Classroom),
		teachers:           make(map[int]*models.Teacher),
	}
}

func (m *mockConflictRepo) GetOrCreateSchedule(ctx context.Context, groupID int, date string) (*models.Schedule, error) {
	return &models.Schedule{ID: groupID, GroupID: groupID, Date: date}, nil
}

func (m *mockConflictRepo) CheckClassroomConflict(ctx context.Context, scheduleID int, classroomID int, pairNumber int, excludeEntryID int) (bool, int, error) {
	key := formatKey(scheduleID, classroomID, pairNumber)
	if conflictingID, exists := m.classroomConflicts[key]; exists && conflictingID != excludeEntryID {
		return true, conflictingID, nil
	}
	return false, 0, nil
}

func (m *mockConflictRepo) CheckTeacherConflict(ctx context.Context, scheduleID int, teacherID int, pairNumber int, excludeEntryID int) (bool, int, error) {
	key := formatKey(scheduleID, teacherID, pairNumber)
	if conflictingID, exists := m.teacherConflicts[key]; exists && conflictingID != excludeEntryID {
		return true, conflictingID, nil
	}
	return false, 0, nil
}

func (m *mockConflictRepo) GetClassroomByID(ctx context.Context, id int) (*models.Classroom, error) {
	return m.classrooms[id], nil
}

func (m *mockConflictRepo) GetTeacherByID(ctx context.Context, id int) (*models.Teacher, error) {
	return m.teachers[id], nil
}

func formatKey(scheduleID, resourceID, pairNumber int) string {
	return string(rune('0'+scheduleID%10)) + "-" + string(rune('0'+resourceID%10)) + "-" + string(rune('0'+pairNumber%10))
}

type testConflictValidator struct {
	repo ConflictChecker
}

func newTestConflictValidator(repo ConflictChecker) *testConflictValidator {
	return &testConflictValidator{repo: repo}
}

func (v *testConflictValidator) CheckConflicts(ctx context.Context, req *models.CreateEntryRequest, excludeEntryID int) (*models.ConflictCheckResult, error) {
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

func TestCheckConflicts_NoConflicts(t *testing.T) {
	repo := newMockConflictRepo()
	repo.classrooms[1] = &models.Classroom{ID: 1, Number: "101"}
	repo.teachers[1] = &models.Teacher{ID: 1, Name: "John"}

	validator := newTestConflictValidator(repo)

	req := &models.CreateEntryRequest{
		GroupID:     1,
		Date:        "2024-01-15",
		ClassroomID: 1,
		TeacherID:   1,
		PairNumber:  1,
	}

	result, err := validator.CheckConflicts(context.Background(), req, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasConflict {
		t.Error("expected no conflicts")
	}
}

func TestCheckConflicts_ClassroomConflict(t *testing.T) {
	repo := newMockConflictRepo()
	repo.classroomConflicts["1-1-1"] = 999
	repo.classrooms[1] = &models.Classroom{ID: 1, Number: "101"}
	repo.teachers[1] = &models.Teacher{ID: 1, Name: "John"}

	validator := newTestConflictValidator(repo)

	req := &models.CreateEntryRequest{
		GroupID:     1,
		Date:        "2024-01-15",
		ClassroomID: 1,
		TeacherID:   1,
		PairNumber:  1,
	}

	result, err := validator.CheckConflicts(context.Background(), req, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasConflict {
		t.Error("expected conflict")
	}
	if len(result.Warnings) == 0 {
		t.Error("expected warnings")
	}
	if result.Warnings[0].Type != "classroom" {
		t.Errorf("expected classroom conflict, got %s", result.Warnings[0].Type)
	}
}

func TestCheckConflicts_TeacherConflict(t *testing.T) {
	repo := newMockConflictRepo()
	repo.teacherConflicts["1-1-1"] = 999
	repo.classrooms[1] = &models.Classroom{ID: 1, Number: "101"}
	repo.teachers[1] = &models.Teacher{ID: 1, Name: "John"}

	validator := newTestConflictValidator(repo)

	req := &models.CreateEntryRequest{
		GroupID:     1,
		Date:        "2024-01-15",
		ClassroomID: 1,
		TeacherID:   1,
		PairNumber:  1,
	}

	result, err := validator.CheckConflicts(context.Background(), req, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasConflict {
		t.Error("expected conflict")
	}
	if len(result.Warnings) == 0 {
		t.Error("expected warnings")
	}
	if result.Warnings[0].Type != "teacher" {
		t.Errorf("expected teacher conflict, got %s", result.Warnings[0].Type)
	}
}

func TestCheckConflicts_BothConflicts(t *testing.T) {
	repo := newMockConflictRepo()
	repo.classroomConflicts["1-1-1"] = 999
	repo.teacherConflicts["1-1-1"] = 998
	repo.classrooms[1] = &models.Classroom{ID: 1, Number: "101"}
	repo.teachers[1] = &models.Teacher{ID: 1, Name: "John"}

	validator := newTestConflictValidator(repo)

	req := &models.CreateEntryRequest{
		GroupID:     1,
		Date:        "2024-01-15",
		ClassroomID: 1,
		TeacherID:   1,
		PairNumber:  1,
	}

	result, err := validator.CheckConflicts(context.Background(), req, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasConflict {
		t.Error("expected conflict")
	}
	if len(result.Warnings) != 2 {
		t.Errorf("expected 2 warnings, got %d", len(result.Warnings))
	}
}

func TestCheckConflicts_ExcludeEntry(t *testing.T) {
	repo := newMockConflictRepo()
	repo.classroomConflicts["1-1-1"] = 1
	repo.classrooms[1] = &models.Classroom{ID: 1, Number: "101"}

	validator := newTestConflictValidator(repo)

	req := &models.CreateEntryRequest{
		GroupID:     1,
		Date:        "2024-01-15",
		ClassroomID: 1,
		TeacherID:   1,
		PairNumber:  1,
	}

	result, err := validator.CheckConflicts(context.Background(), req, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasConflict {
		t.Error("should not conflict when excluding own entry")
	}
}
