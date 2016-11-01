package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/starkandwayne/covalence/db"
	"github.com/starkandwayne/goutils/log"
)

type ErrorResponse struct {
	Error       string `json:"error,omitempty"`
	Description string `json:"description"`
}

type MessageResponse struct {
	Message string `json:"message,omitempty"`
}

type ConnectionHandler struct {
	Data *db.DB
}

func (c ConnectionHandler) create(w http.ResponseWriter, req *http.Request) {
	var connections []db.Connection
	if err := json.NewDecoder(req.Body).Decode(&connections); err != nil {
		c.respond(w, http.StatusUnprocessableEntity, ErrorResponse{
			Description: err.Error(),
		})
		fmt.Println(err)
		return
	}
	for _, connection := range connections {
		_, err := c.Data.CreateConnection(connection.Source.IP, connection.Source.Port, connection.Source.Deployment, connection.Source.Job,
			connection.Source.Index, connection.Source.User, connection.Source.Group, connection.Source.Pid, connection.Source.ProcessName,
			connection.Source.Age, connection.Destination.IP, connection.Destination.Port)
		if err != nil {
			log.Errorf("unable to create connection: %s", err.Error())
			c.respond(w, http.StatusInternalServerError, ErrorResponse{
				Description: err.Error(),
			})
			return
		}
	}
	c.respond(w, http.StatusOK, MessageResponse{
		Message: "Resources created.",
	})
	return
}

func (c ConnectionHandler) get(w http.ResponseWriter, req *http.Request) {
	connections, err := c.Data.GetAllConnections(&db.ConnectionFilter{
		Before: paramUnixTime(req, "before"),
		After:  paramUnixTime(req, "after"),
	})
	if err != nil {
		c.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	c.respond(w, http.StatusOK, connections)
}

func (c ConnectionHandler) respond(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(response)
	if err != nil {
		log.Errorf("unable to encode response %s", "")
	}
}
