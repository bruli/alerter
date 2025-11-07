package message

import "context"

type Publish struct {
	publisher Publisher
	cache     Cache
}

func (p Publish) Handle(ctx context.Context, m *Message) error {
	if !p.cache.Exists(m.resource) && m.IsFailed() {
		if err := p.publisher.Publish(ctx, m); err != nil {
			return err
		}
		p.cache.Set(m.resource)
	}
	return nil
}

func NewPublish(publisher Publisher, cache Cache) *Publish {
	return &Publish{publisher: publisher, cache: cache}
}
