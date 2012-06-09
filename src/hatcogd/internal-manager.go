package main

import "net"

type InternalManager struct {
	port        string
	connections []*Internal
	fromUser    chan Message
	Nick        string // Need to know, to tell client
	lastPrivate []byte // Most recent private message
}

type Message struct {
	channel string
	content string
}

func NewInternalManager(port string, fromUser chan Message) *InternalManager {

	var connections = make([]*Internal, 0)
	return &InternalManager{port, connections, fromUser, "", nil}
}

// Act as a server, forward data to irc connection
func (self *InternalManager) Run() {

	var listener net.Listener
	var netConn net.Conn
	var internalConn *Internal
	var err error

	listener, err = net.Listen("tcp", "127.0.0.1:"+self.port)

	if err != nil {
		LOG.Fatal("Error on internal listen: " + err.Error())
	}
	defer listener.Close()

	for {
		netConn, err = listener.Accept()
		if err != nil {
			LOG.Fatal("Listener accept error: " + err.Error())
			break
		}

		internalConn = &Internal{
			netConn: netConn,
			channel: "",
			manager: self}
		self.connections = append(self.connections, internalConn)
		go internalConn.Run()
	}

}

// Write a message to channel connection
func (self *InternalManager) WriteChannel(channel string, msg []byte) (int, error) {

	var bytesWritten int

	for _, conn := range self.connections {
		if conn.channel == channel {
			conn.netConn.Write(msg)
			bytesWritten += len(msg)
		}
	}

	return bytesWritten, nil
}

// Write a message to all client connections
func (self *InternalManager) WriteAll(msg []byte) (int, error) {

	var bytesWritten int
	for _, conn := range self.connections {
		conn.netConn.Write(msg)
		bytesWritten += len(msg)
	}

	return bytesWritten, nil
}

// Remove a connection from our list, probably because user closed it
func (self *InternalManager) delete(internalConn *Internal) {
	newConnections := make([]*Internal, 0, len(self.connections)-1)
	for _, conn := range self.connections {
		if conn.channel != internalConn.channel {
			newConnections = append(newConnections, conn)
		}
	}
	self.connections = newConnections
}

// The internal connection for given channel, or nil
func (self *InternalManager) GetChannelConnection(channel string) *Internal {
	for _, conn := range self.connections {
		if conn.channel == channel {
			return conn
		}
	}
	return nil
}

// Do we have a connection (a client) open on given channel or nick
func (self *InternalManager) HasChannel(channel string) bool {
	return self.GetChannelConnection(channel) != nil
}

func (self *InternalManager) Close() error {
	for _, conn := range self.connections {
		conn.netConn.Close()
	}
	return nil
}
