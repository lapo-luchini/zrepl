package platformtest

import (
	"context"

	"github.com/dsh2dsh/zrepl/daemon/logging"
	"github.com/dsh2dsh/zrepl/logger"
)

type Logger = logger.Logger

func GetLog(ctx context.Context) Logger {
	return logging.GetLogger(ctx, logging.SubsysPlatformtest)
}
