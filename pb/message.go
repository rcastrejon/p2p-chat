package pb

import "google.golang.org/protobuf/types/known/timestamppb"

func NewMessage(alias string, body string) *Message {
	return &Message{
		Alias:     alias,
		Body:      body,
		Timestamp: timestamppb.Now(),
	}
}
