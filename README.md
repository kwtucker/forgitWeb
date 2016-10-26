![forgit logo](/forgit_md_logo.png)

# Forgit - Never forget to commit
Forgit is a workflow tool that automates basic git flow. The traditional process takes you away from your code and breaks your concentration. Forgit will add, commit, and push your code at set times automatically with custom commit messages based on your code.
***

### Forgit Go Version
* Built with Go v1.7

### Requirements
* Go installed
* Go environmental variables set
* NodeJS /w npm
* Gulp - (Globally)

### Install My Web App
```
  go get github.com/kwtucker/forgitWeb

  or

  git clone https://github.com/kwtucker/forgitWeb.git
```

### Dependencies
```
  go get github.com/julienschmidt/httprouter
  go get github.com/gorilla/context
  go get github.com/gorilla/sessions
  go get golang.org/x/oauth2
  go get golang.org/x/oauth2/github
  go get github.com/google/go-github/github
  go get gopkg.in/mgo.v2
  go get gopkg.in/mgo.v2/bson
  go get github.com/satori/go.uuid
  go get github.com/rs/cors
  npm install -g gulp
  npm install
```

### config.json ( Needs to be at root level in project )  
  * Fill in your own values
```json
{
    "SessionSecret": "string",
    "GithubClientID": "string from Github",
    "GithubClientSecret": "string from Github",
    "GithubState": "string",
    "WebPort": 8100,
    "DbPort": 27017,
    "DbName": "string",
    "TemplateDir": "views path",
    "TemplatePreCompile": true,
    "Debug": false,
    "GoogleAnalyticsTrackingID": "UA-(your GA id)",
    "WebHost": "127.0.0.1",
    "DbHost": "127.0.0.1",
    "StaticPath": "static path"
}
```

### Start App
```
gulp startup
```

___

## API Endpoints

##### GET User Settings
```
/api/users/:forgitId/:init
```
###### Sample
```json
[
  {
    "name": "General",
    "status": 1,
    "notifications": {
      "OnError": 1,
      "OnCommit": 1,
      "OnPush": 1
    },
    "addPullCommit": {
      "TimeMin": 1
    },
    "push": {
      "TimeMin": 2
    },
    "repos": [
      {
        "GithubRepoID": 0,
        "Name": "repo_name",
        "Status": 0
      },
    ]
  }
]
```

##### GET Update Check User Settings

```
/api/users/:forgitId/:no
```

###### Sample
```json
{
  "update": "0"
}
```

##### GET If Bad UUID

###### Sample
```json
{
  "message": "bad credentials",
  "status": 401
}
```
