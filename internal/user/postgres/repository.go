package user_postgres

import (
	"context"
	"database/sql"
	"log"

	"github.com/joaoramos09/go-ai/internal/errs"
	"github.com/joaoramos09/go-ai/internal/user"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}
func (r *Repository) Create(ctx context.Context, u *user.User) (*user.User, error) {

	// if u.Role == nil {
	// 	var defaultRole user.Role
	// 	u.Role = &defaultRole
	// }

	adminRole := user.RoleAdmin
	u.Role = &adminRole

	query := `
	INSERT INTO users (username, email, password, role)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, u.Username, u.Email, u.Password, u.Role).Scan(&u.ID)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	return u, nil
}

func (r *Repository) Get(ctx context.Context, id int) (*user.User, error) {
	query := `
	SELECT id, username, email, password, role
	FROM users
	WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var u user.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			log.Printf("User not found: %d", id)
			return nil, errs.ErrUserNotFound
		default:
			log.Printf("Error getting user: %v", err)
			return nil, err
		}
	}

	return &u, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
	SELECT id, username, email, password, role
	FROM users
	WHERE email = $1
	`

	row := r.db.QueryRowContext(ctx, query, email)

	var u user.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			log.Printf("User not found: %s", email)
			return nil, errs.ErrUserNotFound
		default:
			log.Printf("Error getting user: %v", err)
			return nil, err
		}
	}

	return &u, nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	query := `
	DELETE FROM users
	WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return errs.ErrUserNotDeleted
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return errs.ErrUserNotDeleted
	}

	if rowsAffected == 0 {
		log.Printf("User not found: %d", id)
		return errs.ErrUserNotFound

	}

	return nil
}
