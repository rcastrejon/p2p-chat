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
	// Resolve the peer address
	peer, err := net.ResolveUDPAddr("udp4", peerAddress)
	if err != nil {
		return nil, err
	}

	// Create a packet connection
	conn, err := net.ListenPacket("udp4", ":"+localPort)
	if err != nil {
		return nil, err
	}

	c := &ChatClient{
		conn: conn,
		peer: peer,
		sub:  make(chan *pb.Message),
	}

	go c.listenForMessages()

	return c, nil
}

func (c *ChatClient) listenForMessages() {
	buf := make([]byte, 1024)
	for {
		n, _, err := c.conn.ReadFrom(buf)
		if err != nil {
			continue
		}
		m := pb.Message{}
		err = proto.Unmarshal(buf[:n], &m)
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
	b, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	_, err = c.conn.WriteTo(b, c.peer)
	return err
}

func (c *ChatClient) Receive() *pb.Message {
	return <-c.sub
}
