package api

import (
	"net/http"

	"github.com/starkandwayne/covalence/db"
)

type VMHandler struct {
	Data *db.DB
}

func (v VMHandler) get(w http.ResponseWriter, req *http.Request) {
	vms, err := v.Data.GetVMs(&db.VMsFilter{
		Before: paramUnixTime(req, "before"),
		After:  paramUnixTime(req, "after"),
	})
	if err != nil {
		respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	respond(w, http.StatusOK, vms)
}
