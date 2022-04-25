package transport

import (
	"encoding/binary"
	"io"
	"net"

	"github.com/muque7/my-rpc/codec"
)

// Transport struct
type Transport struct {
	conn net.Conn
}

// NewTransport creates a transport
func NewTransport(conn net.Conn) *Transport {
	return &Transport{conn}
}

// Send data
func (t *Transport) Send(req codec.Data) error {
	b, err := codec.Encode(req) // Encode req into bytes
	if err != nil {
		return err
	}
	buf := make([]byte, 4+len(b))
	binary.BigEndian.PutUint32(buf[:4], uint32(len(b))) // Set Header field
	copy(buf[4:], b)                                    // Set Data field
	_, err = t.conn.Write(buf)
	return err
}

// Receive data
func (t *Transport) Receive() (codec.Data, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(t.conn, header)
	if err != nil {
		return codec.Data{}, err
	}
	dataLen := binary.BigEndian.Uint32(header) // Read Header filed
	data := make([]byte, dataLen)              // Read Data Field
	_, err = io.ReadFull(t.conn, data)
	// log.Println(data)
	if err != nil {
		return codec.Data{}, err
	}
	req, err := codec.Decode(data) // Decode rsp from bytes
	return req, err
}
