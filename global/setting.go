package global

import (
    "github.com/spark8899/deploy-agent/pkg/logger"
    "github.com/spark8899/deploy-agent/pkg/setting"
)

var (
    ServerSetting   *setting.ServerSettingS
    AppSetting      *setting.AppSettingS
    Logger          *logger.Logger
)
