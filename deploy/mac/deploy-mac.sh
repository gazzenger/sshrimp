#!/bin/bash

while true; do
    read -p "This script will install the SSHrimp-Agent as a Launcher for the current user, do you want to proceed? " yn
    case $yn in
        [Yy]* ) break;;
        [Nn]* ) exit;;
        * ) echo "Please answer yes or no.";;
    esac
done

mkdir -p ~/sshrimp
\cp -f sshrimp-agent-mac ~/sshrimp/
\cp -f sshrimp-mac.toml ~/sshrimp/
\cp -f com.user.sshrimp.plist ~/Library/LaunchAgents/
launchctl load ~/Library/LaunchAgents/com.user.sshrimp.plist
launchctl start com.user.sshrimp