package main

import (
	_ "embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

type backend struct {
	tasks     *TaskModel
	natsSub   string
	natscon   *nats.Conn
	logger    *slog.Logger
	broken    bool
	brokenMux sync.RWMutex
}

func main() {
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		fmt.Println("env PORT was unset\nUsing Port 3001 as Backend port")
		port = "3001"
	}

	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	natscon, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer natscon.Drain()

	natsSubject := os.Getenv("NATS_SUBJECT")
	if err != nil {
		log.Fatal("missing NATS_SUBJECT var")
	}

	err = initialseTable(db)
	if err != nil {
		log.Printf("Error initialseTable: %v", err)
		return
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &backend{
		tasks:   &TaskModel{DB: db},
		logger:  logger,
		natsSub: natsSubject,
		natscon: natscon,
	}

	addr := ":" + port
	logger.Info("Starting Todo-App Backend", "address", addr)
	log.Fatal(http.ListenAndServe(addr, app.routes()))
}
