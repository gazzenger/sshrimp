package agent

import (
	"runtime"

	"github.com/magefile/mage/sh"
)

// Build Builds the local ssh agent
func Build() error {

	// if running windows, compile for all platforms
	if runtime.GOOS == "windows" {

		env := map[string]string{
			"GOOS":   "linux",
			"GOARCH": "amd64",
		}
		sh.RunWith(env, "go", "build", "-o", "sshrimp-agent-linux", "./cmd/sshrimp-agent")

		env = map[string]string{
			"GOOS":   "darwin",
			"GOARCH": "amd64",
		}
		sh.RunWith(env, "go", "build", "-o", "sshrimp-agent-mac", "./cmd/sshrimp-agent")

	}

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
