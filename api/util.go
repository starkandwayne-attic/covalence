package api

import (
	"net/http"
	"strconv"
	"time"
)

func paramUnixTime(req *http.Request, name string) *time.Time {
	//utc, _ := time.LoadLocation("UTC")

	value, set := req.URL.Query()[name]
	if !set {
		return nil
	}
	i, err := strconv.ParseInt(value[0], 10, 64)
	if err != nil {
		return nil
	}
	t := time.Unix(i, 0)
	return &t
}
