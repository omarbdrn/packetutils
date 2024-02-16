package packetutils

import "time"

type Connection interface {
	Write([]byte) (int, error)
	Read([]byte) (int, error)
	SetWriteDeadline(time.Time) error
	SetReadDeadline(time.Time) error
	Close() error
}
