package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	pbProducts "github.com/ilkinabd/goods-contracts/gen/go/products/v1"
	_ "github.com/ilkinabd/goods-manager/app/docs"
	"github.com/ilkinabd/goods-manager/app/internal/config"
	product "github.com/ilkinabd/goods-manager/app/internal/controller/grpc/v1/product"
	"github.com/ilkinabd/goods-manager/app/internal/domain/product/dao"
	"github.com/ilkinabd/goods-manager/app/internal/domain/product/policy"
	"github.com/ilkinabd/goods-manager/app/internal/domain/product/service"
	"github.com/ilkinabd/goods-manager/app/pkg/client/postgresql"
	"github.com/ilkinabd/goods-manager/app/pkg/logging"
	"github.com/ilkinabd/goods-manager/app/pkg/metric"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	cfg *config.Config

	router     *httprouter.Router
	httpServer *http.Server
	grpcServer *grpc.Server

	pgClient *pgxpool.Pool

	productServiceServer pbProducts.ProductServiceServer
}

func NewApp(ctx context.Context, config *config.Config) (App, error) {
	logging.Info(ctx, "router initializing")
	router := httprouter.New()

	logging.Info(ctx, "swagger docs initializing")
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	logging.Info(ctx, "heartbeat metric initializing")
	metricHandler := metric.Handler{}
	metricHandler.Register(router)

	pgConfig := postgresql.NewPgConfig(
		config.PostgreSQL.Username, config.PostgreSQL.Password,
		config.PostgreSQL.Host, config.PostgreSQL.Port, config.PostgreSQL.Database,
	)
	pgClient, err := postgresql.NewClient(ctx, 5, time.Second*5, pgConfig)
	if err != nil {
		logging.GetLogger().Fatal(ctx, err)
	}

	productDao := dao.NewProductDAOPostgres(pgClient)
	productService := service.NewProductService(productDao)
	productPolicy := policy.NewProductPolicy(productService)
	productServiceServer := product.NewServer(
		productPolicy,
		pbProducts.UnimplementedProductServiceServer{},
	)

	return App{
		cfg:                  config,
		router:               router,
		pgClient:             pgClient,
		productServiceServer: productServiceServer,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return a.startHTTP(ctx)
	})
	grp.Go(func() error {
		return a.startGRPC(ctx, a.productServiceServer)
	})
	return grp.Wait()
}

func (a *App) startGRPC(ctx context.Context, server pbProducts.ProductServiceServer) error {
	logger := logging.WithFields(ctx, map[string]interface{}{
		"IP":   a.cfg.GRPC.IP,
		"Port": a.cfg.GRPC.Port,
	})
	logger.Info("gRPC Server initializing")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.GRPC.IP, a.cfg.GRPC.Port))
	if err != nil {
		logger.WithError(err).Fatal("failed to create listener")
	}

	serverOptions := []grpc.ServerOption{}

	a.grpcServer = grpc.NewServer(serverOptions...)

	pbProducts.RegisterProductServiceServer(a.grpcServer, server)

	reflection.Register(a.grpcServer)

	return a.grpcServer.Serve(listener)
}

func (a *App) startHTTP(ctx context.Context) error {
	logger := logging.WithFields(ctx, map[string]interface{}{
		"IP":   a.cfg.HTTP.IP,
		"Port": a.cfg.HTTP.Port,
	})
	logger.Info("HTTP Server initializing")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.HTTP.IP, a.cfg.HTTP.Port))
	if err != nil {
		logger.WithError(err).Fatal("failed to create listener")
	}

	logger.WithFields(map[string]interface{}{
		"AllowedMethods":     a.cfg.HTTP.CORS.AllowedMethods,
		"AllowedOrigins":     a.cfg.HTTP.CORS.AllowedOrigins,
		"AllowCredentials":   a.cfg.HTTP.CORS.AllowCredentials,
		"AllowedHeaders":     a.cfg.HTTP.CORS.AllowedHeaders,
		"OptionsPassthrough": a.cfg.HTTP.CORS.OptionsPassthrough,
		"ExposedHeaders":     a.cfg.HTTP.CORS.ExposedHeaders,
		"Debug":              a.cfg.HTTP.CORS.Debug,
	})
	c := cors.New(cors.Options{
		AllowedMethods:     a.cfg.HTTP.CORS.AllowedMethods,
		AllowedOrigins:     a.cfg.HTTP.CORS.AllowedOrigins,
		AllowCredentials:   a.cfg.HTTP.CORS.AllowCredentials,
		AllowedHeaders:     a.cfg.HTTP.CORS.AllowedHeaders,
		OptionsPassthrough: a.cfg.HTTP.CORS.OptionsPassthrough,
		ExposedHeaders:     a.cfg.HTTP.CORS.ExposedHeaders,
		Debug:              a.cfg.HTTP.CORS.Debug,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: a.cfg.HTTP.WriteTimeout,
		ReadTimeout:  a.cfg.HTTP.ReadTimeout,
	}

	if err = a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Warning("server shutdown")
		default:
			logger.Fatal(err)
		}
	}
	err = a.httpServer.Shutdown(context.Background())
	if err != nil {
		logger.Fatal(err)
	}
	return err
}
