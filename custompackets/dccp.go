package custompackets

/*
#include <sys/socket.h>
#include <netinet/in.h>
*/
import "C"

import (
	"net"
	"syscall"
)

func SendDCCPPacket(dstIP net.IP, dstPort int, message []byte) error {
	SOCK_DCCP := 0x6 // https://datatracker.ietf.org/doc/rfc5596/

	sock, err := syscall.Socket(syscall.AF_INET, SOCK_DCCP, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(sock)

	serverAddr := syscall.SockaddrInet4{
		Port: dstPort,
		Addr: [4]byte{dstIP[0], dstIP[1], dstIP[2], dstIP[3]},
	}

	err = syscall.Connect(sock, &serverAddr)
	if err != nil {
		return err
	}

	_, err = syscall.Write(sock, message)
	if err != nil {
		return err
	}

	return nil
}
