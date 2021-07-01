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
	github.com/gazzenger/aws-oidc v0.0.0-20210620061911-39a318848078
	github.com/gazzenger/winssh-pageant v0.0.0-20210701065358-b3cebc44b385
	github.com/jinzhu/copier v0.3.2
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/lxn/win v0.0.0-20191128105842-2da648fda5b4
	github.com/magefile/mage v1.11.0
	github.com/mholt/archiver/v3 v3.5.0
	github.com/pkg/errors v0.9.1
	github.com/pquerna/cachecontrol v0.1.0 // indirect
	golang.org/x/crypto v0.0.0-20191128160524-b544559bb6d1
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	gopkg.in/square/go-jose.v2 v2.4.1 // indirect
)
