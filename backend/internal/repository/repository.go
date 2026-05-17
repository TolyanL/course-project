package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"course-project/internal/models"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetTeacherByLogin(ctx context.Context, login string) (*models.Teacher, error) {
	var t models.Teacher
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, login, password_hash, role, created_at 
		FROM teachers WHERE login = $1
	`, login).Scan(&t.ID, &t.Name, &t.Login, &t.PasswordHash, &t.Role, &t.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &t, err
}

func (r *Repository) CreateTeacher(ctx context.Context, t *models.Teacher) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(t.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	t.PasswordHash = string(hash)

	return r.pool.QueryRow(ctx, `
		INSERT INTO teachers (name, login, password_hash, role) 
		VALUES ($1, $2, $3, $4) RETURNING id
	`, t.Name, t.Login, t.PasswordHash, t.Role).Scan(&t.ID)
}

func (r *Repository) GetTeachers(ctx context.Context) ([]models.Teacher, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, login, role, created_at FROM teachers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []models.Teacher
	for rows.Next() {
		var t models.Teacher
		if err := rows.Scan(&t.ID, &t.Name, &t.Login, &t.Role, &t.CreatedAt); err != nil {
			return nil, err
		}
		teachers = append(teachers, t)
	}
	if teachers == nil {
		teachers = []models.Teacher{}
	}
	return teachers, rows.Err()
}

func (r *Repository) GetTeacherByID(ctx context.Context, id int) (*models.Teacher, error) {
	var t models.Teacher
	err := r.pool.QueryRow(ctx, `SELECT id, name, login, role, created_at FROM teachers WHERE id = $1`, id).
		Scan(&t.ID, &t.Name, &t.Login, &t.Role, &t.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &t, err
}

func (r *Repository) UpdateTeacher(ctx context.Context, t *models.Teacher) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE teachers SET name = $1, role = $2 WHERE id = $3
	`, t.Name, t.Role, t.ID)
	return err
}

func (r *Repository) DeleteTeacher(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM teachers WHERE id = $1`, id)
	return err
}

func (r *Repository) GetSubjects(ctx context.Context) ([]models.Subject, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name FROM subjects`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		var s models.Subject
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, err
		}
		subjects = append(subjects, s)
	}
	if subjects == nil {
		subjects = []models.Subject{}
	}
	return subjects, rows.Err()
}

func (r *Repository) CreateSubject(ctx context.Context, s *models.Subject) error {
	return r.pool.QueryRow(ctx, `INSERT INTO subjects (name) VALUES ($1) RETURNING id`, s.Name).Scan(&s.ID)
}

func (r *Repository) GetSubjectByID(ctx context.Context, id int) (*models.Subject, error) {
	var s models.Subject
	err := r.pool.QueryRow(ctx, `SELECT id, name FROM subjects WHERE id = $1`, id).Scan(&s.ID, &s.Name)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &s, err
}

func (r *Repository) UpdateSubject(ctx context.Context, s *models.Subject) error {
	_, err := r.pool.Exec(ctx, `UPDATE subjects SET name = $1 WHERE id = $2`, s.Name, s.ID)
	return err
}

func (r *Repository) DeleteSubject(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM subjects WHERE id = $1`, id)
	return err
}

func (r *Repository) GetClassrooms(ctx context.Context) ([]models.Classroom, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, number FROM classrooms`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classrooms []models.Classroom
	for rows.Next() {
		var c models.Classroom
		if err := rows.Scan(&c.ID, &c.Number); err != nil {
			return nil, err
		}
		classrooms = append(classrooms, c)
	}
	if classrooms == nil {
		classrooms = []models.Classroom{}
	}
	return classrooms, rows.Err()
}

func (r *Repository) CreateClassroom(ctx context.Context, c *models.Classroom) error {
	return r.pool.QueryRow(ctx, `INSERT INTO classrooms (number) VALUES ($1) RETURNING id`, c.Number).Scan(&c.ID)
}

func (r *Repository) GetClassroomByID(ctx context.Context, id int) (*models.Classroom, error) {
	var c models.Classroom
	err := r.pool.QueryRow(ctx, `SELECT id, number FROM classrooms WHERE id = $1`, id).Scan(&c.ID, &c.Number)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &c, err
}

func (r *Repository) UpdateClassroom(ctx context.Context, c *models.Classroom) error {
	_, err := r.pool.Exec(ctx, `UPDATE classrooms SET number = $1 WHERE id = $2`, c.Number, c.ID)
	return err
}

