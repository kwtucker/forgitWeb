package lib

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/kwtucker/forgit/models"
	"golang.org/x/oauth2"
	"time"
)

// CreateUser ...
func CreateUser(user *github.User, repos []github.Repository, update []models.Setting) *models.User {
	var (
		repoArr             = []models.Repo{}
		settingRepos        = []models.SettingRepo{}
		settings            = []models.Setting{}
		currentUserSettings = models.Setting{}
	)

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

	// if update is a value set setting struct
	if update != nil {
		for u := range update {
			currentUserSettings = models.Setting{
				SettingID: update[u].SettingID,
				Name:      update[u].Name,
				Status:    update[u].Status,
				SettingNotifications: models.SettingNotifications{
					Status:   update[u].SettingNotifications.Status,
					OnError:  update[u].SettingNotifications.OnError,
					OnCommit: update[u].SettingNotifications.OnCommit,
					OnPush:   update[u].SettingNotifications.OnPush,
				},
				SettingAddPullCommit: models.SettingAddPullCommit{
					Status:  update[u].SettingAddPullCommit.Status,
					TimeMin: update[u].SettingAddPullCommit.TimeMin,
				},
				SettingPush: models.SettingPush{
					Status:  update[u].SettingPush.Status,
					TimeMin: update[u].SettingPush.TimeMin,
				},
				Repos: settingRepos,
			}
			settings = append(settings, currentUserSettings)
		}
	}

	// Update is nil means new user
	// Create user setting defualts
	if update == nil {
		currentUserSettings = models.Setting{
			SettingID: 1,
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
	}

	// set the time
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		fmt.Println(err)
	}

	// set time and build user struct
	timenow := &github.Timestamp{time.Now().In(location)}
	currentUser := &models.User{
		GithubID:   *user.ID,
		LastUpdate: timenow.String(),
		LastSync:   timenow.String(),
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
