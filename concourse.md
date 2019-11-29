# Concourse CI

* Overview and Docs: https://concourse-ci.org/
* Tutorial: https://concoursetutorial.com/
* Pipeline: [concourse.yml](concourse.yml)

# Deployment

### Concourse CI

First, download compose file:

    wget https://concourse-ci.org/docker-compose.yml

Now edit that file to set password and url:

    CONCOURSE_EXTERNAL_URL: http://ci.example.com:8080
    CONCOURSE_ADD_LOCAL_USER: ADMIN_LOGIN:ADMIN_PASSWORD
    CONCOURSE_MAIN_TEAM_LOCAL_USER: ADMIN_LOGIN

Start:

    docker-compose up -d

Now open the url and find link to download *CLI tools*.

Install fly CLI:

    sudo mkdir -p /usr/local/bin
    sudo mv ~/Downloads/fly /usr/local/bin
    sudo chmod 0755 /usr/local/bin/fly


Finally, create *Target*:

    fly --target hound login --concourse-url http://ci.example.com:8080 -u ADMIN_LOGIN -p ADMIN_PASSWORD

### Github Token

Create a token with the following access rights:

* ``write:packages``
* ``repo:status``

### Parameters file

Create a file ``concourse-params.yml`` of the the following content:

    github-username: USERNAME
    github-token: TOKEN
    webhook-token: ANY-RANDOM-STRING

* Use TOKEN that created for USERNAME

### Set Pipeline

Execute following command:

    fly --target=hound set-pipeline --pipeline=pull-requests --config=concourse.yml --load-vars-from=concourse-params.yml

### Github Webhook

At Github Repository configure webhook to [Concourse CI API](https://concourse-ci.org/resources.html#resource-webhook-token):

* **Payload URL**:  *http://CI.EXAMPLE.COM:8080/api/v1/teams/main/pipelines/pull-requests/resources/pr/check/webhook?webhook_token=WEBHOOK_TOKEN*
* **Content Type**: any value is ok
* **Secret**: leave empty
* **Which events would you like to trigger this webhook?**: *Let me select individual events.*:

  * **[v] Pull Requests**




