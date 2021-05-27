// +build windows
package main

import (
	"net"

	"github.com/Microsoft/go-winio"
)

func InitPipeListener(socketAddress string, err error) (net.Listener, error) {

	// setup named pipe for use with Windows OpenSSH
	namedPipeFullName := socketAddress
	var cfg = &winio.PipeConfig{}
	return winio.ListenPipe(namedPipeFullName, cfg)

}
