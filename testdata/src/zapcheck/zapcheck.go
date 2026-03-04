package zapcheck

import "go.uber.org/zap"

func loggerCalls() {
	logger := zap.NewNop()

	// rule 1: lowercase
	logger.Info("Starting server") // want `log message should start with a lowercase letter`
	logger.Info("starting server")

	// rule 2: english only
	logger.Error("\u043e\u0448\u0438\u0431\u043a\u0430 \u043f\u043e\u0434\u043a\u043b\u044e\u0447\u0435\u043d\u0438\u044f") // want `log message must be in English, found non-Latin characters`
	logger.Error("connection error")

	// rule 3: special chars
	logger.Warn("warning!") // want `log message should not contain special characters or emoji`
	logger.Warn("warning")

	// rule 4: sensitive data
	tok := "abc"
	logger.Info("token: " + tok) // want `log message may contain sensitive data`
	logger.Info("token validated")
}

func sugarCalls() {
	logger := zap.NewNop()
	sugar := logger.Sugar()

	sugar.Infow("Starting sugar") // want `log message should start with a lowercase letter`
	sugar.Infow("starting sugar")

	sugar.Errorf("\u043e\u0448\u0438\u0431\u043a\u0430") // want `log message must be in English, found non-Latin characters`
	sugar.Errorf("error happened: code %d", 42)
}
