package middleware

import (
	"github.com/gofiber/fiber/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "cookiejar-server/fiber"
)

// OpenTelemetry は Fiber v3 用の OpenTelemetry middleware を返します
func OpenTelemetry() fiber.Handler {
	tracer := otel.Tracer(tracerName)

	return func(c fiber.Ctx) error {
		// コンテキストから trace context を抽出
		ctx := otel.GetTextMapPropagator().Extract(
			c.Context(),
			propagation.HeaderCarrier(c.GetReqHeaders()),
		)

		// span を開始
		spanName := c.Method() + " " + c.Route().Path
		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPMethod(c.Method()),
				semconv.HTTPRoute(c.Route().Path),
				semconv.HTTPTarget(c.OriginalURL()),
				semconv.HTTPScheme(c.Protocol()),
				semconv.NetHostName(c.Hostname()),
			),
		)
		defer span.End()

		// コンテキストを設定
		c.SetContext(ctx)

		// 次のハンドラを実行
		err := c.Next()

		// レスポンス属性を設定
		span.SetAttributes(
			semconv.HTTPStatusCode(c.Response().StatusCode()),
		)

		// エラーがある場合は記録
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else if c.Response().StatusCode() >= 400 {
			span.SetStatus(codes.Error, "HTTP error")
		} else {
			span.SetStatus(codes.Ok, "")
		}

		return err
	}
}
