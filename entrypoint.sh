#!/bin/sh

eval $(ssh-agent -s) && ssh-add /root/.ssh/hound_id_ed25519

ssh-keyscan -t ssh-ed25519 private.vcs.instance >> /root/.ssh/known_hosts

/bin/houndd -conf /data/config.json