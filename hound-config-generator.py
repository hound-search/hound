import requests
from requests.auth import HTTPBasicAuth
from urllib.parse import parse_qsl, urlsplit
import json
import os

git_user = ''
git_token = ''
git_api = 'https://api.github.com/orgs/'
git_org = ''
git_base_url = git_api + git_org + '/repos'

repo_list = {}


def get_from_get_api(url):
    response = requests.get(url, auth=HTTPBasicAuth(git_user, git_token))
    return response


def get_paged_url_from_git_api(number_of_page):
    for page in range(1, int(number_of_page) + 1):
        github_paged_url = git_api + git_org + '/repos?page=' + str(page)
        response = get_from_get_api(github_paged_url).json()
        for repo in response:
            repo_name = repo['name']
            repo_ssh_url = repo['ssh_url']
            repo_url_info = {'url': repo_ssh_url}
            repo_list[repo_name] = repo_url_info


def get_last_page_from_git_respons(git_response):
    last_page = parse_qsl(urlsplit(git_response.links['last']['url']).query)[-1][1]
    return last_page


def move_old_config_file(file_name):
    if os.path.isfile('./' + file_name):
        os.rename(file_name, file_name + '.OLD')


if __name__ == "__main__":
    move_old_config_file('config.json')
    response_object = get_from_get_api(git_base_url)
    number_of_pages = (get_last_page_from_git_respons(response_object))
    get_paged_url_from_git_api(number_of_pages)

    hound_config = {'max-concurrent-indexers': 4, 'dbpath': 'data', 'repos': repo_list}
    config_file = open("config.json", "w")
    config_file.write(json.dumps(hound_config, indent=2, sort_keys=True))
    config_file.close()


