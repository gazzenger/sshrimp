#!/bin/bash

while true; do
    echo "This script will install the SSHrimp-Agent for the current user ONLY"
    echo "The application will be installed to the directory ~/sshrimp/"
    echo "This script will also autostart this application by adding a .desktop file to the directory ~/.config/autostart/"
    echo "To uninstall this application, remove the sshrimp-agent.desktop file from ~/.config/autostart/, as well as the ~/sshrimp folder"
    read -p "Do you want to proceed? " yn
    case $yn in
        [Yy]* ) break;;
        [Nn]* ) exit;;
        * ) echo "Please answer yes or no.";;
    esac
done

chmod +x sshrimp-agent-linux
chmod +x sshrimp-agent-run.sh

# Check if process already is running
SSHRIMP_PID=$(pgrep "sshrimp-agent")
if [[ "" != "$SSHRIMP_PID" ]]
then
    kill -9 $SSHRIMP_PID
fi

mkdir -p ~/sshrimp
\cp -f sshrimp-agent-linux ~/sshrimp/
\cp -f sshrimp-linux.toml ~/sshrimp/

# profile.d folder runs on ALL sessions (which is not ideal)
# sudo \cp -f sshrimp-agent-run.sh /etc/profile.d/

# ~/.config/autostart/
# create autostart folder (if it doesn't exist yet)
mkdir -p ~/.config/autostart
\cp -f sshrimp-agent.desktop ~/.config/autostart/