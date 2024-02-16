package packetutils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"runtime"
	"time"

	"github.com/omarbdrn/eyepx/internal/packetutils/custompackets"
	"github.com/omarbdrn/eyepx/internal/packetutils/utils"
	"github.com/logrusorgru/aurora"
)

type Protocol string

const (
	TCP  Protocol = "tcp"
	UDP  Protocol = "udp"
	DCCP Protocol = "dccp"
)

func Send(conn Connection, data []byte, timeout time.Duration) error {
	err := conn.SetWriteDeadline(time.Now().Add(timeout))
	if err != nil {
		return &WriteTimeoutError{WrappedError: err}
	}
	length, err := conn.Write(data)
	if err != nil {
		return &WriteError{WrappedError: err}
	}
	if length < len(data) {
		return &WriteError{
			WrappedError: fmt.Errorf(
				"failed to write all bytes (%d bytes written, %d bytes expected)",
				length,
				len(data),
			),
		}
	}
	return nil
}

func Recv(conn Connection, timeout time.Duration) ([]byte, error) {
	response := make([]byte, 4096)
	err := conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return []byte{}, &ReadTimeoutError{WrappedError: err}
	}
	length, err := conn.Read(response)
	if err != nil {
		var netErr net.Error
		if (errors.As(err, &netErr) && netErr.Timeout()) { // timeout error
			return []byte{}, nil
		}
		return response[:length], &ReadError{
			Info:         hex.EncodeToString(response[:length]),
			WrappedError: err,
		}
	}
	return response[:length], nil
}

func SendRecv(conn Connection, data []byte, timeout time.Duration) ([]byte, error) {
	err := Send(conn, data, timeout)
	if err != nil {
		return []byte{}, err
	}
	return Recv(conn, timeout)
}

func SendPacket(proto Protocol, destIP net.IP, destPort int) error {
	switch proto {
	case TCP:
		
		conn, err := utils.DialTCP(destIP.String(), uint16(destPort))
		if err != nil{
			return err
		}
		defer conn.Close()

		err = Send(conn, []byte("EyePX-TCP"), time.Duration(10 * time.Second))
		if err != nil{
			return err
		}

		fmt.Printf("%s TCP Packet sent to %s:%d\n", aurora.Blue("[+]"), destIP, destPort)
		return nil
	case UDP:
		conn, err := utils.DialUDP(destIP.String(), uint16(destPort))
		if err != nil{
			return err
		}
		defer conn.Close()

		err = Send(conn, []byte("EyePX-UDP"), time.Duration(10 * time.Second))
		if err != nil{
			return err
		}

		fmt.Printf("%s UDP Packet sent to %s:%d\n", aurora.Blue("[+]"), destIP, destPort)
		return nil
	case DCCP:

		if runtime.GOOS == "darwin"{
			return errors.New("DCCP Protocol is not supported in MacOS idk")
		}
		
		err := custompackets.SendDCCPPacket(destIP, destPort, []byte("EyePX-DCCP"))
		if err != nil {
			return err
		}
		
		fmt.Printf("%s DCCP Packet sent to %s:%d\n", aurora.Blue("[+]"), destIP, destPort)
		return nil
	default:
		return errors.New("unsupported protocol")
	}
}
