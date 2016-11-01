package api

import (
	"fmt"

	"github.com/starkandwayne/covalence/db"
)

type Api struct {
	resync chan int /* api goroutine will send here when the db changes significantly (i.e. new job, updated target, etc.) */

	Database *db.DB

	Web *WebServer /* Webserver that gets spawned to handle http requests */

}

var Version = "(development)"

func NewApi() *Api {
	return &Api{
		resync:   make(chan int),
		Database: &db.DB{},
	}
}

func (a *Api) Run() error {
	if err := a.Database.Connect(); err != nil {
		return fmt.Errorf("failed to connect to %s database at %s: %s\n",
			a.Database.Driver, a.Database.DSN, err)
	}

	a.Web.Start()
	return nil
}
