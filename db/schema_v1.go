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
	case "sqlite3":
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
	case "postgres":
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
		err = db.Exec(`CREATE OR REPLACE FUNCTION
									 public.connections_partition_function()
									 RETURNS TRIGGER AS
									 $BODY$
									 DECLARE
									 _new_time int;
									 _tablename text;
									 _startdate text;
									 _enddate text;
									 _result record;
									 BEGIN
									 --Takes the current inbound "created_at" value and determines when midnight is for the given date
									 _new_time := ((NEW."created_at"/86400)::int)*86400;
									 _startdate := to_char(to_timestamp(_new_time), 'YYYY-MM-DD');
									 _tablename := 'connections_'||_startdate;

									 -- Check if the partition needed for the current record exists
									 PERFORM 1
									 FROM   pg_catalog.pg_class c
									 JOIN   pg_catalog.pg_namespace n ON n.oid = c.relnamespace
									 WHERE  c.relkind = 'r'
									 AND    c.relname = _tablename
									 AND    n.nspname = 'public';

									 -- If the partition needed does not yet exist, then we create it:
									 -- Note that || is string concatenation (joining two strings to make one)
									 IF NOT FOUND THEN
									 _enddate:=_startdate::timestamp + INTERVAL '1 day';
									 EXECUTE 'CREATE TABLE public.' || quote_ident(_tablename) || ' (
										 CHECK ( "created_at" >= EXTRACT(EPOCH FROM DATE ' || quote_literal(_startdate) || ')
										 AND "created_at" < EXTRACT(EPOCH FROM DATE ' || quote_literal(_enddate) || ')
									 )
									 ) INHERITS (public.connections)';

									 -- Table permissions are not inherited from the parent.
									 -- If permissions change on the master be sure to change them on the child also.
									 -- EXECUTE 'ALTER TABLE public.' || quote_ident(_tablename) || ' OWNER TO postgres';
									 -- EXECUTE 'GRANT ALL ON TABLE public.' || quote_ident(_tablename) || ' TO postgres';

									 -- Indexes are defined per child, so we assign a default index that uses the partition columns
									 EXECUTE 'CREATE INDEX ' || quote_ident(_tablename||'_indx1') || ' ON public.' || quote_ident(_tablename) || ' (created_at, uuid)';
									 END IF;

									 -- Insert the current record into the correct partition, which we are sure will now exist.
									 EXECUTE 'INSERT INTO public.' || quote_ident(_tablename) || ' VALUES ($1.*)' USING NEW;
									 RETURN NULL;
									 END;
									 $BODY$
									 LANGUAGE plpgsql;
								`)

		err = db.Exec(`CREATE TRIGGER connections_trigger
						 			 BEFORE INSERT ON public.connections
						 			 FOR EACH ROW EXECUTE PROCEDURE public.connections_partition_function();
								`)

	default:
		err = fmt.Errorf("unsupported database driver '%s'", db.Driver)
	}
	if err != nil {
		return err
	}

	return nil
}
