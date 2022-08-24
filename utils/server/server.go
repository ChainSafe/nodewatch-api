// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

// Package server configures the api server
package server

import (
	"context"
	"eth2-crawler/utils/config"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/cors"
)

// Start starts the service
func Start(ctx context.Context, cfg *config.Server, handler http.Handler) {
	cors := cors.New(cors.Options{
		AllowedOrigins:   cfg.CORS,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD"},
		ExposedHeaders:   []string{"Content-Length", "Content-Type", "Content-Disposition"},
		AllowCredentials: true,
	})

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		ReadTimeout:       time.Duration(cfg.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(cfg.ReadHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(cfg.WriteTimeout) * time.Second,
		Handler:           cors.Handler(handler),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			// log error
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, ecancel := context.WithTimeout(ctx, 10*time.Second)
	defer ecancel()
	if err := server.Shutdown(ctx); err != nil {
		// log error
	}
}
