package log

import (
	"impower/LicenseServer/env"
)

var (
	// LoggerModeDebug : debug mode
	LoggerModeDebug bool = true
	// LoggerLevel : FINEST ,FINE ,DEBUG ,TRACE ,INFO ,WARNING, ERROR, CRITICAL -> 0 ~ 7
	LoggerLevel int = 2
	// LoggerConsoleFormat : logger console print format
	LoggerConsoleFormat string = "[%T %D] [%L] (%S) %M"
	// LoggerFileFormat : logger file print format
	LoggerFileFormat string = "[%T %D] [%L] (%S) %M"
	// LoggerFilePath : logger file path
	LoggerFilePath string = "./log"
	// LoggerFileMaxSize : logger file max size
	LoggerFileMaxSize int = 1048576
	// LoggerFileMaxLines : logger file max lines
	LoggerFileMaxLines int = 10000
	// LoggerFileBackup : logger file backup number
	LoggerFileBackup int = 10
)

func init() {
	LoggerModeDebug = env.Get("debug").(bool)
	LoggerLevel = int(env.Get("logger.level").(int64))
	LoggerFilePath = env.Get("logger.file.path").(string)
	LoggerConsoleFormat = env.Get("logger.console.format").(string)
	LoggerFileFormat = env.Get("logger.file.format").(string)
	LoggerFileMaxSize = int(env.Get("logger.file.size").(int64))
	LoggerFileMaxLines = int(env.Get("logger.file.lines").(int64))
	LoggerFileBackup = int(env.Get("logger.file.maxbackup").(int64))
}
