package broadcast

import (
	"sync"

	"github.com/morrocker/log"
	"github.com/morrocker/utils"
)

type Broadcaster struct {
	lock      sync.Mutex
	listeners map[string]*Listener
}

type Listener struct {
	id string
	b  *Broadcaster
	C  chan interface{}
}

func New() *Broadcaster {
	return &Broadcaster{
		listeners: make(map[string]*Listener),
	}
}

func (b *Broadcaster) ListenTo(id string) *Listener {
	b.lock.Lock()
	defer b.lock.Unlock()
	_, ok := b.listeners[id]
	if !ok {
		l := &Listener{
			id: id,
			C:  make(chan interface{}),
			b:  b,
		}
		b.listeners[id] = l
		return l
	}
	log.Error("Listeners ID already taken")
	return nil
}

func (b *Broadcaster) Listen() *Listener {
	b.lock.Lock()
	defer b.lock.Unlock()
	newId := b.newID()
	l := &Listener{
		id: newId,
		C:  make(chan interface{}),
		b:  b,
	}
	b.listeners[newId] = l

	return l
}

func (b *Broadcaster) Send(id string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	for _, l := range b.listeners {
		if l.id == id {
			l.C <- ""
		}
	}
}

func (b *Broadcaster) Close(id string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	l, ok := b.listeners[id]
	if ok {
		close(l.C)
		delete(b.listeners, id)
	}
}

func (b *Broadcaster) Broadcast() {
	b.lock.Lock()
	defer b.lock.Unlock()
	x := utils.RandString(8)
	log.Bench("Starting Broadcast %s to %d listeners", x, len(b.listeners))
	for _, l := range b.listeners {
		log.Bench("Starting Broadcast %s to %d listeners", x, len(b.listeners))
		l.C <- ""
	}
	log.Bench("Ending Broadcast %s", x)
}

func (b *Broadcaster) CloseAll() {
	b.lock.Lock()
	defer b.lock.Unlock()
	for _, l := range b.listeners {
		close(l.C)
	}
}

func (l *Listener) Close() {
	l.b.Close(l.id)
}
func (l *Listener) ID() string {
	return l.id
}

func (b *Broadcaster) newID() string {
	for {
		newId := utils.RandString(8)
		_, ok := b.listeners[newId]
		if !ok {
			return newId
		}
	}
}
