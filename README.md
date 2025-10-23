# cookiejar-server

## 概要

CookieをPostgreSQLデータベースに保存・取得するためのマイクロサービスアプリケーション

## サービス構成

このアプリケーションは2つのマイクロサービスで構成されています：

- **Writer**: Cookie情報を保存するHTTP REST APIサーバー（ポート3000）
- **Reader**: Cookie情報を取得するgRPCサーバー（ポート50051）

## 機能

### Writer
- Cookie情報の保存（Upsert）

### Reader
- ホスト名によるCookie情報の取得

## アーキテクチャ

Clean Architectureに基づいた設計：

```
.
├── cmd/
│   ├── writer/main.go               # Writer エントリーポイント
│   └── reader/main.go               # Reader エントリーポイント
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
├── proto/v1/                        # gRPC protoファイル
├── gen/v1/                          # gRPC生成コード
├── db/                              # SQLC生成コード
├── queries/                         # SQLクエリ定義
├── e2e/                             # E2Eテスト（runn）
└── schema.sql                       # データベーススキーマ
```

## セットアップ

### Docker Composeを使用する場合（推奨）

```bash
# すべてのサービスを起動
docker compose up -d

# ログを確認
docker compose logs -f

# サービスを停止
docker compose down
```

これにより、以下のサービスが起動します：
- PostgreSQL（ポート5432）
- Writer（ポート3000）
- Reader（ポート50051）
- Jaeger（ポート16686）

### ローカルで実行する場合

#### 必要な環境変数

```bash
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DB=cookiejar
ALLOW_ORIGINS=http://localhost:3000
GRPC_PORT=50051
```

#### データベースの初期化

```bash
psql -U postgres -d cookiejar -f schema.sql
```

#### Writer のビルドと実行

```bash
go build -o cookiejar-writer ./cmd/writer
./cookiejar-writer
```

Writerはポート3000で起動します。

#### Reader のビルドと実行

```bash
go build -o cookiejar-reader ./cmd/reader
./cookiejar-reader
```

Readerはポート50051で起動します。

## API

### Writer API (HTTP REST)

#### POST /

Cookie情報を保存します。

**エンドポイント:** `http://localhost:3000/`

**リクエストボディ:**
```json
[
  {
    "name": "cookie_name",
    "value": "cookie_value",
    "domain": "example.com",
    "path": "/",
    "maxAge": 3600,
    "secure": true,
    "httpOnly": true,
    "sameSite": "Lax"
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

### Reader API (gRPC)

#### GetCookies

ホスト名でCookie情報を取得します。

**エンドポイント:** `localhost:50051`

**リクエスト:**
```protobuf
message GetCookiesRequest {
  string host = 1;
}
```

**レスポンス:**
```protobuf
message GetCookiesResponse {
  string cookies = 1;  // Cookie文字列 (例: "name1=value1; name2=value2")
}
```

## E2Eテスト

[runn](https://github.com/k1LoW/runn)を使用したE2Eテストを提供しています。

### runnのインストール

```bash
# Homebrewを使用
brew install k1LoW/tap/runn

# Goを使用
go install github.com/k1LoW/runn/cmd/runn@latest
```

### テストの実行

サービスが起動している状態で実行してください：

```bash
# すべてのサービスを起動
docker compose up -d

# Writer APIのテスト
runn run e2e/writer.yml

# Reader APIのテスト（gRPC）
runn run e2e/reader.yml

# 統合テスト（Writer → Reader）
runn run e2e/integration.yml

# すべてのテストを実行
runn run e2e/*.yml
```

### テスト内容

- **writer.yml**: Writer APIのテスト（Cookie保存、バリデーション、エラーハンドリング）
- **reader.yml**: Reader APIのテスト（gRPCでのCookie取得）
- **integration.yml**: 統合テスト（WriterでCookieを保存してReaderで取得）

## 技術スタック

- Go 1.25+
- [Fiber v3](https://github.com/gofiber/fiber) - Webフレームワーク
- [gRPC](https://grpc.io/) - RPC framework
- [SQLC](https://github.com/sqlc-dev/sqlc) - 型安全なSQLコード生成
- PostgreSQL - データベース
- [runn](https://github.com/k1LoW/runn) - E2Eテストツール
