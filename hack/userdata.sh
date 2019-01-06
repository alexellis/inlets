#!/bin/bash

# Enables randomly generated authentication token by default.
# Change the value here if you desire a specific token value.
export INLETSTOKEN=$(head -c 16 /dev/urandom | shasum | cut -d" " -f1)

curl -sLO https://github.com/alexellis/inlets/releases/download/0.4.0/inlets && \
  mv inlets /usr/bin/inlets && \
  chmod +x /usr/bin/inlets

curl -sLO https://raw.githubusercontent.com/alexellis/inlets/master/hack/inlets.service  && \
  mv inlets.service /etc/systemd/system/inlets.service && \
  echo "AUTHTOKEN=$INLETSTOKEN" > /etc/default/inlets && \
  systemctl start inlets && \
  systemctl enable inlets




