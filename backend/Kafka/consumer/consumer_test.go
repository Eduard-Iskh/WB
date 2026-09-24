package consumer

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestConsumer(t *testing.T) {
	testCases := []struct {
		name    string
		payload []byte
	}{
		{
			name:    "TestEmptyPayload",
			payload: []byte(``),
		},
		{
			name:    "TestNoValidJSON",
			payload: []byte(`{`),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mess := DLQMessage{
				Payload: tt.payload,
			}
			value, err := json.Marshal(mess)
			if err != nil {
				t.Fatalf("json.Marshal(%q): %v", tt.payload, err)
			}

			var got *DLQMessage
			err = json.Unmarshal(value, &got)
			if err != nil {
				t.Fatalf("Ошибка конвертирования: %v", err)
			}
			if !bytes.Equal(got.Payload, mess.Payload) {
				t.Fatalf(
					"payload mismatch after JSON round trip: got %q, want %q",
					got,
					tt.payload,
				)
			}
		})
	}
}
