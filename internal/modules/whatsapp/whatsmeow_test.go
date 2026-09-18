package whatsapp

import (
	"testing"
	"time"

	"go.mau.fi/whatsmeow"
)

func TestParseJID(t *testing.T) {
	tests := []struct {
		name string
		jid  string
		ok   bool
	}{
		{name: "phone number", jid: "+628123456789", ok: true},
		{name: "full JID", jid: "628123456789@s.whatsapp.net", ok: true},
		{name: "empty", jid: "", ok: false},
		{name: "non numeric user", jid: "not-a-phone", ok: false},
		{name: "empty user", jid: "@s.whatsapp.net", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := (&Whatsmeow{}).ParseJID(tt.jid)
			if ok != tt.ok {
				t.Fatalf("ParseJID(%q) ok = %v, want %v", tt.jid, ok, tt.ok)
			}
		})
	}
}

func TestSendKillChannelDoesNotBlock(t *testing.T) {
	w := &Whatsmeow{
		clientStore: NewStore[*whatsmeow.Client](),
		killchannel: NewStore[chan bool](),
	}
	w.killchannel.Set("unbuffered", make(chan bool))

	done := make(chan struct{})
	go func() {
		w.SendKillChannel("unbuffered")
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("SendKillChannel blocked on an unbuffered channel")
	}
}

func TestNewKillChannelIsBufferedAndReusable(t *testing.T) {
	w := &Whatsmeow{
		clientStore: NewStore[*whatsmeow.Client](),
		killchannel: NewStore[chan bool](),
	}
	w.NewKillChannel("device")

	channel, ok := w.killchannel.Get("device")
	if !ok {
		t.Fatal("kill channel was not created")
	}
	if cap(channel) != 1 {
		t.Fatalf("kill channel capacity = %d, want 1", cap(channel))
	}
	w.SendKillChannel("device")
	select {
	case <-channel:
	default:
		t.Fatal("kill signal was not delivered")
	}
}
