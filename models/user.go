package models

// User ...
type User struct {
	GithubID   int       `bson:"githubID"`
	ForgitID   string    `json:"forgitID"`
	LastUpdate string    `bson:"lastUpdate"`
	LastSync   string    `bson:"lastSync"`
	Login      *string   `bson:"login"`
	Name       *string   `bson:"name"`
	AvatarURL  *string   `bson:"avatar_url"`
	Company    *string   `bson:"company"`
	HTMLURL    *string   `bson:"html_url"`
	ReposURL   *string   `bson:"repos_url"`
	Repos      []Repo    `bson:"repos"`
	Settings   []Setting `bson:"settings,omitempty"`
}

// Repo ...
type Repo struct {
	URL             *string `bson:"url"`
	CommitsURL      *string `bson:"commits_url"`
	ContributorsURL *string `bson:"contributors_url"`
	Description     *string `bson:"description"`
	FullName        *string `bson:"full_name"`
	GitCommitsURL   *string `bson:"git_commits_url"`
	HTMLURL         *string `bson:"html_url"`
	RepoID          *int    `bson:"repo_id"`
	Name            *string `bson:"name"`
	Owner           *string `bson:"owner"`
}

// Setting ...
type Setting struct {
	// SettingID int    `bson:"setting_id"`
	Name   string `bson:"name"`
	Status int    `bson:"status"`
	SettingNotifications
	SettingAddPullCommit
	SettingPush
	Repos []SettingRepo
}

// SettingNotifications ...
type SettingNotifications struct {
	// Status   int `bson:"status"`
	OnError  int `bson:"onError"`
	OnCommit int `bson:"onCommit"`
	OnPush   int `bson:"onPush"`
}

// SettingAddPullCommit ...
type SettingAddPullCommit struct {
	// Status  int `bson:"status"`
	TimeMin int `bson:"timeMinute"`
}

// SettingPush ...
type SettingPush struct {
	// Status  int `bson:"status"`
	TimeMin int `bson:"timeMinute"`
}

// SettingRepo ...
type SettingRepo struct {
	GithubRepoID *int    `bson:"github_repo_id"`
	Name         *string `bson:"name"`
	Status       int     `bson:"status"`
}
