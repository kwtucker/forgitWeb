package lib

import (
	// "fmt"
	"github.com/google/go-github/github"
	"github.com/kwtucker/forgit/models"
	"github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"strconv"
	"time"
)

// CreateUser ...
func CreateUser(user *github.User, repos []github.Repository, settingsUpdate []models.Setting) *models.User {
	var (
		repoArr             = []models.Repo{}
		settingRepos        = []models.SettingRepo{}
		settings            = []models.Setting{}
		currentUserSettings = models.Setting{}
		dateNow             string
		dn                  int64
	)
	// get unix time and convert it to a string for storage
	dn = time.Now().UTC().Unix()
	dateNow = strconv.FormatInt(dn, 10)

	// loop over all repos and set struct
	for k := range repos {
		currentUserRepos := models.Repo{
			URL:             repos[k].URL,
			CommitsURL:      repos[k].CommitsURL,
			ContributorsURL: repos[k].ContributorsURL,
			Description:     repos[k].Description,
			FullName:        repos[k].FullName,
			GitCommitsURL:   repos[k].GitCommitsURL,
			HTMLURL:         repos[k].HTMLURL,
			RepoID:          repos[k].ID,
			Name:            repos[k].Name,
			Owner:           repos[k].Owner.Login,
		}
		repoArr = append(repoArr, currentUserRepos)

		currentUserSettingsRepo := models.SettingRepo{
			GithubRepoID: repos[k].ID,
			Name:         repos[k].Name,
			Status:       0,
		}
		settingRepos = append(settingRepos, currentUserSettingsRepo)
	}

	// Update is nil means new user
	// Create user setting defualts
	currentUserSettings = models.Setting{
		SettingID: 0,
		Name:      "General",
		Status:    1,
		SettingNotifications: models.SettingNotifications{
			Status:   1,
			OnError:  1,
			OnCommit: 1,
			OnPush:   1,
		},
		SettingAddPullCommit: models.SettingAddPullCommit{
			Status:  1,
			TimeMin: 5,
		},
		SettingPush: models.SettingPush{
			Status:  1,
			TimeMin: 60,
		},
		Repos: settingRepos,
	}
	settings = append(settings, currentUserSettings)

	// if update is a value set setting struct
	if settingsUpdate != nil {
		for u := range settingsUpdate {
			currentUserSettings = models.Setting{
				SettingID: settingsUpdate[u].SettingID,
				Name:      settingsUpdate[u].Name,
				Status:    settingsUpdate[u].Status,
				SettingNotifications: models.SettingNotifications{
					Status:   settingsUpdate[u].SettingNotifications.Status,
					OnError:  settingsUpdate[u].SettingNotifications.OnError,
					OnCommit: settingsUpdate[u].SettingNotifications.OnCommit,
					OnPush:   settingsUpdate[u].SettingNotifications.OnPush,
				},
				SettingAddPullCommit: models.SettingAddPullCommit{
					Status:  settingsUpdate[u].SettingAddPullCommit.Status,
					TimeMin: settingsUpdate[u].SettingAddPullCommit.TimeMin,
				},
				SettingPush: models.SettingPush{
					Status:  settingsUpdate[u].SettingPush.Status,
					TimeMin: settingsUpdate[u].SettingPush.TimeMin,
				},
				Repos: settingRepos,
			}
			settings = append(settings, currentUserSettings)
		}
	}

	// Creating UUIDv4
	forgitID := uuid.NewV4()

	currentUser := &models.User{
		GithubID:   *user.ID,
		ForgitID:   forgitID.String(),
		LastUpdate: dateNow,
		LastSync:   dateNow,
		Login:      user.Login,
		Name:       user.Name,
		AvatarURL:  user.AvatarURL,
		Company:    user.Company,
		HTMLURL:    user.HTMLURL,
		ReposURL:   user.ReposURL,
		Repos:      repoArr,
		Settings:   settings,
	}

	return currentUser
}

// GetTokenStruct will create a new token struct and return a pointer to it
func GetTokenStruct(token string) *oauth2.Token {
	// getting new
	var tok = oauth2.Token{
		AccessToken: token,
	}
	var tokpointer = &tok
	return tokpointer
}
