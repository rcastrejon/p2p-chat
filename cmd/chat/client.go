package chat

import (
	"net"

	"github.com/rcastrejon/p2p-chat/pb"
	"google.golang.org/protobuf/proto"
)

type ChatClient struct {
	conn net.PacketConn
	peer *net.UDPAddr
	sub  chan *pb.Message
}

func NewChatClient(peerAddress string, localPort string) (*ChatClient, error) {
	// Resolve the peer address, asume that the peer is already listening. The
	// peer address should be in the form of "host:port".
	peer, err := net.ResolveUDPAddr("udp4", peerAddress)
	if err != nil {
		return nil, err
	}

	// Create a listen packet connection. This will listen for incoming udp
	// packets on the specified port.
	conn, err := net.ListenPacket("udp4", ":"+localPort)
	if err != nil {
		return nil, err
	}

	c := &ChatClient{
		conn: conn,
		peer: peer,
		sub:  make(chan *pb.Message),
	}

	// Start a goroutine to listen for incoming messages without blocking the
	// main thread.
	go c.listenForMessages()

	return c, nil
}

// Listens for udp packets and sends them to the pb.Message channel. This is a
// blocking function, so it is meant to be run in a goroutine.
func (c *ChatClient) listenForMessages() {
	buf := make([]byte, 1024)
	for {
		// TODO: consider handling the address of the sender
		n, _, _ := c.conn.ReadFrom(buf)

		// We want to avoid blocking the listener, so we CAN spawn a goroutine to
		// handle the message (probably not necessary for our small program).
		m := pb.Message{}
		err := proto.Unmarshal(buf[:n], &m)
		if err != nil {
			continue
		}
		c.sub <- &m
	}
}

func (c *ChatClient) Close() {
	c.conn.Close()
}

func (c *ChatClient) Send(message *pb.Message) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	_, err = c.conn.WriteTo(data, c.peer)
	return err
}

func (c *ChatClient) Receive() *pb.Message {
	return <-c.sub
}
