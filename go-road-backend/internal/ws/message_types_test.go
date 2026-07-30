package ws

import (
	"testing"
)

func TestMessageTypeConstants(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"MsgTypeLocation", MsgTypeLocation},
		{"MsgTypeChat", MsgTypeChat},
		{"MsgTypeHeartbeat", MsgTypeHeartbeat},
		{"MsgTypeJoinRoom", MsgTypeJoinRoom},
		{"MsgTypeLeaveRoom", MsgTypeLeaveRoom},
		{"MsgTypeEmergency", MsgTypeEmergency},
		{"MsgTypePresence", MsgTypePresence},
		{"MsgTypeError", MsgTypeError},
		{"MsgTypeNotification", MsgTypeNotification},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("%s should not be empty", tt.name)
			}
		})
	}
}
