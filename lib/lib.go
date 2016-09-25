package lib

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/kwtucker/forgit/models"
	"golang.org/x/oauth2"
	"time"
)

// CreateUser ...
func CreateUser(user *github.User, repos []github.Repository) *models.User {
	var (
		repoArr      = []models.Repo{}
		settingRepos = []models.SettingRepo{}
		settings     = []models.Setting{}
	)

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

	currentUserSettings := models.Setting{
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
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		fmt.Println(err)
	}

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

func GetTokenStruct(token string) *oauth2.Token {
	// getting new
	var tok = oauth2.Token{
		AccessToken: token,
	}
	var tokpointer = &tok
	return tokpointer
}
