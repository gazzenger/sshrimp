#!/bin/bash

rm sshrimp-ca sshrimp-ca.zip
mage
TARGET=$(cat sshrimp.toml | grep Regions | sed -e 's/  Regions = \["\(.*\)"]/\1/')
terraform destroy -target=module.sshrimp-$TARGET.aws_lambda_function.sshrimp_ca -auto-approve
terraform apply -auto-approve
./sshrimp-agent ./sshrimp.toml