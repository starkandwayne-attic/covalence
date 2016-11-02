package db

import (
	"strings"
	"time"

	"github.com/pborman/uuid"
	. "github.com/starkandwayne/goutils/timestamp"
)

type Connection struct {
	UUID        uuid.UUID   `json:"uuid"`
	Source      Source      `json:"source"`
	Destination Destination `json:"destination"`
	CreatedAt   Timestamp   `json:"created_at"`
}

type Source struct {
	IP          string `json:"ip"`
	Port        string `json:"port"`
	Deployment  string `json:"deployment"`
	Job         string `json:"job"`
	Index       int    `json:"index"`
	User        string `json:"user"`
	Group       string `json:"group"`
	Pid         string `json:"pid"`
	ProcessName string `json:"process_name"`
	Age         int    `json:"age"`
}

type Destination struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type ConnectionFilter struct {
	Before *time.Time
	After  *time.Time
}

func ValidateEffectiveUnix(effective time.Time) int64 {
	if effective.Unix() <= 0 {
		return time.Now().Unix()
	}
	return effective.Unix()
}

func (f *ConnectionFilter) Query() (string, []interface{}) {
	var wheres []string = []string{"uuid = uuid"}
	var args []interface{}
	if f.Before != nil {
		wheres = append(wheres, "created_at <= ?")
		args = append(args, f.Before.Unix())
	}
	if f.After != nil {
		wheres = append(wheres, "created_at >= ?")
		args = append(args, f.After.Unix())
	}
	return `
		SELECT DISTINCT uuid, source_ip, source_port, source_deployment,
		                source_job, source_index, source_user, source_group,
										source_pid, source_process_name, source_age, destination_ip,
										destination_port, created_at
			FROM connections
			WHERE ` + strings.Join(wheres, " AND ") + `
			GROUP BY uuid
			ORDER BY uuid ASC
	`, args
}

func (db *DB) GetAllConnections(filter *ConnectionFilter) ([]*Connection, error) {
	if filter == nil {
		filter = &ConnectionFilter{}
	}

	l := []*Connection{}
	query, args := filter.Query()
	r, err := db.Query(query, args...)
	if err != nil {
		return l, err
	}
	defer r.Close()

	for r.Next() {
		ann := &Connection{}
		var this NullUUID
		var createdAt *int64
		if err = r.Scan(&this, &ann.Source.IP, &ann.Source.Port, &ann.Source.Deployment, &ann.Source.Job,
			&ann.Source.Index, &ann.Source.User, &ann.Source.Group, &ann.Source.Pid, &ann.Source.ProcessName,
			&ann.Source.Age, &ann.Destination.IP, &ann.Destination.Port, &createdAt); err != nil {
			return l, err
		}
		ann.UUID = this.UUID
		if createdAt != nil {
			ann.CreatedAt = parseEpochTime(*createdAt)
		}

		l = append(l, ann)
	}

	return l, nil
}

func (db *DB) GetConnection(id uuid.UUID) (*Connection, error) {
	r, err := db.Query(`
		SELECT uuid, source_ip, source_port, source_deployment,
		  source_job, source_index, source_user, source_group,
			source_pid, source_process_name, source_age, destination_ip,
			destination_port
			FROM connections WHERE uuid = ?`, id.String())
	if err != nil {
		return nil, err
	}
	defer r.Close()

	if !r.Next() {
		return nil, nil
	}

	ann := &Connection{}
	var this NullUUID
	if err = r.Scan(&this, &ann.Source.IP, &ann.Source.Port, &ann.Source.Deployment, &ann.Source.Job,
		&ann.Source.Index, &ann.Source.User, &ann.Source.Group, &ann.Source.Pid, &ann.Source.ProcessName,
		&ann.Source.Age, &ann.Destination.IP, &ann.Destination.Port); err != nil {
		return nil, err
	}
	ann.UUID = this.UUID

	return ann, nil
}

func (db *DB) CreateConnection(sourceIP, sourcePort, sourceDeployment, sourceJob string, sourceIndex int,
	sourceUser, sourceGroup string, sourcePid, sourceProcessName string, sourceAge int,
	destinationIP, destinationPort string) (uuid.UUID, error) {
	id := uuid.NewRandom()
	return id, db.Exec(
		`INSERT INTO connections (uuid, source_ip, source_port, source_deployment,
		  source_job, source_index, source_user, source_group,
			source_pid, source_process_name, source_age, destination_ip,
			destination_port, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id.String(), sourceIP, sourcePort, sourceDeployment, sourceJob, sourceIndex,
		sourceUser, sourceGroup, sourcePid, sourceProcessName, sourceAge, destinationIP,
		destinationPort, time.Now().Unix())
}

func (db *DB) DeleteConnection(id uuid.UUID) (bool, error) {
	return true, db.Exec(
		`DELETE FROM connections WHERE uuid = ?`,
		id.String(),
	)
}
