package websocket

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

func Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header,
	readBufSize, writeBufSize int) (*websocket.Conn, error) {
	u := websocket.Upgrader{HandshakeTimeout: time.Duration(time.Second * 5), ReadBufferSize: readBufSize, WriteBufferSize: writeBufSize}
	u.Error = func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		// don't return errors to maintain backwards compatibility
	}
	u.CheckOrigin = func(r *http.Request) bool {
		// allow all connections by default
		return true
	}
	ws, err := u.Upgrade(w, r, responseHeader)
	if _, ok := err.(websocket.HandshakeError); ok {
		return nil, errors.New("not a websocket handshake")
	} else if err != nil {
		return nil, errors.New("Cannot setup WebSocket connection:" + fmt.Sprintf("%s", err))
	}
	return ws, nil
}
