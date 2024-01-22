package logger

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type LoggerV1 interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
}

type Field struct {
	Key string
	Val any
}

func exampleV1() {
	var l LoggerV1
	l.Info("这是一个新用户", Field{Key: "union_id", Val: 123})
}

type LoggerV2 interface {
	//要求args必须是偶数，第一个是key,第二个是value
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}
