package backup

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Service interface {
	CreateBackup(ctx context.Context) (string, string, error)
	RestoreBackup(ctx context.Context, path string) error
}

type service struct {
	dsn       string
	pgDumpBin string
	pgRestore string
}

func NewService(dsn string) Service {
	return &service{
		dsn:       dsn,
		pgDumpBin: getenvDefault("PG_DUMP_PATH", "pg_dump"),
		pgRestore: getenvDefault("PG_RESTORE_PATH", "pg_restore"),
	}
}

func (s *service) CreateBackup(ctx context.Context) (string, string, error) {
	if s.dsn == "" {
		return "", "", fmt.Errorf("backup: missing DATABASE_URL")
	}
	name := "physio_" + time.Now().Format("20060102_150405") + ".dump"
	tmpDir := os.TempDir()
	target := filepath.Join(tmpDir, name)

	cmd := exec.CommandContext(ctx, s.pgDumpBin, "-Fc", "-f", target, s.dsn)
	cmd.Env = os.Environ()
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", "", fmt.Errorf("backup: pg_dump failed: %w (%s)", err, string(out))
	}
	return name, target, nil
}

func (s *service) RestoreBackup(ctx context.Context, path string) error {
	if s.dsn == "" {
		return fmt.Errorf("restore: missing DATABASE_URL")
	}
	cmd := exec.CommandContext(ctx, s.pgRestore, "--clean", "--if-exists", "-d", s.dsn, path)
	cmd.Env = os.Environ()
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("restore: pg_restore failed: %w (%s)", err, string(out))
	}
	return nil
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
