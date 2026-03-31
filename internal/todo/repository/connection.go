package repository

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func ConnectRepository(ctx context.Context) (*pgx.Conn, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return pgx.Connect(ctx, os.Getenv("CONN_STRING"))
}
