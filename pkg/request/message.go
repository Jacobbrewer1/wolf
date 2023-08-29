package request

import "fmt"

// Message represents a message response.
type Message struct {
	Message string `json:"Message" xml:"Message"`
}

// NewMessage creates a new Message.
func NewMessage(message string, args ...any) *Message {
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(message, args...)
	} else {
		msg = message
	}
	return &Message{
		Message: msg,
	}
}

// MessageError represents a message response with an error. It is used when there is a message and error to return.
// An example of this is when trying to unmarshal a request body into a struct. If the request body is invalid, then the
// error will be returned to the client. If the request body is valid, but the struct is invalid, then the message will
// be returned to the client.
type MessageError struct {
	Message string `json:"Message" xml:"Message"`
	Error   string `json:"Error" xml:"Error"`
}

func NewMessageError(message string, err error) *MessageError {
	return &MessageError{
		Message: message,
		Error:   err.Error(),
	}
}
