package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/starkandwayne/goutils/log"
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

func respond(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(response)
	if err != nil {
		log.Errorf("unable to encode response %s", "")
	}
}
