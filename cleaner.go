package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type TaskResult struct {
	Success int      `json:"success"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors"`
}

func logError(taskName, msg string) {
	exe, _ := os.Executable()
	logPath := filepath.Join(filepath.Dir(exe), "error.log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("[%s] [%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), taskName, msg))
}

func RunTask(task CleanTask) (*TaskResult, error) {
	result := &TaskResult{Errors: []string{}}

	if task.Mode == "truncate" {
		f, err := os.OpenFile(task.Path, os.O_RDWR|os.O_TRUNC, 0666)
		if err != nil {
			result.Failed = 1
			errMsg := fmt.Sprintf("%s: %v", task.Path, err)
			result.Errors = append(result.Errors, errMsg)
			logError(task.Name, errMsg)
			return result, nil
		}
		f.Close()
		result.Success = 1
		return result, nil
	}

	if task.Mode == "retention" {
		info, err := os.Stat(task.Path)
		if err != nil {
			return result, err
		}
		// 如果是文件，直接检查并删除
		if !info.IsDir() {
			if matchPattern(task.Path, task.FilePattern) && time.Since(info.ModTime()).Hours()/24 > float64(task.RetentionDays) {
				if err := os.Remove(task.Path); err != nil {
					result.Failed = 1
					errMsg := fmt.Sprintf("%s: %v", task.Path, err)
					result.Errors = append(result.Errors, errMsg)
					logError(task.Name, errMsg)
				} else {
					result.Success = 1
				}
			}
			return result, nil
		}
		// 如果是目录，遍历删除
		filepath.Walk(task.Path, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if matchPattern(path, task.FilePattern) && time.Since(info.ModTime()).Hours()/24 > float64(task.RetentionDays) {
				if err := os.Remove(path); err != nil {
					result.Failed++
					errMsg := fmt.Sprintf("%s: %v", path, err)
					result.Errors = append(result.Errors, errMsg)
					logError(task.Name, errMsg)
				} else {
					result.Success++
				}
			}
			return nil
		})
	}
	return result, nil
}

func matchPattern(path, pattern string) bool {
	if pattern == "" || pattern == "*" {
		return true
	}
	return strings.HasSuffix(strings.ToLower(path), strings.ToLower(pattern))
}