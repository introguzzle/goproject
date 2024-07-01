package migration

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sort"
)

var dir = "migrations"

func Start(db *sql.DB) {
	files, err := getMigrationFiles(dir)
	if err != nil {
		log.Fatalf("Failed to list migration files: %v", err)
	}

	for _, file := range files {
		start(db, file)
	}
}

func getMigrationFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".sql" {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Strings(files)
	return files, nil
}

func start(db *sql.DB, migration string) {
	sqlBytes, err := os.ReadFile(migration)
	if err != nil {
		log.Fatalf("Failed to read migration file %s: %v \n", migration, err)
	}

	sqlQuery := string(sqlBytes)

	_, err = db.Exec(sqlQuery)
	if err != nil {
		log.Fatalf("Failed to execute migration file %s: %v \n", migration, err)
	}

	log.Printf("Migration %s executed successfully \n", migration)
}
