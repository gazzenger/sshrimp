// +build windows

package main

import (
	"fmt"
	"net"
	"os"

	"github.com/Microsoft/go-winio"
)

const (
	AgentMaxMessageLength = 1<<14 - 1
)

// func GetSystemSecurityDescriptor() string {
//
// SDDL encoded.
//
// (system = SECURITY_NT_AUTHORITY | SECURITY_LOCAL_SYSTEM_RID)
// owner: system
// grant: GENERIC_ALL to system
//
// return "O:SYD:(A;;GA;;;SY)"
// return "S:(ML;;NW;;;LW)D:(A;;0x12019f;;;WD)"
// }

func InitListener(socketAddress string, err error) (net.Listener, error) {
	if len(socketAddress) > 9 && socketAddress[0:9] == "\\\\.\\pipe\\" {
		// setup named pipe for use with Windows OpenSSH
		namedPipeFullName := socketAddress
		var cfg = &winio.PipeConfig{
			// SecurityDescriptor: GetSystemSecurityDescriptor(),
		}
		return winio.ListenPipe(namedPipeFullName, cfg)
	} else {
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
}

// QueryAgent provides a way to query the named windows openssh agent pipe
func QueryAgent(pipeName string, buf []byte) (result []byte, err error) {
	return
}
