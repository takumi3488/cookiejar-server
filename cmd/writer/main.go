package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/takumi3488/cookiejar-server/internal/config"
	"github.com/takumi3488/cookiejar-server/internal/middleware"
	"github.com/takumi3488/cookiejar-server/internal/telemetry"

	_ "github.com/lib/pq"
)

func main() {
	// OpenTelemetry の初期化
	tp, err := telemetry.InitTracer("cookiejar-writer")
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

	// 新しいFiberアプリを初期化
	app := fiber.New()

	// OpenTelemetry middleware を追加
	app.Use(middleware.OpenTelemetry())

	// 環境変数からAllowOriginsを取得
	allowOrigins := strings.Split(os.Getenv("ALLOW_ORIGINS"), ",")

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           3600,
		ExposeHeaders:    []string{"Content-Length"},
	}))

	// ルートを登録
	app.Post("/", container.CookieHandler.StoreCookies)

	// ポート3000でサーバーを起動
	log.Fatal(app.Listen(":3000"))
}
