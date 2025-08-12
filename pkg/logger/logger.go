package logger

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

var (
	development = "development"
	production  = "production"
)

// AsyncHandler оборачивает любой slog.Handler и делает его асинхронным
type AsyncHandler struct {
	handler slog.Handler
	ch      chan slog.Record
	wg      sync.WaitGroup
	closed  bool
	mu      sync.Mutex
}

func NewAsyncHandler(ctx context.Context, handler slog.Handler, bufferSize int) *AsyncHandler {
	if bufferSize <= 0 {
		bufferSize = 10000
	}

	ah := &AsyncHandler{
		handler: handler,
		ch:      make(chan slog.Record, bufferSize),
	}

	go ah.worker(ctx)

	return ah
}

func (ah *AsyncHandler) Handle(ctx context.Context, record slog.Record) error {
	ah.mu.Lock()
	if ah.closed {
		ah.mu.Unlock()
		return nil
	}
	ah.mu.Unlock()

	select {
	case ah.ch <- record:
	default:
		fmt.Println("log buffer full")
	}

	return nil
}

func (ah *AsyncHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return ah.handler.Enabled(ctx, level)
}

// WithAttrs проксирует вызов к базовому handler
func (ah *AsyncHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &AsyncHandler{
		handler: ah.handler.WithAttrs(attrs),
		ch:      ah.ch,
	}
}

// WithGroup проксирует вызов к базовому handler
func (ah *AsyncHandler) WithGroup(name string) slog.Handler {
	return &AsyncHandler{
		handler: ah.handler.WithGroup(name),
		ch:      ah.ch,
	}
}

// worker обрабатывает записи из канала
func (ah *AsyncHandler) worker(ctx context.Context) {
	defer ah.wg.Done()
	ah.wg.Add(1)

	for record := range ah.ch {

		ah.handler.Handle(ctx, record)
	}
}

// Close закрывает асинхронный handler и ждет завершения обработки
func (ah *AsyncHandler) Close() {
	ah.mu.Lock()
	if !ah.closed {
		ah.closed = true
		close(ah.ch)
	}
	ah.mu.Unlock()

	ah.wg.Wait()
}

// SetupLogger создает логгер с асинхронным handler
func SetupLogger(ctx context.Context, env string) (*slog.Logger, func()) {
	var baseHandler slog.Handler

	switch env {
	case production:
		baseHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	case development:
		baseHandler = setupPrettySlog().Handler()
	default:
		baseHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	asyncHandler := NewAsyncHandler(ctx, baseHandler, 500)

	return slog.New(asyncHandler), asyncHandler.Close
}

func setupPrettySlog() *slog.Logger {
	opts := PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func NewMiddlewareLogger(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
