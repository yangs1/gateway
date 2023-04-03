package global

import (
	"gateway/initialize"
	"go.uber.org/zap"
)

var Logger *zap.Logger
var ServiceManagerHandler *initialize.ServiceManager
