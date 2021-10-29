package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ilyakaznacheev/cleanenv"
	pb "github.com/seregproj/fibonacci_slice/api/proto"
	"github.com/seregproj/fibonacci_slice/cmd/config"
	"github.com/seregproj/fibonacci_slice/internal/app"
	internallogger "github.com/seregproj/fibonacci_slice/internal/logger"
	grpcservice "github.com/seregproj/fibonacci_slice/internal/server/grpc"
	memorystorage "github.com/seregproj/fibonacci_slice/internal/storage/memory"
	redisstorage "github.com/seregproj/fibonacci_slice/internal/storage/redis"
	"google.golang.org/grpc"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "/etc/fibonacci/config.yml", "Path to configuration file")
	flag.Parse()

	conf := config.NewConfig()
	err := cleanenv.ReadConfig(configFile, conf)
	if err != nil {
		fmt.Printf("cant read config: %v with err %v", configFile, err.Error())
		os.Exit(1)
	}

	var storage app.Cache
	switch conf.Storage.Type {
	default:
		fmt.Printf("invalid storage type:  %v", conf.Storage.Type)
		os.Exit(1)
	case "memory":
		storage = memorystorage.NewMemoryCache()
	case "redis":
		storage = redisstorage.NewRedisCache(net.JoinHostPort(conf.Storage.Redis.Host, conf.Storage.Redis.Port),
			conf.Storage.Redis.DB, conf.Storage.Redis.Expires)
	}

	lvl, err := internallogger.NewLevel(conf.Logger.Level)
	if err != nil {
		fmt.Println(fmt.Errorf("cant create log level: %w", err))
		os.Exit(1)
	}

	logger, err := internallogger.New(conf.Logger.File, lvl)
	if err != nil {
		fmt.Println(fmt.Errorf("cant prepare log: %w", err))
		os.Exit(1)
	}

	fc := app.NewFibonacciCalc(storage, logger)
	grpcService := grpcservice.NewFibonacciService(fc)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	mux := runtime.NewServeMux()
	if err := pb.RegisterFibonacciServiceHandlerServer(ctx, mux, grpcService); err != nil {
		log.Fatalln(fmt.Errorf("cant register service handler: %w", err))
	}
	httpServer := http.Server{
		Addr:    net.JoinHostPort(conf.Server.HTTP.Host, conf.Server.HTTP.Port),
		Handler: mux,
	}
	go func() {
		defer cancel()

		if err = httpServer.ListenAndServe(); err != nil {
			log.Println(fmt.Errorf("cant start HTTP server: %w", err))
		}
	}()

	grpcServer := grpc.NewServer()
	go func() {
		defer cancel()

		pb.RegisterFibonacciServiceServer(grpcServer, grpcService)

		l, err := net.Listen("tcp", net.JoinHostPort(conf.Server.GRPC.Host, conf.Server.GRPC.Port))
		if err != nil {
			log.Println(fmt.Errorf("cant get grpc listener: %w", err))
		}

		err = grpcServer.Serve(l)
		if err != nil {
			log.Println(fmt.Errorf("cant start GRPC server: %w", err))
		}
	}()

	<-ctx.Done()
	grpcServer.GracefulStop()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalln(fmt.Errorf("cant shutdown http server: %w", err))
	}
}
