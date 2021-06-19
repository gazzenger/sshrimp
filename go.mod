module github.com/gazzenger/sshrimp

go 1.13

replace github.com/b-b3rn4rd/gocfn => github.com/stoggi/gocfn v0.0.0-20200214083946-6202cea979b9

require (
	github.com/AlecAivazis/survey/v2 v2.0.5
	github.com/BurntSushi/toml v0.3.1
	github.com/Microsoft/go-winio v0.5.0
	github.com/alecthomas/kong v0.2.2
	github.com/aws/aws-lambda-go v1.13.3
	github.com/aws/aws-sdk-go v1.25.43
	github.com/awslabs/goformation/v4 v4.4.0
	github.com/coreos/go-oidc v2.1.0+incompatible
	github.com/gazzenger/aws-oidc v0.0.0-20210619080922-b29e6b03f557
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/magefile/mage v1.9.0
	github.com/pkg/errors v0.9.1
	github.com/pquerna/cachecontrol v0.1.0 // indirect
	golang.org/x/crypto v0.0.0-20191128160524-b544559bb6d1
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	gopkg.in/square/go-jose.v2 v2.4.1 // indirect
)
