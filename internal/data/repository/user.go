package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/project-app-bioskop-golang/internal/data/entity"
	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/pkg/database"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type UserRepository interface {
	Create(user *entity.User) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	GetAll(q dto.PaginationQuery) ([]entity.User, int, error)
	GetByID(id int) (entity.User, error)
	Update(id int, data *entity.User) error
	Delete(id int) error
}

type userRepository struct {
	db database.PgxIface
	Logger *zap.Logger
}

func NewUserRepository(db database.PgxIface, log *zap.Logger) UserRepository {
	return &userRepository{
		db: db,
		Logger: log,
	}
}

func (r *userRepository) Create(user *entity.User) (*entity.User, error) {
	query := `
		INSERT INTO users (name, email, password, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id
	`
	err := r.db.QueryRow(context.Background(), query, user.Name, user.Email, user.Password, user.Role).Scan(&user.ID)
	if err != nil {
		r.Logger.Error("Error query create user: ", zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	query := `
		SELECT id, created_at, updated_at, name, email, password, role
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`
	var user entity.User
	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
		&user.Name, &user.Email, &user.Password, &user.Role,
	)

	if err != nil {
		return nil, err
	}

	return &user, err
}

func (r *userRepository) GetAll(q dto.PaginationQuery) ([]entity.User, int, error) {
	var offset int
	offset = (q.Page - 1) * q.Limit
	
	// Get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
	err := r.db.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("Error query count users: ", zap.Error(err))
		return nil, 0, err
	}

	// Initiate rows
	var rows pgx.Rows
	
	// Conditional query based on page, limit, and all param
	query := `SELECT id, name, email, role, created_at, updated_at FROM users WHERE deleted_at IS NULL`

	if !q.All && q.Limit > 0 {
		query += ` LIMIT $1 OFFSET $2`
		rows, err = r.db.Query(context.Background(), query, q.Limit, offset)
	} else {
		rows, err = r.db.Query(context.Background(), query)
	}
	
	if err != nil {
		r.Logger.Error("Error query get all users: ", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var u entity.User
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			r.Logger.Error("Error scan user: ", zap.Error(err))
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, nil
}

func (r *userRepository) GetByID(id int) (entity.User, error) {
	var user entity.User
	query := "SELECT id, name, email, role, created_at, updated_at FROM users WHERE id = $1"

	err := r.db.QueryRow(context.Background(), query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err == pgx.ErrNoRows {
		r.Logger.Error("Error not found user: ", zap.Error(err))
		return user, utils.ErrNotFound("user")
	}

	if err != nil {
		r.Logger.Error("Error query get user by id: ", zap.Error(err))
		return user, err
	}

	return user, nil
}

func (r *userRepository) Update(id int, u *entity.User) error {
	query := `
		UPDATE users
		SET name = COALESCE($1, name),
		email = COALESCE($2, email),
		role = COALESCE($3, role),
		updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(context.Background(), query,
		&u.Name, &u.Email, &u.Role, id,
	)

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.Logger.Error("Error not found user: ", zap.Error(err))
		return utils.ErrNotFound("user")
	}

	if err != nil {
		r.Logger.Error("Error query update user: ", zap.Error(err))
		return err
	}

	return nil
}

func (r *userRepository) Delete(id int) error {
	query := `
		UPDATE users
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(context.Background(), query, id)

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.Logger.Error("Error not found user: ", zap.Error(err))
		return utils.ErrNotFound("user")
	}

	if err != nil {
		r.Logger.Error("Error query delete user: ", zap.Error(err))
		return err
	}

	return nil
}