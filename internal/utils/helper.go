package utils

import (
	"encoding/binary"
	"io"
	"net"
)

// ReadMessage reads a message from the connection.
func ReadMessage(conn net.Conn) ([]byte, error) {
	var length uint64
	if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
		return nil, err
	}

	msg := make([]byte, length)
	_, err := io.ReadFull(conn, msg)
	return msg, err
}

// WriteMessage writes a message to the connection.
func WriteMessage(conn net.Conn, msg []byte) error {
	if err := binary.Write(conn, binary.BigEndian, uint64(len(msg))); err != nil {
		return err
	}
	_, err := conn.Write(msg)
	return err
}
