package agent

import (
	"compress/flate"
	"errors"
	"os"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/copier"
	"github.com/magefile/mage/sh"
	"github.com/mholt/archiver/v3"
)

// Agent config for the sshrimp-agent agent
type Agent struct {
	ProviderURL    string
	ClientID       string
	ClientSecret   string
	BrowserCommand []string
	Socket         string
}

// CertificateAuthority config for those few additional fields needed for clients TOML config files
type CertificateAuthority struct {
	AccountID int
	Regions   []string
}

// SSHrimp main configuration struct for sshrimp-agent and sshrimp-ca
type SSHrimp struct {
	Agent                Agent
	CertificateAuthority CertificateAuthority
}

// Build Builds the local ssh agent
func Build() error {
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

// BuildAll Builds the local ssh agent for Windows, Mac and Linux
func BuildAll() error {
	// if running windows, compile for all platforms
	if runtime.GOOS == "windows" {
		// build linux
		env := map[string]string{
			"GOOS":   "linux",
			"GOARCH": "amd64",
		}
		err := sh.RunWith(env, "go", "build", "-o", "./deploy/linux/sshrimp-agent-linux", "./cmd/sshrimp-agent")
		if err != nil {
			return err
		}
		// build mac (darwin)
		env = map[string]string{
			"GOOS":   "darwin",
			"GOARCH": "amd64",
		}
		err = sh.RunWith(env, "go", "build", "-o", "./deploy/mac/sshrimp-agent-mac", "./cmd/sshrimp-agent")
		if err != nil {
			return err
		}
		// build windows
		env = map[string]string{
			"GOOS":   "windows",
			"GOARCH": "amd64",
		}
		err = sh.RunWith(env, "go", "build", "-o", "./deploy/windows/sshrimp-agent-windows.exe", "./cmd/sshrimp-agent")
		return err
	} else {
		return errors.New("Building for all platforms can only be performed on a Windows System")
	}
}

// Package the deploy folder for clients
func PackageFiles(configFile string) error {
	//read in config file
	var c *SSHrimp
	_, err := toml.DecodeFile(configFile, &c)
	if err != nil {
		return err
	}

	// create new configs
	winConfig := SSHrimp{}
	unixConfig := SSHrimp{}

	// deep clone configs from original
	copier.Copy(&winConfig, &c)
	copier.Copy(&unixConfig, &c)

	// update the socket field in each config
	var socketSplit []string
	if len(c.Agent.Socket) > 9 && c.Agent.Socket[0:9] == "\\\\.\\pipe\\" {
		socketSplit = strings.Split(c.Agent.Socket, "\\")
		unixConfig.Agent.Socket = "/tmp/" + socketSplit[len(socketSplit)-1] + ".sock"
	} else {
		socketSplit = strings.Split(c.Agent.Socket, "/")
		winConfig.Agent.Socket = "\\\\\\\\.\\\\pipe\\\\" + socketSplit[len(socketSplit)-1]
	}

	// Create the config files for each platform
	err = createOutputConfigFile("./deploy/windows/sshrimp-windows.toml", winConfig)
	if err != nil {
		return err
	}
	err = createOutputConfigFile("./deploy/mac/sshrimp-mac.toml", unixConfig)
	if err != nil {
		return err
	}
	err = createOutputConfigFile("./deploy/linux/sshrimp-linux.toml", unixConfig)
	if err != nil {
		return err
	}

	// Instantiate the archiver, ensuring overwrite is enabled
	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      true,
		ImplicitTopLevelFolder: false,
	}
	return z.Archive([]string{"./deploy/windows", "./deploy/mac", "./deploy/linux"}, "./deploy.zip")
}

// Function to create a new output TOML file with name filename and using config
// returns any errors
func createOutputConfigFile(fileName string, config SSHrimp) error {
	configFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer configFile.Close()

	// Encode the configuration values as a TOML file
	encoder := toml.NewEncoder(configFile)
	err = encoder.Encode(config)
	return err
}
