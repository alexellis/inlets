#!/usr/bin/env bash

# Source: https://www.docker.com/blog/multi-arch-build-what-about-travis/

sudo rm -rf /var/lib/apt/lists/*
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) edge"
sudo apt-get update -qy
sudo apt-get -y -o Dpkg::Options::="--force-confnew" install docker-ce
mkdir -vp ~/.docker/cli-plugins/
curl --silent -L "https://github.com/docker/buildx/releases/download/v0.3.0/buildx-v0.3.0.linux-amd64" > ~/.docker/cli-plugins/docker-buildx
chmod a+x ~/.docker/cli-plugins/docker-buildx
sudo systemctl start docker

