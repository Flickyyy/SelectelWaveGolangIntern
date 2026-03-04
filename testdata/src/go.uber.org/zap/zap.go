package zap

// Minimal stub of go.uber.org/zap for analysistest type-checking.

type Logger struct{}

type SugaredLogger struct{}

type Field struct{}

func NewNop() *Logger                  { return &Logger{} }
func NewProduction() (*Logger, error)  { return &Logger{}, nil }
func NewDevelopment() (*Logger, error) { return &Logger{}, nil }

func (l *Logger) Info(msg string, fields ...Field)  {}
func (l *Logger) Warn(msg string, fields ...Field)  {}
func (l *Logger) Error(msg string, fields ...Field) {}
func (l *Logger) Debug(msg string, fields ...Field) {}
func (l *Logger) Fatal(msg string, fields ...Field) {}
func (l *Logger) Panic(msg string, fields ...Field) {}

func (l *Logger) Sugar() *SugaredLogger { return &SugaredLogger{} }

func (s *SugaredLogger) Infow(msg string, keysAndValues ...interface{})  {}
func (s *SugaredLogger) Warnw(msg string, keysAndValues ...interface{})  {}
func (s *SugaredLogger) Errorw(msg string, keysAndValues ...interface{}) {}
func (s *SugaredLogger) Debugw(msg string, keysAndValues ...interface{}) {}
func (s *SugaredLogger) Infof(template string, args ...interface{})      {}
func (s *SugaredLogger) Errorf(template string, args ...interface{})     {}

func String(key, val string) Field { return Field{} }