func (r *Repository) DeleteClassroom(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM classrooms WHERE id = $1`, id)
	return err
}

func (r *Repository) GetGroups(ctx context.Context) ([]models.Group, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name FROM groups`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var g models.Group
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	if groups == nil {
		groups = []models.Group{}
	}
	return groups, rows.Err()
}

func (r *Repository) CreateGroup(ctx context.Context, g *models.Group) error {
	return r.pool.QueryRow(ctx, `INSERT INTO groups (name) VALUES ($1) RETURNING id`, g.Name).Scan(&g.ID)
}

func (r *Repository) GetGroupByID(ctx context.Context, id int) (*models.Group, error) {
	var g models.Group
	err := r.pool.QueryRow(ctx, `SELECT id, name FROM groups WHERE id = $1`, id).Scan(&g.ID, &g.Name)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &g, err
}

func (r *Repository) UpdateGroup(ctx context.Context, g *models.Group) error {
	_, err := r.pool.Exec(ctx, `UPDATE groups SET name = $1 WHERE id = $2`, g.Name, g.ID)
	return err
}

func (r *Repository) DeleteGroup(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM groups WHERE id = $1`, id)
	return err
}

func (r *Repository) GetOrCreateSchedule(ctx context.Context, groupID int, date string) (*models.Schedule, error) {
	var s models.Schedule
	var dateText string
	err := r.pool.QueryRow(ctx, `
		INSERT INTO schedules (group_id, date) VALUES ($1, $2::date)
		ON CONFLICT (group_id, date) DO UPDATE SET date = EXCLUDED.date
		RETURNING id, group_id, date::text
	`, groupID, date).Scan(&s.ID, &s.GroupID, &dateText)
	if err != nil {
		return nil, err
	}
	s.Date = dateText
	return &s, nil
}

func (r *Repository) GetScheduleByFilters(ctx context.Context, groupID int, teacherID int, date string) ([]models.Schedule, error) {
	query := `
		SELECT s.id, s.group_id, s.date::text, g.name
		FROM schedules s
		JOIN groups g ON g.id = s.group_id
		WHERE 1=1`
	args := []any{}
	idx := 1

	if groupID > 0 {
		query += fmt.Sprintf(" AND s.group_id = $%d", idx)
		args = append(args, groupID)
		idx++
	}
	if teacherID > 0 {
		query += fmt.Sprintf(" AND s.id IN (SELECT schedule_id FROM schedule_entries WHERE teacher_id = $%d)", idx)
		args = append(args, teacherID)
		idx++
	}
	if date != "" {
		query += fmt.Sprintf(" AND s.date = $%d", idx)
		args = append(args, date)
	}

	query += " ORDER BY s.date"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var s models.Schedule
		var gName string
		if err := rows.Scan(&s.ID, &s.GroupID, &s.Date, &gName); err != nil {
			return nil, err
		}
		s.Group = &models.Group{ID: s.GroupID, Name: gName}
		schedules = append(schedules, s)
	}
	return schedules, rows.Err()
}

func (r *Repository) GetScheduleEntries(ctx context.Context, scheduleID int) ([]models.ScheduleEntry, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT se.id, se.schedule_id, se.subject_id, se.teacher_id, se.classroom_id, se.pair_number,
			   sub.name, sub.id, t.name, t.id, c.number, c.id
		FROM schedule_entries se
		JOIN subjects sub ON sub.id = se.subject_id
		JOIN teachers t ON t.id = se.teacher_id
		JOIN classrooms c ON c.id = se.classroom_id
		WHERE se.schedule_id = $1
		ORDER BY se.pair_number
	`, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.ScheduleEntry
	for rows.Next() {
		var e models.ScheduleEntry
		var subName, teachName, classNum string
		var subID, teachID, classID int
		if err := rows.Scan(&e.ID, &e.ScheduleID, &e.SubjectID, &e.TeacherID, &e.ClassroomID, &e.PairNumber,
			&subName, &subID, &teachName, &teachID, &classNum, &classID); err != nil {
			return nil, err
		}
		e.Subject = &models.Subject{ID: subID, Name: subName}
		e.Teacher = &models.Teacher{ID: teachID, Name: teachName}
		e.Classroom = &models.Classroom{ID: classID, Number: classNum}
		entries = append(entries, e)
	}
	if entries == nil {
		entries = []models.ScheduleEntry{}
	}
	return entries, rows.Err()
}

