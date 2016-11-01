package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/starkandwayne/covalence/db"
	"github.com/starkandwayne/goutils/log"
)

type WebServer struct {
	Database *db.DB
	Addr     string
	WebRoot  string
	Api      *Api
}

func (ws *WebServer) Setup() error {
	log.Debugf("Configuring WebServer...")
	if err := ws.Database.Connect(); err != nil {
		log.Errorf("Failed to connect to %s database at %s: %s", ws.Database.Driver, ws.Database.DSN, err)
		return err
	}

	connectionHandler := ConnectionHandler{Data: ws.Database}

	r := mux.NewRouter()
	r.HandleFunc("/connections", connectionHandler.create).Methods("POST")
	r.HandleFunc("/connections", connectionHandler.get).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(ws.WebRoot)))
	http.Handle("/", r)
	return nil
}

func (ws *WebServer) Start() {
	err := ws.Setup()
	if err != nil {
		panic("Could not set up WebServer for Covalence: " + err.Error())
	}
	log.Debugf("Starting WebServer on '%s'...", ws.Addr)
	err = http.ListenAndServe(ws.Addr, nil)
	if err != nil {
		log.Errorf("HTTP API failed %s", err.Error())
		panic("Cannot setup WebServer, aborting. Check logs for details.")
	}
}
