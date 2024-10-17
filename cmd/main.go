package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/v1adhope/flights/internal/configs"
	v1 "github.com/v1adhope/flights/internal/controllers/http/v1"
	"github.com/v1adhope/flights/internal/usecases"
	"github.com/v1adhope/flights/internal/usecases/infrastructure/repository"
	"github.com/v1adhope/flights/pkg/httpsrv/httpsrv"
	"github.com/v1adhope/flights/pkg/logger"
	"github.com/v1adhope/flights/pkg/postgresql"
)

func main() {
	mainCtx := context.Background()

	configs.MustConfig()

	pd, err := postgresql.Build(
		mainCtx,
		postgresql.WithConnStr(configs.Global.Postgres.ConnStr),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer pd.Close()

	repo := repository.New(pd)

	usecases := usecases.New(repo)

	log := logger.New(
		logger.WithLevel("debug"),
	)

	v1.SetMode(configs.Global.Srv.Mode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	v1.Register(&v1.Router{
		Handler:  router,
		Usecases: usecases,
		Log:      log,
	})
	v1Srv := httpsrv.New(
		router,
		httpsrv.WithShutdownTimeout(configs.Global.Srv.ShutdownTimeout),
	)
	v1Srv.Run()
}
