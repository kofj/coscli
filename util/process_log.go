package util

import (
	logger "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
)

var (
	processLogMu sync.Mutex
)

func writeProcessLog(errString string, fo *FileOperations) {
	if !fo.Operation.ProcessLog {
		return
	}
	var err error
	if fo.ProcessLogger.Path == "" {
		fo.ProcessLogger.Path = filepath.Join(fo.Operation.ProcessLogPath, fo.OutPutDirName)
		_, err := os.Stat(fo.ProcessLogger.Path)
		if os.IsNotExist(err) {
			err := os.MkdirAll(fo.ProcessLogger.Path, 0755)
			if err != nil {
				logger.Errorf("Failed to create process log dir: %v", err)
				return
			}
		}
	}

	if fo.ProcessLogger.logFile == nil {
		// 创建进程日志文件
		processLoggerFilePath := filepath.Join(fo.ProcessLogger.Path, "process.log")
		fo.ProcessLogger.logFile, err = os.OpenFile(processLoggerFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logger.Errorf("Failed to create process log file:%v", err)
			return
		}
	}

	processLogMu.Lock()

	_, writeErr := fo.ProcessLogger.logFile.WriteString(errString)

	if writeErr != nil {
		logger.Errorf("Failed to write process log  file : %v\n", writeErr)
	}
	processLogMu.Unlock()
}

// CloseProcessLoggerFile closes the process log file if it is not nil.
func CloseProcessLoggerFile(fo *FileOperations) {
	if fo.ProcessLogger.logFile != nil {
		defer fo.ProcessLogger.logFile.Close()
	}
}
