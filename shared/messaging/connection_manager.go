package messaging

import (
	"errors"
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	ErrConnectionNotFound = errors.New("connection not found")
)

// connWrapper is wrapper around websocket to allow thread-safe operation
// remmeber websocket connection is not thread safe
type connWrapper struct {
	conn  *websocket.Conn // its from gurilla package ?
	mutex sync.Mutex
}

type ConnManager struct {
	connections map[string]*connWrapper
	mutex       sync.RWMutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true //allow all origins
	},
}

// TODO
// if multiple instanceof connectionMnager exist , you should think of different shared storage !
// cool
func NewConnManger() *ConnManager {
	return &ConnManager{
		connections: make(map[string]*connWrapper, 0),
	}
}

func (cm *ConnManager) Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (cm *ConnManager) Add(id string, conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.connections[id] = &connWrapper{
		conn: conn,
	}

	log.Printf("conn added for user %s", id)
}

func (cm *ConnManager) Remove(id string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	delete(cm.connections, id)
}

func (cm *ConnManager) Get(id string) (*websocket.Conn, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	wrapper, ok := cm.connections[id]
	if !ok {
		return nil, false
	}
	return wrapper.conn, true
}

func (cm *ConnManager) SendMessage(id string, message contracts.WSMessage) error {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	wrapper, ok := cm.connections[id]
	if !ok {
		return ErrConnectionNotFound
	}

	wrapper.mutex.Lock()
	defer wrapper.mutex.Unlock()

	return wrapper.conn.WriteJSON(message)
}
