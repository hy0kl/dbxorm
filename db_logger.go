package dbdao

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/hy0kl/logger"
	"xorm.io/core"
)

var dbLogger = &dbLog{}

var (
	skip      = 3
	callerKey = "x_runtime_caller"
	tag       = "dbxorm"
	ctx       = context.Background()
)

type dbLog struct{}

func (d *dbLog) Debug(v ...interface{}) {
	logger.Dx(ctx, tag, fmt.Sprint(v...), callerKey, runtimeCaller(skip))
}

func (d *dbLog) Debugf(format string, v ...interface{}) {
	logger.Dx(ctx, tag, fmt.Sprintf(format, v...), callerKey, runtimeCaller(skip))
}

func (d *dbLog) Info(v ...interface{}) {
	logger.Ix(ctx, tag, fmt.Sprint(v...), callerKey, runtimeCaller(skip))
}

func (d *dbLog) Infof(format string, v ...interface{}) {
	if len(v) > 0 {
		duration, ok := v[len(v)-1].(time.Duration)
		if ok && duration < slowDuration {
			return
		}
	}
	logger.Ix(ctx, tag, fmt.Sprintf(format, v...), callerKey, runtimeCaller(skip))
}

func (d *dbLog) Warn(v ...interface{}) {
	logger.Wx(ctx, tag, fmt.Sprint(v...))
}

func (d *dbLog) Warnf(format string, v ...interface{}) {
	logger.Wx(ctx, tag, fmt.Sprintf(format, v...), callerKey, runtimeCaller(skip))
}

func (d *dbLog) Error(v ...interface{}) {
	logger.Ex(ctx, tag, fmt.Sprint(v...), callerKey, runtimeCaller(skip))
}

func (d *dbLog) Errorf(format string, v ...interface{}) {
	logger.Ex(ctx, tag, fmt.Sprintf(format, v...), callerKey, runtimeCaller(skip))
}

func (d *dbLog) Level() core.LogLevel {
	return core.LOG_INFO
}

func (d *dbLog) SetLevel(l core.LogLevel) {
	return
}

func (d *dbLog) ShowSQL(show ...bool) {
	return
}

func (d *dbLog) IsShowSQL() bool {
	return showSql
}

func runtimeCaller(s int) (position string) {
	_, file, line, ok := runtime.Caller(s)
	if ok {
		position = file + ":" + strconv.Itoa(line)
	} else {
		position = "EMPTY"
	}

	return
}
