package ca

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
	"github.com/gazzenger/sshrimp/internal/config"
	"golang.org/x/crypto/ssh"
)

// Config Generate a sshrimp configuration file if it doesn't exsit
func Config() error {
	c := config.NewSSHrimpWithDefaults()

	// Read the existing config if it doesn't exist, generate a new one
	if err := c.Read(config.GetPath()); err != nil {
		configPath, err := config.Wizard(config.GetPath(), c)
		if err != nil {
			return err
		}
		// If a different config path is chosen, use it for the rest of the build
		os.Setenv("SSHRIMP_CONFIG", configPath)
	}

	return nil
}

// Build Builds the certificate authority
func Build() error {
	env := map[string]string{
		"GOOS": "linux",
	}
	return sh.RunWith(env, "go", "build", "./cmd/sshrimp-ca")
}

// Package Packages the certificate authority files into a zip archive
func Package() error {
	if modified, err := target.Path("sshrimp-ca", config.GetPath()); err == nil && !modified {
		return nil
	}

	mg.Deps(Build, Config)

	zipFile, err := os.Create("sshrimp-ca.zip")
	if err != nil {
		return err
	}
	defer zipFile.Close()

	if err := lambdaCreateArchive(zipFile, "sshrimp-ca", config.GetPath()); err != nil {
		return err
	}
	return nil
}

// Generate Generates a CloudFormation template used to deploy the certificate authority
func Generate() error {
	mg.Deps(Config)

	c := config.NewSSHrimp()
	if err := c.Read(config.GetPath()); err != nil {
		return err
	}

	if modified, err := target.Path("sshrimp-ca.tf.json", config.GetPath()); err == nil && modified {
		template, err := generateTerraform(c)
		if err != nil {
			return err
		}
		ioutil.WriteFile("sshrimp-ca.tf.json", template, 0644)
	}

	if modified, err := target.Path("./terraform/policy-variables.tf", config.GetPath()); err == nil && modified {
		// Generate policy variable file
		variableDefinitionsFile := generateVariableDefinitionsFile(c)
		ioutil.WriteFile("./terraform/policy-variables.tf", variableDefinitionsFile, 0644)
	}

	return nil
}

// Keys Get the public keys of all configured KMS keys in OpenSSH format
func Keys() error {

	c := config.NewSSHrimp()
	if err := c.Read(config.GetPath()); err != nil {
		return err
	}

	// For each configured region, get the public key from KMS and format it in an OpenSSH authorized_keys format
	for _, region := range c.CertificateAuthority.Regions {

		// Create a new session in the correct region
		session := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(region),
		}))
		svc := kms.New(session)

		// Get the public key from KMS
		response, err := svc.GetPublicKey(&kms.GetPublicKeyInput{
			KeyId: aws.String(c.CertificateAuthority.KeyAlias),
		})
		if err != nil {
			return err
		}

		// Parse the public key from KMS
		publicKey, err := x509.ParsePKIXPublicKey(response.PublicKey)
		if err != nil {
			return err
		}

		// Convert the public key into an SSH public key
		sshPublicKey, err := ssh.NewPublicKey(publicKey)
		if err != nil {
			return err
		}

		// Generate the final string to output on stdout
		authorizedKey := strings.TrimSuffix(string(ssh.MarshalAuthorizedKey(sshPublicKey)), "\n")
		fmt.Printf("%s sshrimp-ca@%s\n", authorizedKey, region)
	}

	return nil
}

// Clean Cleans the output files for sshrimp-ca
func Clean() error {
	if err := sh.Rm("sshrimp-ca"); err != nil {
		return err
	}
	if err := sh.Rm("sshrimp-ca.tf.json"); err != nil {
		return err
	}
	if err := sh.Rm("sshrimp-ca.zip"); err != nil {
		return err
	}
	return nil
}
