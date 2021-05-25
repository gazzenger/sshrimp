package agent

import (
	"github.com/magefile/mage/sh"
)

// Build Builds the local ssh agent
func Build() error {
	// sh.Run("env GOOS=windows GOARCH=amd64")

	// env := map[string]string{
	// 	"GOOS":   "windows",
	// 	"GOARCH": "amd64",
	// }
	// sh.RunWith(env, "go", "build", "./cmd/sshrimp-agent")

	return sh.Run("go", "build", "./cmd/sshrimp-agent")
}

// Clean Cleans the output files for sshrimp-agent
func Clean() error {
	return sh.Rm("sshrimp-agent")
}

// Install Installs the sshrimp-agent
func Install() error {
	return sh.Run("go", "install", "./cmd/sshrimp-agent")
}
