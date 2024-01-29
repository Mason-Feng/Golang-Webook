package logger

func Error(err error) Field {
	return Field{
		Key: "error",
		Val: err,
	}
}
func Int64(key string, val int64) Field {
	return Field{Key: key, Val: val}
}
