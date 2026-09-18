package logger

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	gormlogger "gorm.io/gorm/logger"
)

func testLogger() (*Logger, *bytes.Buffer) {
	var output bytes.Buffer
	return &Logger{core: zerolog.New(&output)}, &output
}

func TestLoggerRedactsSensitiveFieldsAndValues(t *testing.T) {
	log, output := testLogger()
	log.Info("provider request",
		"webhook", "https://hooks.example.test/secret-path",
		"url", "https://hooks.example.test/secret-path",
		"token", "sk-live-secret",
		"message_body", "a customer phone number",
		"error", errors.New("request failed: https://user:password@example.test/hook?token=secret"),
		"status", "ok",
	)

	logged := output.String()
	for _, secret := range []string{"hooks.example.test/secret-path", "user:password", "sk-live-secret", "a customer phone number", "?token=secret"} {
		if strings.Contains(logged, secret) {
			t.Fatalf("log contains secret %q: %s", secret, logged)
		}
	}
	for _, expected := range []string{"[REDACTED]", "\"status\":\"ok\""} {
		if !strings.Contains(logged, expected) {
			t.Fatalf("log does not contain %q: %s", expected, logged)
		}
	}
}

func TestGormLoggerRedactsSQLLiteralValues(t *testing.T) {
	log, output := testLogger()
	gormLog := log.GetGormLogger(gormlogger.Info)
	gormLog.Trace(
		context.Background(),
		time.Now(),
		func() (string, int64) {
			return "SELECT * FROM users WHERE email = 'customer@example.test' AND token = 'secret' AND id = 42", 1
		},
		nil,
	)

	logged := output.String()
	for _, secret := range []string{"customer@example.test", "'secret'", "42"} {
		if strings.Contains(logged, secret) {
			t.Fatalf("SQL log contains value %q: %s", secret, logged)
		}
	}
	if !strings.Contains(logged, "SQL query") || !strings.Contains(logged, "[REDACTED]") {
		t.Fatalf("SQL log was not emitted in sanitized form: %s", logged)
	}
}

func TestRedactTextSanitizesStructuredSecrets(t *testing.T) {
	got := redactText(`provider returned password="secret" and dsn=postgres://user:pass@db/app`)
	if strings.Contains(got, "secret") || strings.Contains(got, "user:pass") {
		t.Fatalf("redactText() leaked secret: %s", got)
	}
	if !strings.Contains(got, "[REDACTED]") {
		t.Fatalf("redactText() = %q, want redaction marker", got)
	}
}
