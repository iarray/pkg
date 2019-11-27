package eventbus

type Event interface {
	EventName() string
}

type EventStore interface {
	Save(e Event)
}

type EventHandler func(e Event)

type Subscription struct {
	Topic string
	ID    uint64
}

type Subscriber interface {
	Subscribe(topic string, handler EventHandler) Subscription
	Unsubscribe(info *Subscription)
}

type Publisher interface {
	Publish(event Event)
	PublishAsync(event Event)
}

type Bus interface {
	Subscriber
	Publisher
	CommandSubscriber
	CommandPublisher
}

var busSubscriber Subscriber
var busPublisher Publisher
var store EventStore

func Register(bus Bus) {
	busSubscriber = bus
	busPublisher = bus
	busCommandSubscriber = bus
	busCommandPublisher = bus
}

func RegisterStore(s EventStore) {
	store = s
}

func Subscribe(topic string, cb EventHandler) {
	if busSubscriber != nil {
		busSubscriber.Subscribe(topic, cb)
	}
}

func Publish(event Event) {
	if busPublisher != nil {
		busPublisher.Publish(event)
		if store != nil {
			store.Save(event)
		}
	}
}

func PublishAsync(event Event) {
	if busPublisher != nil {
		busPublisher.PublishAsync(event)
		if store != nil {
			store.Save(event)
		}
	}
}
