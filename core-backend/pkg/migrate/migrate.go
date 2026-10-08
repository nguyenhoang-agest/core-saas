package migrate

import (
	"context"
	"fmt"
	"io/fs"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

var migrateMu sync.Mutex

// Run executes database migrations for a specific module from an embedded filesystem.
// Migrations are tracked in a dedicated version table `goose_version_<moduleName>`,
// keeping module schemas isolated and ready for future service extraction.
func Run(ctx context.Context, pool *pgxpool.Pool, moduleName string, fsys fs.FS, dir string) error {
	migrateMu.Lock()
	defer migrateMu.Unlock()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("setting goose dialect: %w", err)
	}

	tableName := fmt.Sprintf("goose_version_%s", moduleName)
	goose.SetTableName(tableName)
	goose.SetBaseFS(fsys)
	defer goose.SetBaseFS(nil)

	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()

	if err := goose.UpContext(ctx, sqlDB, dir); err != nil {
		return fmt.Errorf("running migrations for module %q (table %s): %w", moduleName, tableName, err)
	}

	return nil
}
