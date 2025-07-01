# Tiny logging API

This library provides a basic API to be integrated with any logging framework. It
allows to use logging in your libraries without enforcing a particular logging
framework.

Here is it:
```golang
// Logger interface
type Logger interface {

	// Debug logging: Every details
	Debug(event string, keyvals ...interface{})

	// Info logging: Core events
	Info(event string, keyvals ...interface{})

	// Warning logging: Anything out of the ordinary but non-life threatening
	Warn(event string, keyvals ...interface{})

	// Error logging: Major issue
	Error(event string, keyvals ...interface{})

	// Panic logging: We want to crash
	Panic(event string, keyvals ...interface{})

	// Context extending interface
	With(keyvals ...interface{}) Logger
}
```

It was designed for the [ftpserverlib](https://github.com/fclairamb/ftpserverlib).

## Compatible logging frameworks

### [slog](https://pkg.go.dev/log/slog)
```golang
import (
	adapter "github.com/fclairamb/go-log/slog"
)

func main() {
	logger := adapter.NewWrap(slog.New())

	logger.Info("Hello world !")
}
```

### [go-kit/log](https://github.com/go-kit/log)
```golang
import (
	"os"

	gklog "github.com/go-kit/log"
	adapter "github.com/fclairamb/go-log/gokit"
)

func main() {
	gkLogger := gklog.NewLogfmtLogger(gklog.NewSyncWriter(os.Stdout))
	logger := adapter.NewWrap(gkLogger)

	logger.Info("Hello world !")
}
```

### [log15](https://github.com/inconshreveable/log15)
```golang
import (
	"github.com/inconshreveable/log15"
	adapter "github.com/fclairamb/go-log/log15"
)

func main() {
	log15Logger := log15.New()
	logger := adapter.NewWrap(log15Logger)

	logger.Info("Hello world !")
}
```

### [zap](https://github.com/uber-go/zap)
```golang
import (
	"go.uber.org/zap"
	adapter "github.com/fclairamb/go-log/zap"
)

func main() {
	innerLogger, _ := zap.NewProduction()
	logger := adapter.NewWrap(innerLogger.Sugar())

	logger.Info("Hello world !")
}
```

### [zerolog](https://github.com/rs/zerolog/)
```golang
import (
	zerolog "github.com/rs/zerolog/log"
	adapter "github.com/fclairamb/go-log/zerolog"
)

func main() {
	logger := adapter.NewWrap(&zerolog.Logger)

	logger.Info("Hello world !")
}

```

### [logrus](https://github.com/sirupsen/logrus) (not recommended)
```golang
import (
	"github.com/sirupsen/logrus" //nolint: depguard
	adapter "github.com/fclairamb/go-log/logrus"
)

func main() {
	logrusLogger := logrus.New()
	logger := adapter.NewWrap(logrusLogger)

	logger.Info("Hello world !")
}
```
