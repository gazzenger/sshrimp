// +build darwin linux

package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

func init() {
	sigs = append(sigs, syscall.SIGTERM, syscall.SIGHUP)
}

func InitListener(socketAddress string, err error) (net.Listener, error) {
	// testing to ensure nothing else is using the AF_UNIX domain socket file
	// only used on unix systems, or using WSL
	if _, err = os.Stat(socketAddress); err == nil {
		conn, sockErr := net.Dial("unix", socketAddress)
		if sockErr == nil { // socket is accepting connections
			conn.Close()
			return nil, fmt.Errorf("socket %s already exists", socketAddress)
		}
		os.Remove(socketAddress) // socket is not accepting connections, assuming safe to remove
	}
	// setup AF_UNIX domain socket for use with Linux, Mac or WSL
	// This affects all files created for the process. Since this is a sensitive
	// socket, only allow the current user to write to the socket.
	// syscall.Umask(0077)
	return net.Listen("unix", socketAddress)
}
