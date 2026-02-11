package crud

const LogLevelInfo = "info"
const LogLevelWarn = "warn"
const LogLevelError = "error"
const LogLevelDebug = "debug"

func (crud *Crud) log(level string, message string, attrs map[string]any) {
	if crud.funcLog == nil {
		return
	}
	crud.funcLog(level, message, attrs)
}
