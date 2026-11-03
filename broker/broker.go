package broker

import "sync"

type Broker[T any] struct {
	mu          sync.RWMutex
	messages    []Message[T]
	subscribers map[*Subscriber[T]]struct{}
}

type Message[T any] struct {
	Payload T
	Offset  int64
}

func NewBroker[T any]() *Broker[T] {
	return &Broker[T]{subscribers: make(map[*Subscriber[T]]struct{})}
}

func (b *Broker[T]) append(payload T) Message[T] {
	b.mu.Lock()
	defer b.mu.Unlock()

	offset := int64(len(b.messages))
	message := Message[T]{Payload: payload, Offset: offset}
	b.messages = append(b.messages, message)
	return message
}

func (b *Broker[T]) messageAt(offset int64) (Message[T], bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if offset >= int64(len(b.messages)) {
		return Message[T]{}, false
	}
	return b.messages[offset], true
}

func (b *Broker[T]) registerSubscriber(s *Subscriber[T]) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[s] = struct{}{}
}

func (b *Broker[T]) SubscriberCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.subscribers)
}

type Publisher[T any] struct {
	broker *Broker[T]
}

func NewPublisher[T any](b *Broker[T]) *Publisher[T] {
	return &Publisher[T]{broker: b}
}

func (p *Publisher[T]) Publish(payload T) Message[T] {
	return p.broker.append(payload)
}

type Subscriber[T any] struct {
	name       string
	mu         sync.RWMutex
	broker     *Broker[T]
	nextOffset int64
}

func NewSubscriber[T any](name string) *Subscriber[T] {
	return &Subscriber[T]{name: name}
}

func (s *Subscriber[T]) Subscribe(b *Broker[T]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.broker != nil || b == nil {
		return
	}

	s.broker = b
	b.registerSubscriber(s)
}

func (s *Subscriber[T]) Next() (Message[T], bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.broker == nil {
		return Message[T]{}, false
	}

	msg, ok := s.broker.messageAt(s.nextOffset)

	if !ok {
		return Message[T]{}, false
	}

	s.nextOffset++

	return msg, true
}

func (s *Subscriber[T]) Name() string { return s.name }
