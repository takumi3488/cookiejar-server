package telemetry

import (
	"context"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InitTracer はOTLPエクスポーターとトレーサープロバイダーを初期化します
func InitTracer(serviceName string) (func(context.Context) error, error) {
	ctx := context.Background()

	// OTLP エンドポイントを環境変数から取得（デフォルト: jaeger:4317）
	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "jaeger:4317"
	}

	// OTLP gRPC エクスポーターを作成
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(otlpEndpoint),
		otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		otlptracegrpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}

	// リソースを作成（サービス名を設定）
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	// トレーサープロバイダーを作成
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// グローバルトレーサープロバイダーを設定
	otel.SetTracerProvider(tp)

	log.Printf("OpenTelemetry tracer initialized for service: %s, endpoint: %s", serviceName, otlpEndpoint)

	// シャットダウン関数を返す
	return tp.Shutdown, nil
}
