package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	pb "github.com/takumi3488/cookiejar-server/gen/cookiejar/v1"
	"github.com/takumi3488/cookiejar-server/internal/config"
	"github.com/takumi3488/cookiejar-server/internal/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type cookieServiceServer struct {
	pb.UnimplementedCookieServiceServer
	container *config.Container
}

func (s *cookieServiceServer) GetCookies(ctx context.Context, req *pb.GetCookiesRequest) (*pb.GetCookiesResponse, error) {
	// hostでCookieを取得
	cookies, err := s.container.CookieUsecase.GetCookiesByHost(ctx, req.Host)
	if err != nil {
		log.Printf("Failed to get cookies for host %s: %v", req.Host, err)
		return nil, status.Errorf(codes.NotFound, "cookies not found for host: %s", req.Host)
	}

	// Cookieをhttp.Cookieに変換してからString形式に変換
	var cookieStrings []string
	for _, cookie := range cookies {
		httpCookie := cookie.ToHTTPCookie()
		cookieStrings = append(cookieStrings, httpCookie.String())
	}

	return &pb.GetCookiesResponse{
		Cookies: strings.Join(cookieStrings, "; "),
	}, nil
}

func main() {
	// OpenTelemetry の初期化
	tp, err := telemetry.InitTracer("cookiejar-reader")
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := telemetry.Shutdown(ctx, tp); err != nil {
			log.Printf("Failed to shutdown tracer: %v", err)
		}
	}()

	// 起動時のinfoスパンを送信
	func() {
		ctx := context.Background()
		tracer := otel.Tracer("cookiejar-reader")
		_, span := tracer.Start(ctx, "application.startup")
		defer span.End()

		port := os.Getenv("GRPC_PORT")
		if port == "" {
			port = "50051"
		}

		span.SetAttributes(
			attribute.String("service.name", "cookiejar-reader"),
			attribute.String("service.type", "grpc"),
			attribute.String("service.port", port),
		)
		span.SetStatus(otelcodes.Ok, "Application started successfully")
		log.Println("Startup info span sent")
	}()

	// データベース接続を初期化
	dbClient, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := dbClient.Close(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		}
	}()

	// 依存性注入コンテナを初期化
	container := config.NewContainer(dbClient)

	// gRPCサーバーを初期化（otelgrpc interceptorを追加）
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	pb.RegisterCookieServiceServer(grpcServer, &cookieServiceServer{
		container: container,
	})

	// gRPC reflectionを有効化（開発/テスト用）
	reflection.Register(grpcServer)

	// ポート50051でリスナーを作成
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("gRPC server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
