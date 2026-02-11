---
path: modules/logging.md
page-type: module
summary: Documentation for the structured logging system including log levels, the FuncLog callback, and logged events.
tags: [module, logging, funclog, levels, observability]
created: 2025-02-11
updated: 2025-02-11
version: 2.0.0
---

# Module: Logging

## Purpose

The logging module provides optional structured logging for key events during request handling. When `FuncLog` is configured, the package logs request handling, before-action aborts, and CSRF validation failures.

**File:** `v2/log.go`

## Log Levels

```go
const LogLevelInfo  = "info"
const LogLevelWarn  = "warn"
const LogLevelError = "error"
const LogLevelDebug = "debug"
```

## FuncLog Callback

```go
FuncLog: func(level string, message string, attrs map[string]any) {
    slog.Log(context.Background(), toSlogLevel(level), message, attrsToArgs(attrs)...)
}
```

The callback receives:
- **level** - One of `"info"`, `"warn"`, `"error"`, `"debug"`
- **message** - Human-readable event description
- **attrs** - Structured context as key-value pairs

## Logged Events

| Event | Level | Message | Attributes |
|-------|-------|---------|------------|
| Request received | `info` | `"handling request"` | `action`, `method`, `url` |
| Before-action abort | `warn` | `"request aborted by before-action hook"` | `action` |
| CSRF validation failure | `warn` | `"CSRF validation failed"` | `action`, `error` |

## Internal Helper

```go
func (crud *Crud) log(level string, message string, attrs map[string]any) {
    if crud.funcLog == nil {
        return
    }
    crud.funcLog(level, message, attrs)
}
```

The `log()` method is nil-safe - it silently returns when `FuncLog` is not configured.

## Integration Example

### With Go's slog

```go
import "log/slog"

FuncLog: func(level string, message string, attrs map[string]any) {
    args := make([]any, 0, len(attrs)*2)
    for k, v := range attrs {
        args = append(args, k, v)
    }
    switch level {
    case crud.LogLevelInfo:
        slog.Info(message, args...)
    case crud.LogLevelWarn:
        slog.Warn(message, args...)
    case crud.LogLevelError:
        slog.Error(message, args...)
    case crud.LogLevelDebug:
        slog.Debug(message, args...)
    }
},
```

### With fmt (simple)

```go
FuncLog: func(level string, message string, attrs map[string]any) {
    fmt.Printf("[%s] %s %v\n", level, message, attrs)
},
```

## See Also

- [Configuration](../configuration.md) - FuncLog configuration
- [Modules: Crud Core](crud_core.md) - Handler middleware pipeline
- [Architecture](../architecture.md) - Request lifecycle
