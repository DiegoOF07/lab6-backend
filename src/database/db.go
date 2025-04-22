package database

import (
    "database/sql"
    "log"
    "os"
    "strings"


    _ "github.com/mattn/go-sqlite3"
)

func SetupDatabase(dbPath string) (*sql.DB, error) {
    log.Printf("Intentando conectarse a : %s", dbPath)
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }

    if err = db.Ping(); err != nil {
        db.Close()
        return nil, err
    }

    ddlPath := os.Getenv("DDL_PATH")
    if ddlPath == "" {
        ddlPath = "./ddl.sql" // Valor por defecto
    }

    ddlBytes, err := os.ReadFile(ddlPath)
	if err != nil {
		log.Fatalf("Error al leer el archivo DDL: %v", err)
	}

	ddl := string(ddlBytes)
	statements := strings.Split(ddl, ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		_, err = db.Exec(stmt)
		if err != nil {
			log.Printf("Error al ejecutar: %s\nError: %v", stmt, err)
		}
	}
    log.Println("Conexión existosa")
    return db, nil
}

