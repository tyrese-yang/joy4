package rtmp

import (
	"context"
	"net"
	"time"

	quic "github.com/lucas-clemente/quic-go"
)

type qconn struct {
	session quic.Session

	receiveStream quic.Stream
	sendStream    quic.Stream
}

func newConn(sess quic.Session) (*qconn, error) {
	stream, err := sess.OpenStream()
	if err != nil {
		return nil, err
	}
	return &qconn{
		session:       sess,
		sendStream:    stream,
		receiveStream: stream,
	}, nil
}

func (c *qconn) Read(b []byte) (int, error) {
	if c.receiveStream == nil {
		var err error
		c.receiveStream, err = c.session.AcceptStream(context.Background())
		if err != nil {
			return 0, err
		}

		// quic.Stream.Close() closes the stream for writing
		err = c.receiveStream.Close()
		if err != nil {
			return 0, err
		}
	}

	return c.receiveStream.Read(b)
}

func (c *qconn) Write(b []byte) (int, error) {
	return c.sendStream.Write(b)
}

// LocalAddr returns the local network address.
// needed to fulfill the net.Conn interface
func (c *qconn) LocalAddr() net.Addr {
	return c.session.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *qconn) RemoteAddr() net.Addr {
	return c.session.RemoteAddr()
}

func (c *qconn) Close() error {
	return c.session.CloseWithError(0, "")
}

func (c *qconn) SetDeadline(t time.Time) error {
	return nil
}

func (c *qconn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *qconn) SetWriteDeadline(t time.Time) error {
	return nil
}

var _ net.Conn = &qconn{}

type qserver struct {
	quicServer quic.Listener
}

var _ net.Listener = &qserver{}

// Accept waits for and returns the next connection to the listener.
func (s *qserver) Accept() (net.Conn, error) {
	sess, err := s.quicServer.Accept(context.Background())
	if err != nil {
		return nil, err
	}
	qconn, err := newConn(sess)
	if err != nil {
		return nil, err
	}
	return qconn, nil
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (s *qserver) Close() error {
	return s.quicServer.Close()
}

// Addr returns the listener's network address.
func (s *qserver) Addr() net.Addr {
	return s.quicServer.Addr()
}
