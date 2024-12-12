package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	DB *pgxpool.Pool
}

func (r *AuthRepository) CreateUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3)"
	_, err := r.DB.Exec(ctx, query, user.Username, user.Email, user.Password)
	if err != nil {
		log.Printf("CreateUser error: %v", err)
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (r *AuthRepository) GetUserByUsername(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT username, email, password FROM users WHERE username=$1"
	row := r.DB.QueryRow(ctx, query, username)

	var user User
	err := row.Scan(&user.Username, &user.Email, &user.Password)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("user not found")
		}
		log.Printf("GetUserByUsername error: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}
