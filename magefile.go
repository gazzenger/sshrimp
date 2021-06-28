//+build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"

	// mage:import ca
	"github.com/gazzenger/sshrimp/tools/mage/ca"
	// mage:import agent
	"github.com/gazzenger/sshrimp/tools/mage/agent"
)

var Default = All

// Builds all the targets
func Build() {
	mg.Deps(ca.Build, agent.Build)
}

// Builds for all platforms and package
// providing the config file is required, as it needs to be read
// and used for generating config files for each platform
func BuildAndPackage(configFile string) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("Config File %s doesn't exist\n", configFile)
		return
	}
	mg.Deps(agent.BuildAll, mg.F(agent.PackageFiles, configFile))
	fmt.Println("All done.")
}

// Remove all build output (except generated configuration files)
func Clean() {
	mg.Deps(ca.Clean, agent.Clean)
}

// Build and deploy the ca and agent
func All() {
	mg.Deps(agent.Build, ca.Package, ca.Generate)

	if _, err := os.Stat("./terraform"); os.IsNotExist(err) {
		fmt.Println("All done. Run `terraform init` then `terraform apply` to deploy.")
	} else {
		fmt.Println("All done. Run `terraform apply` to deploy.")
	}
}
