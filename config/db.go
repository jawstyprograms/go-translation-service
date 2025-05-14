package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

// ConnectDB creates connection with PostgreSQL database
func ConnectDB() (*pgx.Conn,error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL not set")
	}

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	fmt.Println("Connected to PostgreSQL")
	return conn, nil
}


