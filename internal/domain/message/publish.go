package message

import "context"

type Publish struct {
	publisher Publisher
	cache     Cache
}

func (p Publish) Handle(ctx context.Context, m *Message) error {
	switch {
	case p.shouldSendFailed(m):
		if err := p.send(ctx, m.Message()); err != nil {
			return err
		}
		p.cache.Set(m.resource)
		return nil
	case p.shouldSendReadyAgain(m):
		if err := p.send(ctx, m.Message()); err != nil {
			return err
		}
		p.cache.Remove(m.resource)
		return nil
	default:
		return nil
	}
}

func (p Publish) shouldSendReadyAgain(m *Message) bool {
	return p.cache.Exists(m.resource) && !m.IsFailed()
}

func (p Publish) shouldSendFailed(m *Message) bool {
	return !p.cache.Exists(m.resource) && m.IsFailed()
}

func (p Publish) send(ctx context.Context, message string) error {
	return p.publisher.Publish(ctx, message)
}

func NewPublish(publisher Publisher, cache Cache) *Publish {
	return &Publish{publisher: publisher, cache: cache}
}
