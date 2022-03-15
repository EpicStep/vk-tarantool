package main

import (
	"context"
	"errors"
	vk_tarantool "github.com/EpicStep/vk-tarantool"
	"github.com/EpicStep/vk-tarantool/internal/router"
	"github.com/EpicStep/vk-tarantool/internal/shorter"
	"github.com/EpicStep/vk-tarantool/pkg/database"
	"github.com/EpicStep/vk-tarantool/pkg/server"
	"github.com/EpicStep/vk-tarantool/pkg/tracing"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	log.Println("Service started")

	tp, err := tracing.TracerProvider("")
	if err != nil {
		return err
	}

	otel.SetTracerProvider(tp)

	tp.Tracer("vk-tarantool")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	db, err := database.New("127.0.0.1:3301", "guest", "")
	if err != nil {
		return err
	}

	r := router.New()

	service := shorter.New(db, tp.Tracer(""))

	front, err := vk_tarantool.GetFrontendAssets()
	if err != nil {
		return err
	}

	r.Route("/", func(r chi.Router) {
		r.Handle("/metrics", promhttp.Handler())
		r.Handle("/ui/*", http.FileServer(http.FS(front)))

		r.Route("/api", func(r chi.Router) {
			service.Routes(r)
		})
	})

	addr := ":80"

	srv := server.New(addr, r)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdown()

	if err := srv.Shutdown(ctx); err != nil {
		return errors.New("server shutdown failed")
	}

	if err := db.Close(); err != nil {
		return errors.New("database close failed")
	}

	return nil
}
