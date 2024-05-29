package db

type Config struct {
	TableName string
}

var (
	levelError = "error"
	spanKey    = "span"
	traceKey   = "trace"
)
