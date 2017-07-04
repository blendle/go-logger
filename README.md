# go-logger

This package contains shared logic for Go-based logging.

This package uses the [`zap`][zap] logging package under the hood. See its
[README][zap] for more details on the logging API.

See [this wiki page][wiki] for more details on logging at Blendle.

[zap]: https://github.com/uber-go/zap#zap-zap---
[wiki]: https://www.notion.so/blendle/Structured-Logging-5e6fb67cd17c42acb7f180c37f436c2a

## Usage

Add to your `main.go`:

```golang
import logger "github.com/blendle/go-logger"

func init() {
  c := &logger.Config{
    App:         "my-app",
    Tier:        "api",
    Production:  false,
    Version:     "cf89f839",
    Environment: "staging",
  }

  logger.Init(c)
}
```

Then use it throughout your application:

```golang
logger.L.Warn("Something happened!")
```

### Custom Zap Configuration

You can also provide custom Zap configuration on initialization, if you need it:

```golang
options := func(c zap.Config) {
  c.Sampling = &zap.SamplingConfig{
    Initial:    100,
    Thereafter: 100,
  }
}

logger.Init(c, options)
```

## Debugging

You can send the `usr1` signal to your application to switch the log level
between the default `INFO` and `DEBUG` level on runtime.

This allows you to capture debug logs during anomalies and find the problem.

You can also set the `DEBUG` environment variable to `true` to have the
application launch with the default log level set to `DEBUG` instead of `INFO`.

Again, you can send `usr1` to toggle back to `INFO` as well.
