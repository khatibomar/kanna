package core

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	dateFormat    = "2006-01-02 15-04-05"
	hoursPerMonth = float64(24 * 31)
)

var loggingDir = filepath.Join(getConfDir(), "logs")

func (t *Kanna) setUpLogging() error {
	if err := os.MkdirAll(loggingDir, os.ModePerm); err != nil {
		return err
	}

	now := time.Now()
	_ = filepath.Walk(loggingDir, func(path string, info fs.FileInfo, err error) error {
		fileDate := info.ModTime()
		if !info.IsDir() && now.Sub(fileDate).Hours() >= hoursPerMonth {
			_ = os.Remove(path)
		}
		return nil
	})

	formattedDate := now.Format(dateFormat)
	logFilePath := filepath.Join(loggingDir, fmt.Sprintf("%s.log", formattedDate))

	var err error
	if t.LogFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm); err != nil {
		return err
	}
	log.SetOutput(t.LogFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Printf("Session started at %s\n", formattedDate)

	return nil
}

func (t *Kanna) stopLogging() error {
	return t.LogFile.Close()
}
