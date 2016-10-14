package models

// APIError ...
type APIError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

//UpdateStatus ...
type UpdateStatus struct {
	Update string `json:"update"`
}

// APIUser ..
type APIUser struct {
	GithubID   int       `json:"githubID"`
	ForgitID   string    `json:"forgitID"`
	ForgitPath string    `json:"forgitPath"`
	UpdateTime string    `json:"updateTime"`
	Settings   []Setting `json:"settings,omitempty"`
}

// APISetting ...
type APISetting struct {
	Name                 string `json:"name"`
	Status               int    `json:"status"`
	SettingNotifications `json:"notifications"`
	SettingAddPullCommit `json:"addPullCommit"`
	SettingPush          `json:"push"`
	Repos                []SettingRepo `json:"repos"`
}

// APISettingNotifications ...
type APISettingNotifications struct {
	OnError  int `json:"onError"`
	OnCommit int `json:"onCommit"`
	OnPush   int `json:"onPush"`
}

// APISettingAddPullCommit ...
type APISettingAddPullCommit struct {
	TimeMin int `json:"timeMinute"`
}

// APISettingPush ...
type APISettingPush struct {
	TimeMin int `json:"timeMinute"`
}

// APISettingRepo ...
type APISettingRepo struct {
	GithubRepoID int    `json:"github_repo_id"`
	Name         string `json:"name"`
	Status       int    `json:"status"`
}
