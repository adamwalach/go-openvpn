package mi_test

import (
	"bufio"
	"net"
	"strings"
	"testing"

	mi "github.com/adamwalach/go-openvpn/server/mi"
	"github.com/stretchr/testify/assert"
)

var cResponsePid = `SUCCESS: pid=10869
`

func TestReadResponsePid(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader(cResponsePid))
	response, err := mi.ReadResponse(reader)
	assert.Nil(t, err)
	pid, err := mi.ParsePid(response)
	assert.Nil(t, err)
	assert.Equal(t, int64(10869), pid)
}

func TestReadResponseFailure(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader(""))
	_, err := mi.ReadResponse(reader)
	assert.NotNil(t, err)
}

func TestSendCommandFailure(t *testing.T) {
	conn := &net.IPConn{}
	err := mi.SendCommand(conn, "dummy")
	assert.NotNil(t, err)
}
