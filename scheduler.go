package main

import (
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron     *cron.Cron
	entryMap map[string]cron.EntryID
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		cron:     cron.New(),
		entryMap: make(map[string]cron.EntryID),
	}
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}

func (s *Scheduler) AddTask(task CleanTask) error {
	if !task.Enabled {
		return nil
	}
	entryID, err := s.cron.AddFunc(task.CronSpec, func() {
		RunTask(task)
	})
	if err != nil {
		return err
	}
	s.entryMap[task.ID] = entryID
	return nil
}

func (s *Scheduler) RemoveTask(taskID string) {
	if entryID, ok := s.entryMap[taskID]; ok {
		s.cron.Remove(entryID)
		delete(s.entryMap, taskID)
	}
}

func (s *Scheduler) ReloadTasks(tasks []CleanTask) {
	for id := range s.entryMap {
		s.RemoveTask(id)
	}
	for _, task := range tasks {
		s.AddTask(task)
	}
}