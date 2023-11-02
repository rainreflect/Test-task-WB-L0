package application

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rec/internal/config"
	"rec/internal/pkg/cache"
	"rec/internal/pkg/handler.go"
	"rec/internal/pkg/postgresql"
	"rec/internal/pkg/postgresql/caching"
	"rec/internal/pkg/server"
	"syscall"
	"time"
)

type App struct {
	cfg config.Config
}

func NewApp(cfg config.Config) *App {
	app := App{}
	app.cfg = cfg
	return &app
}

func (app *App) Run() {
	db, err := postgresql.Conn(app.cfg)
	if err != nil {
		log.Println("Error while connnecting to database")
		panic(err)
	}

	cacheServ := cache.NewCacheServ()
	storeService := caching.NewCacheService(*cacheServ, db)
	err = storeService.RestoreCache()
	if err != nil {
		log.Println("Cache for restoring wasn't found")
	}

	sh := handler.NewStreamingHandler(db)

	server := server.NewServer(*storeService, app.cfg.Http.Host+":"+app.cfg.Http.Port, db)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
		sh.Finish()
		defer server.Stop()
		defer db.Close()
	}()

	if err := server.Serv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown Failed:", err)
	}
	log.Print("Server EXIT")
}
