package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	files, err := filepath.Glob(filepath.Join("db", "migrations", "*.sql"))
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return nil
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", file, err)
		}
		if _, err := pool.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("execute migration %s: %w", file, err)
		}
	}
	return nil
}
