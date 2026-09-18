package pubsub

import (
	"sapasora/platform/config"

	messagebus "github.com/vardius/message-bus"
	"go.uber.org/fx"
)

type PubSub struct {
	bus messagebus.MessageBus
}

type PubSubIn struct {
	fx.In

	Config config.Config
}

func New(in PubSubIn) *PubSub {
	size := in.Config.GetInt("pubsub.size", 100)
	return &PubSub{
		bus: messagebus.New(size),
	}
}

func (p *PubSub) Subscribe(topic string, fn any) error {
	return p.bus.Subscribe(topic, fn)
}

func (p *PubSub) Unsubscribe(topic string, fn any) error {
	return p.bus.Unsubscribe(topic, fn)
}

func (p *PubSub) Publish(topic string, args ...any) {
	p.bus.Publish(topic, args...)
}

func (p *PubSub) Close(topic string) {
	p.bus.Close(topic)
}
