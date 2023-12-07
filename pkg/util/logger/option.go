package logger

type Option func(*Logger)

func WithLevel(lvl Level) Option {
	return func(l *Logger) {
		l.level = lvl
	}
}

func WithCid(cid string) Option {
	return func(l *Logger) {
		l.correlationId = cid
	}
}
