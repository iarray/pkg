package inmemorybus

import (
	"errors"
	"log"
	"sync"

	"github.com/iarray/pkg/ddd/infrastruct/eventbus"
)

var EMPTY_COMMAND eventbus.CommandResult

func init() {
	bus := NewBus()
	EMPTY_COMMAND = eventbus.NewCommandResult(nil, errors.New("Unprocessed"))
	eventbus.Register(bus)
}

type SubscriberList []*subscriberInfo

type subscriberInfo struct {
	eventbus.Subscription
	handler eventbus.EventHandler
}

type EventBus struct {
	sync.Mutex
	cmdLock         sync.Mutex
	subscribers     map[string]SubscriberList
	commandHandlers map[string]eventbus.CommandHandler
	count           uint64
}

func NewBus() eventbus.Bus {
	return &EventBus{
		subscribers:     make(map[string]SubscriberList),
		commandHandlers: make(map[string]eventbus.CommandHandler),
	}
}

func (b *EventBus) Subscribe(topic string, handler eventbus.EventHandler) eventbus.Subscription {
	b.Lock()
	defer b.Unlock()
	log.Printf("subscribe [%s] ", topic)
	b.count++
	subc := &subscriberInfo{
		Subscription: eventbus.Subscription{
			Topic: topic,
			ID:    b.count,
		},
		handler: handler,
	}

	b.subscribers[topic] = append(b.subscribers[topic], subc)
	return subc.Subscription
}

func (b *EventBus) Unsubscribe(info *eventbus.Subscription) {
	b.Lock()
	defer b.Unlock()
	if list, ok := b.subscribers[info.Topic]; ok {
		for idx, item := range list {
			if item.ID == info.ID {
				b.subscribers[info.Topic] = append(list[:idx], list[idx+1:]...)
				break
			}
		}
	}
}

func (b *EventBus) copySubscriptions(topic string) SubscriberList {
	b.Lock()
	defer b.Unlock()
	if infos, ok := b.subscribers[topic]; ok {
		return infos
	}
	return nil
}

func (b *EventBus) Publish(event eventbus.Event) {
	subs := b.copySubscriptions(event.EventName())
	for _, sub := range subs {
		sub.handler(event)
	}
}

func (b *EventBus) PublishAsync(event eventbus.Event) {
	subs := b.copySubscriptions(event.EventName())
	for _, sub := range subs {
		go sub.handler(event)
	}
}

func (b *EventBus) SendCommandAsync(cmd eventbus.Command) (eventbus.Feature, error) {
	b.cmdLock.Lock()
	defer b.cmdLock.Unlock()
	if handler, ok := b.commandHandlers[cmd.Action()]; ok {
		ret := make(chan eventbus.CommandResult)
		go func() {
			res := handler(cmd)
			ret <- res
		}()
		return ret, nil
	}
	return nil, errors.New("command handler not registed !")

}

func (b *EventBus) SendCommand(cmd eventbus.Command) eventbus.CommandResult {
	b.cmdLock.Lock()
	defer b.cmdLock.Unlock()
	if handler, ok := b.commandHandlers[cmd.Action()]; ok {
		return handler(cmd)
	}
	return EMPTY_COMMAND

}

func (b *EventBus) RegisterCommandHandler(action string, handler eventbus.CommandHandler) {
	b.cmdLock.Lock()
	defer b.cmdLock.Unlock()
	log.Printf("registed [%s] handler", action)
	b.commandHandlers[action] = handler
}
