package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	Debug *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

func Init() {
	logLevel := os.Getenv("LOG_LEVEL")

	// 创建 log 目录
	logDir := filepath.Join(getProjectRoot(), "log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
	}

	// 根据日期创建日志文件
	logFile := filepath.Join(logDir, fmt.Sprintf("app-%s.log", time.Now().Format("20060102")))
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		// 回退到只输出到 stdout
		Debug = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
		Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
		Warn = log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
		Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		// 同时输出到文件和控制台
		writers := io.MultiWriter(file, os.Stdout)
		Debug = log.New(writers, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
		Info = log.New(writers, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
		Warn = log.New(writers, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
		Error = log.New(io.MultiWriter(file, os.Stderr), "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Set Debug logger to discard if LOG_LEVEL != debug
	if logLevel != "debug" && logLevel != "DEBUG" {
		Debug.SetOutput(io.Discard)
	}
}

func init() {
	Init()
}

func getProjectRoot() string {
	// 从环境变量获取项目根目录
	if root := os.Getenv("PROJECT_ROOT"); root != "" {
		return root
	}
	// 回退方案：获取当前工作目录
	wd, _ := os.Getwd()
	return wd
}
