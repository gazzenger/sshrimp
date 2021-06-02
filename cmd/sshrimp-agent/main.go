package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"net"
	"os"
	"os/signal"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/gazzenger/sshrimp/internal/config"
	"github.com/gazzenger/sshrimp/internal/sshrimpagent"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

var sigs = []os.Signal{os.Kill, os.Interrupt}

var cli struct {
	Config string `kong:"arg,type='string',help='sshrimp config file (default: ${config_file} or ${env_var_name} environment variable)',default='${config_file}',env='SSHRIMP_CONFIG'"`
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Name("sshrimp-agent"),
		kong.Description("An SSH Agent that renews SSH certificates automatically from a SSHrimp Certificate Authority."),
		kong.Vars{
			"config_file":  config.DefaultPath,
			"env_var_name": config.EnvVarName,
		},
	)

	c := config.NewSSHrimpWithDefaults()
	ctx.FatalIfErrorf(c.Read(cli.Config))
	ctx.FatalIfErrorf(launchAgent(c, ctx))
}

func launchAgent(c *config.SSHrimp, ctx *kong.Context) error {
	var (
		err        error
		listener   net.Listener
		privateKey crypto.Signer
		signer     ssh.Signer
	)

	// // testing to ensure nothing else is using the AF_UNIX domain socket file
	// // only used on unix systems, or using WSL
	// if _, err = os.Stat(c.Agent.Socket); err == nil {
	// 	conn, sockErr := net.Dial("unix", c.Agent.Socket)
	// 	if sockErr == nil { // socket is accepting connections
	// 		conn.Close()
	// 		return fmt.Errorf("socket %s already exists", c.Agent.Socket)
	// 	}
	// 	os.Remove(c.Agent.Socket) // socket is not accepting connections, assuming safe to remove
	// }

	// //on windows, without using WSL a pipe must be used instead

	// // setup AF_UNIX domain socket for use with Linux, Mac or WSL
	// // This affects all files created for the process. Since this is a sensitive
	// // socket, only allow the current user to write to the socket.
	// syscall.Umask(0077)
	// listener, err = net.Listen("unix", c.Agent.Socket)
	// if err != nil {
	// 	return err
	// }

	// // setup named pipe for use with Windows OpenSSH
	// namedPipeFullName := c.Agent.Socket
	// var cfg = &winio.PipeConfig{}
	// listener, err = winio.ListenPipe(namedPipeFullName, cfg)
	// if err != nil {
	// 	return err
	// }

	// usage of a pipe is taken from repository https://github.com/benpye/wsl-ssh-pageant
	if len(c.Agent.Socket) > 9 && c.Agent.Socket[0:9] == "\\\\.\\pipe\\" {
		listener, err = InitPipeListener(c.Agent.Socket, err)
	} else {
		listener, err = InitSocketListener(c.Agent.Socket, err)
	}
	if err != nil {
		return err
	}

	defer listener.Close()

	ctx.Printf("listening on %s", c.Agent.Socket)

	// Generate a new SSH private/public key pair
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	signer, err = ssh.NewSignerFromKey(privateKey)
	if err != nil {
		return err
	}

	// Create the sshrimp agent with our configuration and the private key signer
	sshrimpAgent := sshrimpagent.NewSSHrimpAgent(c, signer)

	// Listen for signals so that we can close the listener and exit nicely
	osSignals := make(chan os.Signal)
	signal.Notify(osSignals, sigs...)
	go func() {
		_ = <-osSignals
		listener.Close()
	}()

	// Accept connections and serve the agent
	for {
		var conn net.Conn
		conn, err = listener.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				// Occurs if the user interrupts the agent with a ctrl-c signal
				return nil
			}
			return err
		}

		go agent.ServeAgent(sshrimpAgent, conn)
		// if err = agent.ServeAgent(sshrimpAgent, conn); err != nil && !errors.Is(err, io.EOF) {
		// 	return err
		// }
	}
}
