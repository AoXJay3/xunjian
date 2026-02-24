package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

func initLogger(logDir string) {
	if err := os.MkdirAll(logDir, 0750); err != nil {
		log.Fatalf("创建日志目录失败: %v", err)
	}

	today := time.Now().Format("20060102")
	logFile := filepath.Join(logDir, "xunjian-"+today+".log")

	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		log.Fatalf("打开日志文件失败: %v", err)
	}

	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
}
