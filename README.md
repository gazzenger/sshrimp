# sshrimp 🦐

SSH Certificate Authority in a lambda, automated by an OpenID Connect enabled agent.

Why? Check out this presentation [Zero Trust SSH - linux.conf.au 2020](http://youtu.be/lYzklWPTbsQ).

## ~~ Warning ~~

This is still in very early development. Only use for testing. Not suitable for use in production yet. PR's welcome ;)

## Quickstart

This project uses [mage](https://magefile.org/) as a build tool. Install it.

Build the agent, lambda, and generate terraform code ready for deployment:

    mage

## Deployment

[Terraform](https://www.terraform.io/) files are defined in `/terraform` and the generated `sshrimp-ca.tf.json` file can be used to automatically deploy sshrimp into multiple AWS regions.

    terraform init
    terraform apply

> You will need [AWS credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html) in your environment to run `terraform apply`. You can also use [aws-vault](https://github.com/99designs/aws-vault) or [aws-oidc](https://github.com/stoggi/aws-oidc) to more securely manage AWS credentials on the command line.


## sshd_config (on your server)

Server configruation is minimal. Get the public keys from KMS (using AWS credentials):

    mage ca:keys

Put these keys in a file on your server `/etc/ssh/trusted_user_ca_keys`, owned by `root` permissions `0644`.

Modify `/etc/ssh/sshd_config` to add the line:

    TrustedUserCAKeys /etc/ssh/trusted_user_ca_keys


## ssh_config (on your local computer)

Since OpenSSH (>= 7.3), you can use the [IdentityAgent](https://man.openbsd.org/ssh_config.5#IdentityAgent) option in your ssh config file to set the socketname you configured:

    Host *.sshrimp.io
        User jeremy
        IdentityAgent /tmp/sshrimp-agent.sock

This has the advantage of only using the agent for the group of hosts you need, and let other hosts use your regular agent (like github.com for cloning git repos). In fact, you can't add other identities to the sshrimp-agent. It's meant to be used for only the hosts you need it for.

> For other SSH clients or older versions, set the `SSH_AUTH_SOCK` environment variable when invoking ssh: `SSH_AUTH_SOCK=/tmp/sshrimp-agent.sock ssh user@host`

## Let's go!

Start the agent:

    sshrimp-agent /path/to/sshrimp.toml

SSH to your host:

    ssh example.server.sshrimp.io

🎉

## Why sshrimp?

* Shrimp have shells.
* Shrimp are lightweight.
* Has a [backronym](https://en.wikipedia.org/wiki/Backronym): SSH. Really. Isn't. My. Problem.
* Shrimp on a barbie?
* Yeah...


## Usage on Windows

On Windows - you can use WSL and use the server agent exactly the same as on Linux
or
If using OpenSSH Client (installed via Windows Features) this currently only supports Pipes, therefore to get this working, configure the socket field in the config file to be 
```
Socket = "\\\\.\\pipe\\sshrimp"
```
And this can either be referenced by doing the following
```cmd
set SSH_AUTH_SOCK=\\.\pipe\ssh-pageant
ssh username@ipaddress
```
or setup the config file 
```cmd
Host [HOSTNAME]
    HostName [IPADDRESS]
    User [USERNAME]
    IdentityAgent \\.\pipe\sshrimp
    ForwardAgent yes
```
Then run the ssh command
```cmd
ssh [HOSTNAME]
```

## Deployment Config File
Some of the default config settings are listed below
```toml
[Agent] 
  ProviderURL = ""
  ClientID = ""
  BrowserCommand = [""]
  Socket = ""

[CertificateAuthority]
  AccountID = XXXXX
  Regions = [""]
  FunctionName = "sshrimp"
  KeyAlias = "alias/sshrimp"
  ForceCommandRegex = "^$"
  SourceAddressRegex = ""
  UsernameRegex = "^(.*)@someemail\\.com"
  UsernameClaim = "email"
  ValidAfterOffset = "-5m"
  ValidBeforeOffset = "+2m"
  Extensions = ["permit-agent-forwarding", "permit-port-forwarding", "permit-pty", "permit-user-rc", "permit-x11-forwarding"]
```
> Please note the Browser command is no longer needed, with the updated AWS-OIDC, as the default browser can be utilised.

> Also Socket can be defined as a UNIX socket file, OR a Windows Named Pipe

## Client Distribution

To build the SSHrimp-Agent for distribution across Windows, Linux and Mac, run the following mage command.
*Please note this only runs on Windows, as the Windows Build requires packages only found on Windows*

    mage buildandpackage
> This will build the sshrimp-agent for Windows, Mac and Linux and place these in the deploy folder, and zip it up


The output zip file is ```deploy.zip```, and contains the built sshrimp-agent executable for Windows, Mac and Linux in separate folders, with the config files, and deploy scripts.

The config file contains bare minimum parameters for clients, these are shown below.
```toml
[Agent] 
  ProviderURL = ""
  ClientID = ""
  BrowserCommand = [""]
  Socket = ""

[CertificateAuthority]
  AccountID = 11111111
  Regions = ["aaaaaaaa"]
```
## Code Sources
A thanks for help from
- Main body of this repository is forked from Jeremy Stott's SSHrimp project - https://github.com/stoggi/sshrimp - MIT
- Usage with Pipes on Windows - https://github.com/benpye/wsl-ssh-pageant - BSD 2-Clause
- Minor updates to Serve Agent for allowing multiple connections - https://github.com/daveadams/vaulted/blob/56a9a631ececd4610d83d6499725b34d64285ccc/lib/proxy_keyring.go#L82 - MIT
- Recursively zip folder - https://stackoverflow.com/a/49057861
- Create a launch agent on Mac - https://stackoverflow.com/questions/6442364/running-script-upon-login-mac


## TODO
* Connect with a provisioning user
* Make zipping deploy.zip a .tar.gz to maintain ACL