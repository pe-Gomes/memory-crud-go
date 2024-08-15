package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pe-Gomes/memory-crud-go/api"
	"github.com/pe-Gomes/memory-crud-go/infra"
)

func main() {
	h := api.NewHandler(infra.NewAppDB())

	s := http.Server{
		Addr:           ":8080",
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    time.Minute,
		MaxHeaderBytes: 0,
	}

	if err := s.ListenAndServe(); err != nil {
		slog.Error("failed to serve", "err", err)
	}

	id := uuid.New()
	id.String()
}
