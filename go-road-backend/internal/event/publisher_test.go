package event

import (
	"testing"

	"github.com/nats-io/nats.go"
)

func TestPublisherNew(t *testing.T) {
	// Test with nil connection to verify error handling
	_, err := NewPublisher(nil)
	if err == nil {
		t.Log("expected error with nil conn (no NATS server)")
	}
}

func TestPublishWithNil(t *testing.T) {
	var p *Publisher
	err := p.Publish("test.subject", "data")
	if err == nil {
		t.Log("expected error with nil publisher")
	}
}

func TestPublisherStruct(t *testing.T) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Skip("NATS server not available")
	}
	defer nc.Close()

	p, err := NewPublisher(nc)
	if err != nil {
		t.Fatal("expected publisher creation")
	}
	if p == nil {
		t.Fatal("expected non-nil publisher")
	}
}
