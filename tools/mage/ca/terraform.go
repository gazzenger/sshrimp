package ca

import (
	"encoding/json"
	"strconv"

	"github.com/gazzenger/sshrimp/internal/config"
)

// Provider describes an AWS provider
type Provider struct {
	Version           string   `json:"version"`
	Alias             string   `json:"alias"`
	Region            string   `json:"region"`
	AllowedAccountIDs []string `json:"allowed_account_ids"`
}

// Module describes a terraform module
type Module struct {
	Source    string            `json:"source"`
	Providers map[string]string `json:"providers"`
}

// TerraformOutput represents the main.tf.json struct
type TerraformOutput struct {
	Provider map[string][]Provider `json:"provider"`
	Module   map[string]Module     `json:"module"`
}

func generateTerraform(c *config.SSHrimp) ([]byte, error) {

	providers := make([]Provider, len(c.CertificateAuthority.Regions))
	modules := make(map[string]Module, len(c.CertificateAuthority.Regions))
	for index, region := range c.CertificateAuthority.Regions {
		providers[index].Version = "~> 2.49"
		providers[index].Alias = region
		providers[index].Region = region
		providers[index].AllowedAccountIDs = []string{
			strconv.Itoa(c.CertificateAuthority.AccountID),
		}
		modules["sshrimp-"+region] = Module{
			Source: "./terraform",
			Providers: map[string]string{
				"aws": "aws." + region,
			},
		}
	}

	output := TerraformOutput{
		Provider: map[string][]Provider{
			"aws": providers,
		},
		Module: modules,
	}

	return json.MarshalIndent(output, "", "  ")
}

func generateVariableDefinitionsFile(c *config.SSHrimp) []byte {

	webidentityPrincipaIdentifiers := "arn:aws:iam::" + strconv.Itoa(c.CertificateAuthority.AccountID) + ":oidc-provider/" + c.CertificateAuthority.IdentityProviderURI
	webidentityProviderUrl := c.CertificateAuthority.IdentityProviderURI + ":aud"
	webidentityClientId := c.CertificateAuthority.IdentityProviderClientID

	outputString := "variable \"webidentity_principal_identifiers\" {\n" +
		"  type    = string\n" +
		"  default = \"" + webidentityPrincipaIdentifiers + "\"\n" +
		"}\n" +
		"variable \"webidentity_provider_url\" {\n" +
		"  type    = string\n" +
		"  default = \"" + webidentityProviderUrl + "\"\n" +
		"}\n" +
		"variable \"webidentity_client_id\" {\n" +
		"  type    = string\n" +
		"  default = \"" + webidentityClientId + "\"\n" +
		"}"

	outputArray := []byte(outputString)

	return outputArray
}
