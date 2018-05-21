# go-logger

This package contains shared logic for Go-based structured logging.

This package uses the [`zap`][zap] logging package under the hood. See its
[README][zap] for more details on the logging API.

[zap]: https://github.com/uber-go/zap#zap-zap---

## Usage

Add to your `main.go`:

```golang
import logger "github.com/blendle/go-logger"

func main() {
  logger := logger.New("my-service", "cf89f839")
}
```

Then use it throughout your application:

```golang
logger.Warn("Something happened!")
```

### Custom Zap Options

You can also provide custom Zap options on initialization, if you need it:

```golang
sampler := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
  return zapcore.NewSampler(core, time.Second, 100, 100)
})

fields := zap.Fields(zap.String("alwaysAdd", "this"))

logger := logger.New("my-service", "cf89f839", sampler, fields)
```

## Debugging

You can send the `USR1` signal to your application to switch the log level
between the default `INFO` and `DEBUG` level on runtime.

This allows you to capture debug logs during anomalies and find the problem.

You can also set the `DEBUG` environment variable to `true` to have the
application launch with the default log level set to `DEBUG` instead of `INFO`.

Again, you can send `USR1` to toggle back to `INFO` as well.
