package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"go-gin-webapi/internal/auth"
	"go-gin-webapi/internal/config"
	"go-gin-webapi/internal/handler"
	"go-gin-webapi/internal/repo"
	"go-gin-webapi/schemas"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("mysql", cfg.DB.DSN())
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	idtk := auth.NewIdentityToolkitClient(cfg.Firebase.APIKey)
	fbAdmin, fbAdminErr := auth.NewFirebaseAdmin(ctx, cfg.Firebase)

	repos := repo.New(db)
	h := handler.NewAPI(repos, idtk, fbAdmin, fbAdminErr, cfg.Auth)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	schemas.RegisterHandlersWithOptions(r, h, schemas.GinServerOptions{
		BaseURL: "/api/v1",
	})

	// Optional swagger spec endpoint
	r.GET("/swagger.json", func(c *gin.Context) {
		spec, err := schemas.GetSwagger()
		if err != nil {
			msg := err.Error()
			c.JSON(http.StatusInternalServerError, schemas.InternalServerErrorJSONResponse{Error: &msg})
			return
		}
		c.JSON(http.StatusOK, spec)
	})

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
}


