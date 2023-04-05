package mysql

import (
	"context"
	"database/sql"
	"hunt/pkg/model"
	"time"

	"hunt/internal/repository"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

// constructor function for instantiate a repository struct
func New() (*Repository, error) {
	db, err := sql.Open(os.Getenv("DNS_ACCESS_DATABASE"))
	if err != nil {
		return &Repository{}, err
	}

	// For increased flexibility, it can be configured using environment variables or flags
	db.SetConnMaxIdleTime(45 * time.Second)
	db.SetMaxIdleConns(7)
	db.SetMaxOpenConns(25)

	return &Repository{db}, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *model.User) error {
	_, err := r.db.ExecContext(ctx, `
	INSERT INTO users (id,
		level,
		created_at,
		updated_at,
		profile_picture,
		first_name,
		last_name,
		email,
		password
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`, user.ID, user.Level, user.CreatedAt, user.UpdatedAt, user.ProfilePicture, user.Name, user.Name, user.Email, user.Password.Hash)

	return err
}

func (r *Repository) GetUser(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	var level int
	var name, profile_picture string
	res := r.db.QueryRowContext(ctx, "SELECT level, first_name, profile_picture FROM users WHERE id = ?", userID)

	if err := res.Scan(&level, &name, &profile_picture); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &model.User{
		Level:          level,
		ProfilePicture: profile_picture,
		Name:           name,
	}, nil

}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var level int
	var name, profile_picture, password string
	var ID uuid.UUID
	res := r.db.QueryRowContext(ctx, "SELECT id, level, first_name, profile_picture, password FROM users WHERE email = ?", email)

	if err := res.Scan(&ID, &level, &name, &profile_picture, &password); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	var pass model.Password
	pass.Hash = []byte(password)

	return &model.User{
		ID:             ID,
		Level:          level,
		ProfilePicture: profile_picture,
		Name:           name,
		Password:       pass,
	}, nil

}
