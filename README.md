![forgit logo](/forgit_md_logo.png)

# Forgit - Never forget to commit
Targeting junior to senior level developers who always look to increase their productivity. Forgit is a workflow tool that automates basic git flow. The traditional process takes you away from your code and breaks your concentration. Forgit will add, commit, and push your code at set times automatically with custom commit messages based on your code.
***

### Forgit Go Version
* Built with Go v1.7

### Requirements
* Go installed
* Go environmental variables set
* config.json file included
* NodeJS /w npm
* Gulp - (Globally)

### Install My Web App
```
  go get github.com/kwtucker/forgit
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
  npm install
```

### config.json ( needs to be at root level in project
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
    "StaticPath": "static"
}
```

## Start App
```
gulp startup
```