func (r *Repository) CheckClassroomConflict(ctx context.Context, scheduleID int, classroomID int, pairNumber int, excludeEntryID int) (bool, int, error) {
	var exists bool
	var entryID *int
	err := r.pool.QueryRow(ctx, `
		SELECT 
			EXISTS(
				SELECT 1 FROM schedule_entries 
				WHERE schedule_id = $1 AND classroom_id = $2 AND pair_number = $3 AND id != $4
			),
			(SELECT id FROM schedule_entries 
				WHERE schedule_id = $1 AND classroom_id = $2 AND pair_number = $3 AND id != $4 LIMIT 1)
	`, scheduleID, classroomID, pairNumber, excludeEntryID).Scan(&exists, &entryID)
	if entryID == nil {
		return exists, 0, err
	}
	return exists, *entryID, err
}

func (r *Repository) CheckTeacherConflict(ctx context.Context, scheduleID int, teacherID int, pairNumber int, excludeEntryID int) (bool, int, error) {
	var exists bool
	var entryID *int
	err := r.pool.QueryRow(ctx, `
		SELECT 
			EXISTS(
				SELECT 1 FROM schedule_entries 
				WHERE schedule_id = $1 AND teacher_id = $2 AND pair_number = $3 AND id != $4
			),
			(SELECT id FROM schedule_entries 
				WHERE schedule_id = $1 AND teacher_id = $2 AND pair_number = $3 AND id != $4 LIMIT 1)
	`, scheduleID, teacherID, pairNumber, excludeEntryID).Scan(&exists, &entryID)
	if entryID == nil {
		return exists, 0, err
	}
	return exists, *entryID, err
}

func (r *Repository) CreateScheduleEntry(ctx context.Context, req *models.CreateEntryRequest) (*models.ScheduleEntry, error) {
	schedule, err := r.GetOrCreateSchedule(ctx, req.GroupID, req.Date)
	if err != nil {
		return nil, err
	}

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	_, err = conn.Conn().Exec(ctx, `
		INSERT INTO schedule_entries (schedule_id, subject_id, teacher_id, classroom_id, pair_number)
		VALUES ($1, $2, $3, $4, $5)
	`, schedule.ID, req.SubjectID, req.TeacherID, req.ClassroomID, req.PairNumber)
	if err != nil {
		return nil, err
	}

	var entryID int
	err = conn.Conn().QueryRow(ctx, `SELECT lastval()`).Scan(&entryID)
	if err != nil {
		return nil, err
	}

	return &models.ScheduleEntry{
		ID:          entryID,
		ScheduleID:  schedule.ID,
		SubjectID:   req.SubjectID,
		TeacherID:   req.TeacherID,
		ClassroomID: req.ClassroomID,
		PairNumber:  req.PairNumber,
	}, nil
}

func (r *Repository) GetScheduleEntryByID(ctx context.Context, id int) (*models.ScheduleEntry, error) {
	var entry models.ScheduleEntry
	err := r.pool.QueryRow(ctx, `
		SELECT id, schedule_id, subject_id, teacher_id, classroom_id, pair_number
		FROM schedule_entries WHERE id = $1
	`, id).Scan(&entry.ID, &entry.ScheduleID, &entry.SubjectID, &entry.TeacherID, &entry.ClassroomID, &entry.PairNumber)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &entry, err
}

func (r *Repository) DB() *pgxpool.Pool {
	return r.pool
}

func (r *Repository) UpdateScheduleEntry(ctx context.Context, id int, req *models.CreateEntryRequest) error {
	schedule, err := r.GetOrCreateSchedule(ctx, req.GroupID, req.Date)
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, `
		UPDATE schedule_entries 
		SET schedule_id = $1, subject_id = $2, teacher_id = $3, classroom_id = $4, pair_number = $5
		WHERE id = $6
	`, schedule.ID, req.SubjectID, req.TeacherID, req.ClassroomID, req.PairNumber, id)
	return err
}

func (r *Repository) DeleteScheduleEntry(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM schedule_entries WHERE id = $1`, id)
	return err
}

func (r *Repository) CheckTeacherPassword(ctx context.Context, login, password string) (*models.Teacher, error) {
	teacher, err := r.GetTeacherByLogin(ctx, login)
	if err != nil || teacher == nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(teacher.PasswordHash), []byte(password)); err != nil {
		return nil, nil
	}
	return teacher, nil
}
