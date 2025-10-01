# cookiejar-server

MoneyForward用のCookie保存サーバー

## 概要

このサーバーは、ブラウザ拡張機能から送信されたCookieをPostgreSQLデータベースに保存するためのREST APIサーバーです。
元々は [mf-go](https://github.com/takumi3488/mf-go) リポジトリの `cmd/jar` として存在していましたが、独立したリポジトリとして分離されました。

## 機能

- Cookie情報の保存（Upsert）
- 保存されたCookieの取得
- CORS対応（MoneyForward向け）

## アーキテクチャ

Clean Architectureに基づいた設計：

```
.
├── main.go                          # エントリーポイント
├── internal/
│   ├── config/                      # 依存性注入コンテナ
│   ├── domain/
│   │   ├── entity/                  # ドメインエンティティ
│   │   └── repository/              # リポジトリインターフェース
│   ├── infrastructure/
│   │   └── persistence/             # データベース実装
│   ├── interface/
│   │   └── handler/                 # HTTPハンドラー
│   └── usecase/                     # ビジネスロジック
├── db/                              # SQLC生成コード
├── queries/                         # SQLクエリ定義
└── schema.sql                       # データベーススキーマ
```

## セットアップ

### 必要な環境変数

```bash
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DB=cookiejar
```

### データベースの初期化

```bash
psql -U postgres -d cookiejar -f schema.sql
```

### ビルドと実行

```bash
go build -o cookiejar-server .
./cookiejar-server
```

サーバーはポート3000で起動します。

## API

### POST /

Cookie情報を保存します。

**リクエストボディ:**
```json
[
  {
    "name": "cookie_name",
    "value": "cookie_value",
    "domain": ".moneyforward.com",
    "path": "/",
    "secure": true,
    "httpOnly": true,
    "sameSite": "None"
  }
]
```

**レスポンス:**
```json
{
  "status": "success",
  "count": 1
}
```

## 技術スタック

- Go 1.25+
- [Fiber v3](https://github.com/gofiber/fiber) - Webフレームワーク
- [SQLC](https://github.com/sqlc-dev/sqlc) - 型安全なSQLコード生成
- PostgreSQL - データベース

## ライセンス

元のプロジェクト [mf-go](https://github.com/takumi3488/mf-go) から分離
