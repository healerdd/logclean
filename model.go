package main

type CleanTask struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Path          string `json:"path"`
	Mode          string `json:"mode"`
	RetentionDays int    `json:"retentionDays"`
	FilePattern   string `json:"filePattern"`
	CronSpec      string `json:"cronSpec"`
	Enabled       bool   `json:"enabled"`
}

type Config struct {
	Tasks []CleanTask `json:"tasks"`
}