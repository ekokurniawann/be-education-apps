package repository

import (
	"be-education/dto"
	"be-education/models"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id int64) error
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	UpdateProfileURL(ctx context.Context, userID int64, profileURL string) error
	GetStudentCountsByClass(ctx context.Context) (map[string]int, error)
	GetAdmins(ctx context.Context) ([]*models.User, error)
	GetTotalAdmins(ctx context.Context) (int, error)
	GetMahasiswaUsers(ctx context.Context) ([]*models.User, error)
}

type userRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (name, email, password, class, birthday, role, profile_url, created_at, updated_at)
		VALUES (:name, :email, :password, :class, :birthday, :role, :profile_url, :created_at, :updated_at)
		RETURNING id, created_at, updated_at`

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare named query for user creation: %w", err)
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, user, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepositoryImpl) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT id, name, email, password, class, birthday, role, profile_url, created_at, updated_at
		FROM users
		WHERE id = $1`

	user := &models.User{}
	err := r.db.GetContext(ctx, user, query, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

func (r *userRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, name, email, password, class, birthday, role, profile_url, created_at, updated_at
		FROM users
		WHERE email = $1`

	user := &models.User{}
	err := r.db.GetContext(ctx, user, query, email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *userRepositoryImpl) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET name = :name, email = :email, password = :password, class = :class, birthday = :birthday,
		    role = :role, profile_url = :profile_url, updated_at = :updated_at
		WHERE id = :id`

	user.UpdatedAt = time.Now()

	res, err := r.db.NamedExecContext(
		ctx,
		query,
		user,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found for update", user.ID)
	}
	return nil
}

func (r *userRepositoryImpl) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found for deletion", id)
	}
	return nil
}

func (r *userRepositoryImpl) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	query := `
		SELECT id, name, email, password, class, birthday, role, profile_url, created_at, updated_at
		FROM users`

	users := []*models.User{}
	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	return users, nil
}

func (r *userRepositoryImpl) UpdateProfileURL(ctx context.Context, userID int64, profileURL string) error {
	query := `
		UPDATE users
		SET profile_url = $1, updated_at = $2
		WHERE id = $3`

	updatedAt := time.Now()

	res, err := r.db.ExecContext(ctx, query, profileURL, updatedAt, userID)
	if err != nil {
		return fmt.Errorf("failed to update user profile URL: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for profile URL update: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found for profile URL update", userID)
	}
	return nil
}

func (r *userRepositoryImpl) GetStudentCountsByClass(ctx context.Context) (map[string]int, error) {
	query := `
        SELECT TRIM(class) as class, COUNT(id) as total
        FROM users
        WHERE role = 'mahasiswa'
        GROUP BY TRIM(class)
    `

	var counts []dto.StudentCount
	err := r.db.SelectContext(ctx, &counts, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get student counts by class: %w", err)
	}

	studentCountsMap := make(map[string]int)
	for _, count := range counts {
		studentCountsMap[count.Class] = count.Total
	}

	return studentCountsMap, nil
}

func (r *userRepositoryImpl) GetAdmins(ctx context.Context) ([]*models.User, error) {
	query := `
		SELECT id, name, email, password, class, birthday, role, profile_url, created_at, updated_at
		FROM users
		WHERE role = 'admin'`

	admins := []*models.User{}
	err := r.db.SelectContext(ctx, &admins, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin users: %w", err)
	}
	return admins, nil
}

func (r *userRepositoryImpl) GetTotalAdmins(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(id)
		FROM users
		WHERE role = 'admin'`

	var total int
	err := r.db.GetContext(ctx, &total, query)
	if err != nil {
		return 0, fmt.Errorf("failed to get total admin count: %w", err)
	}
	return total, nil
}

func (r *userRepositoryImpl) GetMahasiswaUsers(ctx context.Context) ([]*models.User, error) {
	query := `
		SELECT id, name, email, password, class, birthday, role, profile_url, created_at, updated_at
		FROM users
		WHERE role = 'mahasiswa'`

	mahasiswaUsers := []*models.User{}
	err := r.db.SelectContext(ctx, &mahasiswaUsers, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get mahasiswa users: %w", err)
	}
	return mahasiswaUsers, nil
}
