package slogcheck

import (
	"context"
	"log/slog"
)

func packageLevelCalls() {
	// ---- rule 1: lowercase first letter ----
	slog.Info("Starting server on port 8080") // want `log message should start with a lowercase letter`
	slog.Error("Failed to connect")           // want `log message should start with a lowercase letter`
	slog.Info("starting server on port 8080")
	slog.Debug("connection established")

	// ---- rule 2: english only ----
	slog.Info("\u0437\u0430\u043f\u0443\u0441\u043a \u0441\u0435\u0440\u0432\u0435\u0440\u0430")                                                                                               // want `log message must be in English, found non-Latin characters`
	slog.Error("\u043e\u0448\u0438\u0431\u043a\u0430 \u043f\u043e\u0434\u043a\u043b\u044e\u0447\u0435\u043d\u0438\u044f \u043a \u0431\u0430\u0437\u0435 \u0434\u0430\u043d\u043d\u044b\u0445") // want `log message must be in English, found non-Latin characters`
	slog.Info("starting server")

	// ---- rule 3: special characters / emoji ----
	slog.Info("server started!\U0001f680") // want `log message should not contain special characters or emoji`
	slog.Error("connection failed!!!")     // want `log message should not contain special characters or emoji`
	slog.Warn("something went wrong...")   // want `log message should not contain special characters or emoji`
	slog.Info("server started successfully")

	// ---- rule 4: sensitive data ----
	password := "secret123"
	apiKey := "key-abc"
	token := "tok-456"
	slog.Info("user password: " + password) // want `log message may contain sensitive data`
	slog.Debug("api_key=" + apiKey)         // want `log message may contain sensitive data`
	slog.Info("token: " + token)            // want `log message may contain sensitive data`
	slog.Info("user authenticated successfully")
	slog.Info("token validated")
}

func methodCalls() {
	logger := slog.Default()
	logger.Info("Starting via logger method") // want `log message should start with a lowercase letter`
	logger.Info("starting via logger method")
}

func contextVariants() {
	ctx := context.Background()
	slog.InfoContext(ctx, "Starting with context") // want `log message should start with a lowercase letter`
	slog.InfoContext(ctx, "starting with context")
}
