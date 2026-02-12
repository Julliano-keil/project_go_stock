package migrations

import (
	"database/sql"
	"embed"
	"log"
	"sort"
	"strings"
)

// Arquivos .up.sql embutidos no binário (pasta migrations).
//
//go:embed *.up.sql
var fs embed.FS

const schemaTable = "schema_migrations"

// Run aplica todas as migrations ainda não aplicadas no banco.
// Cria a tabela schema_migrations se não existir e executa cada arquivo .up.sql
// em ordem alfabética pelo nome do arquivo. Cada migration é executada em transação.
func Run(db *sql.DB) error {
	if err := createSchemaTable(db); err != nil {
		return err
	}

	applied, err := appliedVersions(db)
	if err != nil {
		return err
	}

	entries, err := fs.ReadDir(".")
	if err != nil {
		return err
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, name := range files {
		if applied[name] {
			continue
		}

		log.Printf("migration: applying %s", name)
		body, err := fs.ReadFile(name)
		if err != nil {
			return err
		}

		if err := runMigration(db, name, body); err != nil {
			return err
		}
		log.Printf("migration: applied %s", name)
	}

	return nil
}

func createSchemaTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + schemaTable + ` (
			version VARCHAR(255) NOT NULL PRIMARY KEY,
			applied_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`)
	return err
}

func appliedVersions(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT version FROM " + schemaTable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := make(map[string]bool)
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		m[v] = true
	}
	return m, rows.Err()
}

func runMigration(db *sql.DB, version string, body []byte) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	content := string(body)
	statements := splitStatements(content)

	for _, s := range statements {
		s = strings.TrimSpace(s)
		if s == "" || strings.HasPrefix(s, "--") {
			continue
		}
		if _, err := tx.Exec(s); err != nil {
			return err
		}
	}

	if _, err := tx.Exec(
		"INSERT INTO "+schemaTable+" (version) VALUES (?)",
		version,
	); err != nil {
		return err
	}

	return tx.Commit()
}

// splitStatements divide o SQL por ";" respeitando quebras de linha.
func splitStatements(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	var out []string
	for _, part := range strings.Split(s, ";") {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
