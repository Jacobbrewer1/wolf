package discordgo

import "fmt"

type Operation int

const (
	// OperationDispatch - An event was dispatched.
	OperationDispatch Operation = 0

	// OperationHeartbeat - Fired periodically by the client to keep the connection alive.
	OperationHeartbeat Operation = 1

	// OperationIdentify - Starts a new session during the initial handshake.
	OperationIdentify Operation = 2

	// OperationPresenceUpdate - Update the client's presence.
	OperationPresenceUpdate Operation = 3

	// OperationVoiceStateUpdate - Used to join/leave or move between voice channels.
	OperationVoiceStateUpdate Operation = 4

	// OperationResume - Resume a previous session that was disconnected.
	OperationResume Operation = 6

	// OperationReconnect - You should attempt to reconnect and resume immediately.
	OperationReconnect Operation = 7

	// OperationRequestGuildMembers - Request information about offline guild members in a large guild.
	OperationRequestGuildMembers Operation = 8

	// OperationInvalidSession - The session has been invalidated. You should reconnect and identify/resume accordingly.
	OperationInvalidSession Operation = 9

	// OperationHello - Sent immediately after connecting, contains the heartbeat_interval to use.
	OperationHello Operation = 10

	// OperationHeartbeatACK - Sent in response to receiving a heartbeat to acknowledge that it has been received.
	OperationHeartbeatACK Operation = 11
)

// String returns a string representation of the Operation
func (o Operation) String() string {
	switch o {
	case OperationDispatch:
		return "Dispatch"
	case OperationHeartbeat:
		return "Heartbeat"
	case OperationIdentify:
		return "Identify"
	case OperationPresenceUpdate:
		return "Presence_Update"
	case OperationVoiceStateUpdate:
		return "Voice_State_Update"
	case OperationResume:
		return "Resume"
	case OperationReconnect:
		return "Reconnect"
	case OperationRequestGuildMembers:
		return "Request_Guild_Members"
	case OperationInvalidSession:
		return "Invalid_Session"
	case OperationHello:
		return "Hello"
	case OperationHeartbeatACK:
		return "Heartbeat_ACK"
	}
	return fmt.Sprintf("Unknown_Operation_(%d)", o)
}

// MarshalJSON returns the int code for the Operation
func (o *Operation) MarshalJSON() ([]byte, error) {
	i, err := Marshal(*o)
	if err != nil {
		return nil, fmt.Errorf("error marshaling operation: %w", err)
	}
	return i, nil
}

// UnmarshalJSON receives the int code from discord and converts it to the Operation
func (o *Operation) UnmarshalJSON(b []byte) error {
	var i int
	if err := Unmarshal(b, &i); err != nil {
		return err
	}
	*o = Operation(i)
	return nil
}
