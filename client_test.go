package source

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	s := newServer(t)
	if s == nil {
		return
	}
	defer func() {
		assert.NoError(t, s.Close())
	}()

	c, err := NewClient(s.Addr, Timeout(time.Second*2))
	if !assert.NoError(t, err) {
		return
	}

	defer func() {
		assert.NoError(t, c.Close())
	}()

	_, err = c.Exec("status")
	assert.NoError(t, err)

	_, err = c.ExecCmd(NewCmd("echo").WithArgs("my test"))
	assert.NoError(t, err)

	resp, err2 := c.Exec("invalid")
	assert.NoError(t, err2)
	assert.Equal(t, "unknown command 2:invalid", resp)
}

func TestClientNilOption(t *testing.T) {
	_, err := NewClient("", nil)
	if !assert.Error(t, err) {
		return
	}

	assert.Equal(t, ErrNilOption, err)
}

func TestClientOptionError(t *testing.T) {
	errBadOption := errors.New("bad option")
	_, err := NewClient("", func(c *Client) error { return errBadOption })
	if !assert.Error(t, err) {
		return
	}

	assert.Equal(t, errBadOption, err)
}

func TestClientDisconnect(t *testing.T) {
	s := newServer(t)
	if s == nil {
		return
	}
	defer func() {
		assert.NoError(t, s.Close())
	}()

	c, err := NewClient(s.Addr, Timeout(time.Second*2))
	if !assert.NoError(t, err) {
		return
	}

	assert.NoError(t, c.Close())

	_, err = c.Exec("version")
	assert.Error(t, err)
}

func TestClientWriteFail(t *testing.T) {
	s := newServer(t)
	if s == nil {
		return
	}
	defer func() {
		assert.NoError(t, s.Close())
	}()

	c, err := NewClient(s.Addr, Timeout(time.Second*2))
	if !assert.NoError(t, err) {
		return
	}
	assert.NoError(t, c.conn.(*net.TCPConn).CloseWrite())

	_, err = c.Exec("version")
	assert.Error(t, err)
}

func TestClientDialFail(t *testing.T) {
	c, err := NewClient("127.0.0.1", Timeout(time.Nanosecond))
	if assert.Error(t, err) {
		return
	}

	// Should never get here
	assert.NoError(t, c.Close())
}

func TestClientFailConn(t *testing.T) {
	s := newServerStopped(t)
	if s == nil {
		return
	}
	s.failConn = true
	s.Start()
	defer func() {
		assert.NoError(t, s.Close())
	}()

	c, err := NewClient(s.Addr, Timeout(time.Second))
	if !assert.NoError(t, err) {
		return
	}

	assert.NoError(t, c.Close())
}
