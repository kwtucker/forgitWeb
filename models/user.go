package models

import ()

// User ...
type User struct {
	ID         int `json:"id"`
	LastUpdate `json:"lastUpdate"`
	LastSync   `json:"lastSync"`
	Login      `json:"login"`
	Name       `json:"name"`
	AvatarURL  `json:"avatar_url"`
	Company    `json:"company"`
	HTMLURL    `json:"html_url"`
	ReposURL   `json:"repos_url"`
	Repos      []Repo    `json:"repos"`
	Settings   []Setting `json:"settings"`
}

// Repo ...
type Repo struct {
	URL             string `json:"url"`
	CommitsURL      string `json:"commits_url"`
	ContributorsURL string `json:"contributors_url"`
	Description     string `json:"description"`
	FullName        string `json:"full_name"`
	GitCommitsURL   string `json:"git_commits_url"`
	HTMLURL         string `json:"html_url"`
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Owner           string `json:"owner"`
}

// Setting ...
type Setting struct {
	ID     int
	Name   string
	Status int
	SettingNotifications
	SettingAddPullCommit
	SettingPush
	Repos []SettingRepo
}

// SettingNotifications ...
type SettingNotifications struct {
	Status          int
	TimeNoCommitMin int
	OnCommit        int
	OnPush          int
}

// SettingAddPullCommit ...
type SettingAddPullCommit struct {
	Status  int
	TimeMin int
}

// SettingPush ...
type SettingPush struct {
	Status  int
	TimeMin int
}

// SettingRepo ...
type SettingRepo struct {
	ID     int
	Name   string
	Status int
}
