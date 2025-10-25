package telemetry

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// InitTracer は OpenTelemetry の TracerProvider を初期化します
func InitTracer(serviceName string) (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	// OTLP gRPC exporter を作成（Jaeger用）
	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	isSecure := false
	if otlpEndpoint == "" {
		otlpEndpoint = "jaeger:4317" // デフォルト値（gRPC）
	} else {
		// 環境変数からスキームを取り除き、HTTPSかどうかを判定
		otlpEndpoint, isSecure = stripSchemeAndDetectSecure(otlpEndpoint)
	}

	// エクスポーターのオプションを構築
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(otlpEndpoint),
	}
	if !isSecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Resource を作成
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// TracerProvider を作成
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// グローバルに設定
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	scheme := "http"
	if isSecure {
		scheme = "https"
	}
	log.Printf("OpenTelemetry initialized for service: %s, endpoint: %s://%s", serviceName, scheme, otlpEndpoint)

	return tp, nil
}

// Shutdown は TracerProvider をシャットダウンします
func Shutdown(ctx context.Context, tp *sdktrace.TracerProvider) error {
	if tp == nil {
		return nil
	}
	return tp.Shutdown(ctx)
}

// stripSchemeAndDetectSecure は URL からスキームを取り除き、HTTPSかどうかを返します
func stripSchemeAndDetectSecure(endpoint string) (string, bool) {
	// https:// の場合は secure = true
	if stripped, ok := strings.CutPrefix(endpoint, "https://"); ok {
		return stripped, true
	}
	// http:// の場合は secure = false
	if stripped, ok := strings.CutPrefix(endpoint, "http://"); ok {
		return stripped, false
	}
	// スキームなしの場合はデフォルトで HTTP (secure = false)
	return endpoint, false
}
