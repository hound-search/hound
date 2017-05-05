#!/usr/bin/env python

#########################
# Grab porchdotcom's list of repos and generate a config file so they can be automatically indexed into Hound.
#########################

import requests

# need to add authentication so we can grab private repos
github_url = 'https://api.github.com/users/porchdotcom/repos'


# replace this with an authenticated version (TODO: below) so we can grab all of the repos
repos = requests.get(github_url).json()

# http://docs.python-requests.org/en/latest/user/authentication/
# repos = requests.get(github_url, auth=('user', 'pass'))

size = len (repos)
reposAdded = 0

file = open('config.json', 'w+')

# non-repo specific portion of the config
configStart = '{\n    "dbpath" : "data",\n    "repos" : {\n'
#print configStart

file.write(configStart)


# add each repo to the config
for repo in repos:
    reposAdded = reposAdded + 1
    repoData = '        "' + repo["name"] + '" : {\n            "url" : "' + repo["clone_url"] + '"\n        }'
    if reposAdded != size:
        repoData = repoData + ',\n'
    #print repoData
    file.write(repoData)


# non-repo specific ending of config
configEnd = '\n    }\n}'
#print configEnd

file.write(configEnd)

file.close()
