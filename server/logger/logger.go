package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var Business = log.New(io.Discard, "", log.LstdFlags|log.Lmicroseconds)
var Error = log.New(io.Discard, "", log.LstdFlags|log.Lmicroseconds)

func Init(filename string) error {
	if strings.TrimSpace(filename) == "" {
		filename = "logs/server.log"
	}
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}
	businessFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	errorName := strings.TrimSuffix(filename, filepath.Ext(filename)) + ".error" + filepath.Ext(filename)
	errorFile, err := os.OpenFile(errorName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		_ = businessFile.Close()
		return err
	}
	Business = log.New(io.MultiWriter(os.Stdout, businessFile), "", log.LstdFlags|log.Lmicroseconds)
	Error = log.New(io.MultiWriter(os.Stderr, errorFile), "", log.LstdFlags|log.Lmicroseconds)
	return nil
}

func SQL(statement string, args ...interface{}) { Business.Printf("sql=%q args=%v", statement, args) }
