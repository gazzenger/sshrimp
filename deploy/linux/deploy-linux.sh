#!/bin/bash

while true; do
    read -p "This script will install the SSHrimp-Agent for the current user, do you want to proceed? " yn
    case $yn in
        [Yy]* ) break;;
        [Nn]* ) exit;;
        * ) echo "Please answer yes or no.";;
    esac
done

chmod +x sshrimp-agent-linux
chmod +x sshrimp-agent-run.sh
mkdir -p ~/sshrimp
\cp -f sshrimp-agent-linux ~/sshrimp/
\cp -f sshrimp-linux.toml ~/sshrimp/

sudo \cp -f sshrimp-agent-run.sh /etc/profile.d/
