package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx       context.Context
	scheduler *Scheduler
	config    *Config
}

func NewApp() *App {
	return &App{
		scheduler: NewScheduler(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	cfg, _ := LoadConfig()
	a.config = cfg
	a.scheduler.ReloadTasks(cfg.Tasks)
	a.scheduler.Start()
}

func (a *App) shutdown(ctx context.Context) {
	a.scheduler.Stop()
}

func (a *App) GetTasks() []CleanTask {
	return a.config.Tasks
}

func (a *App) SaveTask(task CleanTask) error {
	if task.ID == "" {
		task.ID = uuid.New().String()
		a.config.Tasks = append(a.config.Tasks, task)
	} else {
		for i, t := range a.config.Tasks {
			if t.ID == task.ID {
				a.config.Tasks[i] = task
				break
			}
		}
	}
	err := SaveConfig(a.config)
	if err == nil {
		a.scheduler.ReloadTasks(a.config.Tasks)
	}
	return err
}

func (a *App) DeleteTask(id string) error {
	for i, t := range a.config.Tasks {
		if t.ID == id {
			a.config.Tasks = append(a.config.Tasks[:i], a.config.Tasks[i+1:]...)
			break
		}
	}
	err := SaveConfig(a.config)
	if err == nil {
		a.scheduler.ReloadTasks(a.config.Tasks)
	}
	return err
}

func (a *App) RunTaskNow(id string) (*TaskResult, error) {
	for _, t := range a.config.Tasks {
		if t.ID == id {
			return RunTask(t)
		}
	}
	return &TaskResult{Errors: []string{}}, nil
}

func (a *App) SelectDirectory() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择目录",
	})
}

func (a *App) SelectFile() (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择文件",
	})
}

func (a *App) GetErrorLogs() (string, error) {
	exe, _ := os.Executable()
	logPath := filepath.Join(filepath.Dir(exe), "error.log")
	data, err := os.ReadFile(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(data), nil
}

func (a *App) ClearErrorLogs() error {
	exe, _ := os.Executable()
	logPath := filepath.Join(filepath.Dir(exe), "error.log")
	return os.WriteFile(logPath, []byte{}, 0644)
}