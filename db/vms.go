package db

import (
	"strings"
	"time"
)

type VMsFilter struct {
	Before *time.Time
	After  *time.Time
}

type VM struct {
	DeploymentName string `json:"deployment_name"`
	JobName        string `json:"job_name"`
	Index          int    `json:"index"`
	IP             string `json:"ip"`
}

func (f *VMsFilter) Query() (string, []interface{}) {
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
		SELECT DISTINCT source_deployment, source_job, source_index, source_ip
			FROM connections
			WHERE ` + strings.Join(wheres, " AND ") + `
			ORDER BY source_job ASC
	`, args
}

func (db *DB) GetVMs(filter *VMsFilter) ([]*VM, error) {
	if filter == nil {
		filter = &VMsFilter{}
	}

	l := []*VM{}
	query, args := filter.Query()
	r, err := db.Query(query, args...)
	if err != nil {
		return l, err
	}
	defer r.Close()

	for r.Next() {
		ann := &VM{}
		if err = r.Scan(&ann.DeploymentName, &ann.JobName, &ann.Index, &ann.IP); err != nil {
			return l, err
		}

		l = append(l, ann)
	}

	return l, nil
}
