package config

import (
	"database/sql"

	"github.com/takumi3488/cookiejar-server/db"
	"github.com/takumi3488/cookiejar-server/internal/domain/repository"
	"github.com/takumi3488/cookiejar-server/internal/infrastructure/persistence"
	"github.com/takumi3488/cookiejar-server/internal/interface/handler"
	"github.com/takumi3488/cookiejar-server/internal/usecase"
)

type Container struct {
	// データベース
	DB      *sql.DB
	Queries *db.Queries

	// リポジトリ
	CookieRepo repository.CookieRepository

	// ユースケース
	CookieUsecase usecase.CookieUsecase

	// ハンドラー
	CookieHandler *handler.CookieHandler
}

func NewContainer(dbConn *sql.DB) *Container {
	queries := db.New(dbConn)

	// リポジトリを初期化
	cookieRepo := persistence.NewCookieRepository(queries)

	// ユースケースを初期化
	cookieUsecase := usecase.NewCookieUsecase(cookieRepo)

	// ハンドラーを初期化
	cookieHandler := handler.NewCookieHandler(cookieUsecase)

	return &Container{
		DB:      dbConn,
		Queries: queries,

		CookieRepo: cookieRepo,

		CookieUsecase: cookieUsecase,

		CookieHandler: cookieHandler,
	}
}
