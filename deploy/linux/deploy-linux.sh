#!/bin/bash

while true; do
    read -p "This script will install the SSHrimp-Agent for the current user, do you want to proceed? " yn
    case $yn in
        [Yy]* ) break;;
        [Nn]* ) exit;;
        * ) echo "Please answer yes or no.";;
    esac
done

mkdir -p ~/sshrimp
\cp -f sshrimp-agent-linux ~/sshrimp/
\cp -f sshrimp-linux.toml ~/sshrimp/

if grep -q "~/sshrimp/sshrimp-agent-linux ~/sshrimp/sshrimp-linux.toml" "~/.bashrc" ; then
else
    echo '~/sshrimp/sshrimp-agent-linux ~/sshrimp/sshrimp-linux.toml' >> ~/.bashrc
fi
