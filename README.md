# Zeus

Zeus is a very simple utility package for context-based logging. In its core, it provides syntactic
sugar for associating a logger with `context.Context` objects and retrieving it from these
contexts.

Zeus uses the "blazing-fast and structured" logging library [zap](https://github.com/uber-go/zap)
under the hood. By not using `zap` directly, interfaces in your code become less dependent on a
logging library as you do not need to pass a logger around. Rather, the logger is "silently"
attached to a context object that is deeply rooted in the Go language.

## Installation

```bash
go get github.com/borchero/zeus
```

## Logging Quickstart

```go
func main() {
    // Ensures that all logs are flushed prior to program exit
    defer zeus.Sync()

    ctx := context.Background()
    // You can obtain a global logger instance from *any* context object in your application!
    logger := zeus.Logger(ctx)
    logger.Debug("hello")
    // >>> prints "DEBUG   hello"

    // You can do structured logging by customizing the logger that is associated with a context.
    // For example, you can give it a name...
    ctx = zeus.WithName(ctx, "test")
    zeus.Logger(ctx).Info("world")
    // >>> prints "INFO    test    world"

    // ...or add some fields to the logger.
    ctx = zeus.WithFields(ctx, zap.String("name", "borchero"))
    zeus.Logger(ctx).Info("new world")
    // >>> prints "INFO    test    new world       {"name": "borchero"}"
}
```

Note that `zeus.Logger` always returns a `*zap.Logger` instance which allows Zeus to be equally
powerful as `zap`.

### Log Configuration

In the example above, you see that we can use a logger without explicitly creating one beforehand.
This logger can be customized either from environment variables (recommended) or from code.

When using environment variables, you may set `GO_LOG` and `GO_LOG_FORMAT`. `GO_LOG` allows for
`debug`, `info`, `warn`, `error`, and `fatal` (and their uppercase counterparts) while
`GO_LOG_FORMAT` may be set to `human` or `json`. By default, Zeus uses the following configuration:

- `GO_LOG=info`
- `GO_LOG_FORMAT=human`

When you need utmost customization, you can create your own `*zap.Logger` and attach it to a
context via `zeus.WithLogger`. Be aware that such a logger is only available for the context passed
to `zeus.WithLogger` and the context's descendents. `zeus.Sync` will not ensure that all logs
managed by this logger are printed and you will need to call `Sync` on this logger yourself.
