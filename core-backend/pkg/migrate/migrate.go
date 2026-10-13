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
