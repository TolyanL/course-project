package models

import "time"

type Teacher struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Login        string    `json:"login"`
	Password     string    `json:"password,omitempty"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type Subject struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Classroom struct {
	ID     int    `json:"id"`
	Number string `json:"number"`
}

type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ScheduleEntry struct {
	ID          int        `json:"id"`
	ScheduleID  int        `json:"schedule_id"`
	SubjectID   int        `json:"subject_id"`
	TeacherID   int        `json:"teacher_id"`
	ClassroomID int        `json:"classroom_id"`
	PairNumber  int        `json:"pair_number"`
	Subject     *Subject   `json:"subject,omitempty"`
	Teacher     *Teacher   `json:"teacher,omitempty"`
	Classroom   *Classroom `json:"classroom,omitempty"`
}

type Schedule struct {
	ID      int             `json:"id"`
	GroupID int             `json:"group_id"`
	Date    string          `json:"date"`
	Pairs   []ScheduleEntry `json:"pairs"`
	Group   *Group          `json:"group,omitempty"`
}

type PairSlot struct {
	SubjectID   int `json:"subject_id"`
	TeacherID   int `json:"teacher_id"`
	ClassroomID int `json:"classroom_id"`
	PairNumber  int `json:"pair_number"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

type CreateEntryRequest struct {
	GroupID     int    `json:"group_id"`
	Date        string `json:"date"`
	SubjectID   int    `json:"subject_id"`
	TeacherID   int    `json:"teacher_id"`
	ClassroomID int    `json:"classroom_id"`
	PairNumber  int    `json:"pair_number"`
	ForceSave   bool   `json:"force_save"`
}

type ConflictWarning struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Field   string `json:"field"`
	EntryID int    `json:"entry_id,omitempty"`
}

type ConflictCheckResult struct {
	HasConflict bool              `json:"has_conflict"`
	Warnings    []ConflictWarning `json:"warnings,omitempty"`
}
