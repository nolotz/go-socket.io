package websocket

import (
	"github.com/nolotz/go-socket.io/engineio/transport"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebsocketSetReadDeadline(t *testing.T) {
	at := assert.New(t)

	tran := &Transport{}
	conn := make(chan transport.Conn, 1)
	handler := func(w http.ResponseWriter, r *http.Request) {
		c, err := tran.Accept(w, r)
		require.NoError(t, err)

		conn <- c
	}
	httpSvr := httptest.NewServer(http.HandlerFunc(handler))
	defer httpSvr.Close()

	u, err := url.Parse(httpSvr.URL)
	require.NoError(t, err)

	u.Scheme = "ws"

	header := make(http.Header)
	cc, err := tran.Dial(u, header)
	require.NoError(t, err)

	defer cc.Close()

	sc := <-conn
	defer sc.Close()

	err = cc.SetReadDeadline(time.Now().Add(time.Second / 10))
	require.NoError(t, err)

	_, _, _, err = cc.NextReader()
	require.Error(t, err)

	timeout, ok := err.(net.Error)
	at.True(ok)
	at.True(timeout.Timeout())

	op, ok := err.(net.Error)
	at.True(ok)
	at.True(op.Timeout())
}
