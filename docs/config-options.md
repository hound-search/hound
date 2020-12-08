ConfigOption | Description
:------ | :-----
MaxConcurrentIndexers | defines the total number of indexers required to be used for indexing code. If not provided defaults `2`
HealthCheckURI |  health check url for hound , if not provided defaults to `/healthz`
DbPath | absolute file path where the `config.json` file exists. By default is `data`
title | Title used for the application.Defaulted to 'Hound'
url-pattern | composed of base url and anchor values in form of key value pairs.
vs-config | holds the version control config, default VCS used in Hound is git.Other options for VCS are svn,mercurial,bitbucket,hg, etc.Refer to `config-example.json` to get the list of vcs and usage
Repos | holds the list of repos which are required to be indexed by Hound . Each Repo is added with reponame as a Json Key with options associated with repo as values similar to example provided in `config-example.json`


gitOptions  | Description
:------ | :-----
ms-between-polls | time interval to poll the repo url ,default is `30s`
detect-ref    | used to determine branch , defaults to `master` branch 
ref | used to provide reference for the branch for repo.

svn-options  | Descriptions
------ | -----
username  | user name for the svn repo
password | password to authenticate use for svn repo.

url-options | Description
------ | -----
url-pattern | when provided used by Hound for config, else defaults to `{url}/blob/{rev}/{path}{anchor}`
anchor | when provided used for vcs config else set to `#L{line}`

