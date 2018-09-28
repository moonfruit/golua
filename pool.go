package golua

import (
	"sync"
	"unsafe"
)

var pool = objectPool{}

type objectPool struct {
	store map[uintptr]interface{}
	mutex sync.RWMutex
}

func (p *objectPool) Ref(object interface{}) uintptr {
	id := uintptr(unsafe.Pointer(&object))

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.store == nil {
		p.store = make(map[uintptr]interface{})
	}

	p.store[id] = object
	return id
}

func (p *objectPool) UnRef(id uintptr) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	object, ok := p.store[id]
	if !ok {
		return
	}

	delete(p.store, id)

	type Closer interface {
		Close()
	}
	if closer, ok := object.(Closer); ok {
		closer.Close()
	}
}

func (p *objectPool) Get(id uintptr) (interface{}, bool) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	object, ok := p.store[id]
	return object, ok
}
