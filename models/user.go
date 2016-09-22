package models

// User ...
type User struct {
	ID         *int      `json:"id"`
	LastUpdate string    `json:"lastUpdate"`
	LastSync   string    `json:"lastSync"`
	Login      *string   `json:"login"`
	Name       *string   `json:"name"`
	AvatarURL  *string   `json:"avatar_url"`
	Company    *string   `json:"company"`
	HTMLURL    *string   `json:"html_url"`
	ReposURL   *string   `json:"repos_url"`
	Repos      []Repo    `json:"repos"`
	Settings   []Setting `json:"settings"`
}

// Repo ...
type Repo struct {
	URL             *string `json:"url"`
	CommitsURL      *string `json:"commits_url"`
	ContributorsURL *string `json:"contributors_url"`
	Description     *string `json:"description"`
	FullName        *string `json:"full_name"`
	GitCommitsURL   *string `json:"git_commits_url"`
	HTMLURL         *string `json:"html_url"`
	ID              *int    `json:"id"`
	Name            *string `json:"name"`
	Owner           *string `json:"owner"`
}

// Setting ...
type Setting struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`
	SettingNotifications
	SettingAddPullCommit
	SettingPush
	Repos []SettingRepo
}

// SettingNotifications ...
type SettingNotifications struct {
	Status   int `json:"status"`
	OnError  int `json:"onError"`
	OnCommit int `json:"onCommit"`
	OnPush   int `json:"onPush"`
}

// SettingAddPullCommit ...
type SettingAddPullCommit struct {
	Status  int `json:"status"`
	TimeMin int `json:"timeMinute"`
}

// SettingPush ...
type SettingPush struct {
	Status  int `json:"status"`
	TimeMin int `json:"timeMinute"`
}

// SettingRepo ...
type SettingRepo struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}
