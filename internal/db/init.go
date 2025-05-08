package db

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type MigrationState struct {
	Filename  string
	Applied   bool
	AppliedAt time.Time
}

func (p *PostgreSQL) InitDatabase() error {
	if err := p.db.Ping(); err != nil {
		return fmt.Errorf("ошибка соединения с базой данных: %w", err)
	}

	migrationsDir := filepath.Join("internal", "db", "migrations")

	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return fmt.Errorf("директория с миграциями не найдена: %s", migrationsDir)
	}

	if err := p.ensureBasicStructureExists(); err != nil {
		return err
	}

	migrations, err := p.getMigrationStates(migrationsDir)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if !migration.Applied {
			if err := p.applyMigration(migrationsDir, migration.Filename); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *PostgreSQL) ensureBasicStructureExists() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS migration_history (
		id SERIAL PRIMARY KEY,
		filename VARCHAR(255) NOT NULL UNIQUE,
		applied_at TIMESTAMP NOT NULL DEFAULT NOW()
	);
	`

	_, err := p.db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы migration_history: %w", err)
	}

	var count int
	err = p.db.QueryRow("SELECT COUNT(*) FROM migration_history WHERE filename = '0001_schema.sql'").Scan(&count)

	if err != nil || count == 0 {
		schemaPath := filepath.Join("internal", "db", "migrations", "0001_schema.sql")
		schemaSQL, err := os.ReadFile(schemaPath)
		if err != nil {
			return fmt.Errorf("ошибка при чтении файла 0001_schema.sql: %w", err)
		}

		tx, err := p.db.Begin()
		if err != nil {
			return fmt.Errorf("ошибка при начале транзакции для схемы: %w", err)
		}

		_, err = tx.Exec(string(schemaSQL))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("ошибка при выполнении схемы: %w", err)
		}

		_, err = tx.Exec("INSERT INTO migration_history (filename) VALUES ($1) ON CONFLICT (filename) DO NOTHING", "0001_schema.sql")
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("ошибка при записи информации о схеме: %w", err)
		}

		if err = tx.Commit(); err != nil {
			return fmt.Errorf("ошибка при фиксации транзакции для схемы: %w", err)
		}

		fmt.Printf("Успешно применена основная схема базы данных\n")
	}

	return nil
}

func (p *PostgreSQL) getMigrationStates(migrationsDir string) ([]MigrationState, error) {
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении директории миграций: %w", err)
	}

	var sqlFiles []string
	for _, file := range files {
		filename := file.Name()
		if !file.IsDir() && filepath.Ext(filename) == ".sql" {
			sqlFiles = append(sqlFiles, filename)
		}
	}
	sort.Strings(sqlFiles)

	appliedMigrations := make(map[string]time.Time)
	rows, err := p.db.Query("SELECT filename, applied_at FROM migration_history")
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
		} else {
			return nil, fmt.Errorf("ошибка при получении списка примененных миграций: %w", err)
		}
	} else {
		defer rows.Close()

		for rows.Next() {
			var filename string
			var appliedAt time.Time
			if err := rows.Scan(&filename, &appliedAt); err != nil {
				return nil, fmt.Errorf("ошибка при чтении данных о миграции: %w", err)
			}
			appliedMigrations[filename] = appliedAt
		}
	}

	var result []MigrationState
	for _, filename := range sqlFiles {
		appliedAt, isApplied := appliedMigrations[filename]
		result = append(result, MigrationState{
			Filename:  filename,
			Applied:   isApplied,
			AppliedAt: appliedAt,
		})
	}

	return result, nil
}

func (p *PostgreSQL) applyMigration(migrationsDir, filename string) error {
	filePath := filepath.Join(migrationsDir, filename)
	sqlContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("ошибка при чтении файла %s: %w", filename, err)
	}

	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка при начале транзакции для миграции %s: %w", filename, err)
	}

	_, err = tx.Exec(string(sqlContent))
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка при выполнении миграции %s: %w", filename, err)
	}

	_, err = tx.Exec("INSERT INTO migration_history (filename) VALUES ($1) ON CONFLICT (filename) DO NOTHING", filename)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка при записи информации о миграции %s: %w", filename, err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка при фиксации транзакции для миграции %s: %w", filename, err)
	}

	fmt.Printf("Успешно применена миграция: %s\n", filename)
	return nil
}
