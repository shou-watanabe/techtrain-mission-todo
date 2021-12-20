package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	// set http handlers
	mux := http.NewServeMux()

	// TODO: ここから実装を行う
	hh := handler.NewHealthzHandler()
	mux.Handle("/healthz", middleware.AuthLayers(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hh.ServeHTTP(rw, r)
	})))

	ts := service.NewTODOService(todoDB)
	th := handler.NewTODOHandler(ts)
	mux.Handle("/todos", middleware.AuthLayers(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		th.ServeHTTP(rw, r)
	})))

	ph := handler.NewPanicHandler()
	mux.Handle("/do-panic", middleware.Layers(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ph.ServeHTTP(rw, r)
	})))

	if err := http.ListenAndServe(port, mux); err != nil {
		return err
	}

	return nil
}
