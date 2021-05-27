package identity

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/stoggi/sshrimp/internal/config"
)

// Identity holds information required to verify an OIDC identity token
type Identity struct {
	ctx           context.Context
	verifier      *oidc.IDTokenVerifier
	usernameRE    *regexp.Regexp
	usernameClaim string
}

// NewIdentity return a new Identity, with default values and oidc proivder information populated
func NewIdentity(c *config.SSHrimp) (*Identity, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, c.Agent.ProviderURL)
	if err != nil {
		return nil, err
	}

	oidcConfig := &oidc.Config{
		ClientID:             c.Agent.ClientID,
		SupportedSigningAlgs: []string{"RS256"},
	}

	return &Identity{
		ctx:           ctx,
		verifier:      provider.Verifier(oidcConfig),
		usernameRE:    regexp.MustCompile(c.CertificateAuthority.UsernameRegex),
		usernameClaim: c.CertificateAuthority.UsernameClaim,
	}, nil
}

// Validate an identity token
func (i *Identity) Validate(token string) (string, []string, error) {

	idToken, err := i.verifier.Verify(i.ctx, token)
	if err != nil {
		return "", []string{}, errors.New("failed to verify identity token: " + err.Error())
	}

	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		return "", []string{}, errors.New("failed to parse claims: " + err.Error())
	}

	claimedUsername, ok := claims[i.usernameClaim].(string)
	if !ok {
		return "", []string{}, errors.New("configured username claim not in identity token")
	}

	return i.parseUsername(claimedUsername, splitRoles(claims["roles"]))
}

func (i *Identity) parseUsername(username string, roles []string) (string, []string, error) {
	if match := i.usernameRE.FindStringSubmatch(username); match != nil {
		return match[1], roles, nil
	}
	return "", []string{}, errors.New("unable to parse username from claim")
}

// Split the roles array from the roles field in the claims
func splitRoles(roles interface{}) []string {
	roleStr := fmt.Sprintf("%v", roles)
	roleStr = strings.ReplaceAll(roleStr, "[", "")
	roleStr = strings.ReplaceAll(roleStr, "]", "")
	return strings.Split(roleStr, " ")
}
