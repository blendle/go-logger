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
  logger := logger.Must(logger.New("my-service", "cf89f839"))
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

logger := logger.Must(logger.New("my-service", "cf89f839", sampler, fields))
```

### Stackdriver logging

You can optionally add Stackdriver specific fields to your logs. These can be
used by Stackdriver to [improve log readability/grouping][sd].

```golang
import zapdriver "github.com/blendle/zapdriver"
```

```golang
logger.Info("Hello", zapdriver.Label("hello", "world"))
```

See [here][zd] for all available Stackdriver fields.

[sd]: https://cloud.google.com/error-reporting/docs/formatting-error-messages
[zd]: https://github.com/blendle/zapdriver#special-purpose-logging-fields

## Debugging

You can send the `USR1` signal to your application to switch the log level
between the default `INFO` and `DEBUG` level on runtime.

This allows you to capture debug logs during anomalies and find the problem.

You can also set the `DEBUG` environment variable to `true` to have the
application launch with the default log level set to `DEBUG` instead of `INFO`.

Again, you can send `USR1` to toggle back to `INFO` as well.

## Testing

This package contains a public testing API you can use if you need to assert
a log entry exists.

```golang
// TestNew calls New, but returns both the logger, and an observer that can be
// used to fetch and compare delivered logs.
TestNew(tb testing.TB, options ...zap.Option) (*zap.Logger, *observer.ObservedLogs)
```

```golang
// TestNewWithLevel is equal to TestNew, except that it takes an extra argument,
// dictating the minimum log level required to record an entry in the recorder.
TestNewWithLevel(tb testing.TB, level zapcore.LevelEnabler, options ...zap.Option) (*zap.Logger, *observer.ObservedLogs)
```

see `testing.go` for more details.
