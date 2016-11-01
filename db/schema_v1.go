package db

import (
	"fmt"
)

type v1Schema struct{}

func (s v1Schema) Deploy(db *DB) error {
	err := db.Exec(`CREATE TABLE schema_info (
               version INTEGER
             )`)
	if err != nil {
		return err
	}

	err = db.Exec(`INSERT INTO schema_info VALUES (1)`)
	if err != nil {
		return err
	}

	switch db.Driver {
	case "mysql":
		err = db.Exec(`CREATE TABLE connections (
               uuid      VARCHAR(36) NOT NULL,
               source_ip            TEXT,
               source_port          TEXT,
               source_deployment    TEXT,
               source_job           TEXT,
							 source_index         INTEGER,
               source_user          TEXT,
               source_group         TEXT,
               source_pid           TEXT,
							 source_process_name  TEXT,
               source_age           INTEGER,
							 destination_ip       TEXT,
               destination_port     TEXT,
							 created_at           INTEGER,
               PRIMARY KEY (uuid)
             )`)
	case "postgres", "sqlite3":
		err = db.Exec(`CREATE TABLE connections (
               uuid      UUID PRIMARY KEY,
							 source_ip            TEXT,
               source_port          TEXT,
               source_deployment    TEXT,
               source_job           TEXT,
							 source_index         INTEGER,
               source_user          TEXT,
               source_group         TEXT,
               source_pid           TEXT,
							 source_process_name  TEXT,
               source_age           INTEGER,
							 destination_ip       TEXT,
               destination_port     TEXT,
							 created_at           INTEGER
             )`)
	default:
		err = fmt.Errorf("unsupported database driver '%s'", db.Driver)
	}
	if err != nil {
		return err
	}

	return nil
}
