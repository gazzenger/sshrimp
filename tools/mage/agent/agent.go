package agent

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/magefile/mage/sh"
)

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
func PackageFiles() error {

	return RecursiveZip("./deploy", "./deploy.zip")
}

// Function to recursively zip a folder
// taken from https://stackoverflow.com/a/49057861
func RecursiveZip(pathToZip, destinationPath string) error {
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	err = filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))
		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	return nil
}
