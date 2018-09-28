package golua

import (
	"reflect"
	"sync"
)

var pool = objPool{}

type objRef struct {
	obj interface{}
	cnt int
}

type objPool struct {
	store map[uintptr]*objRef
	mutex sync.RWMutex
}

func (p *objPool) Ref(obj interface{}) uintptr {
	id := reflect.ValueOf(obj).Pointer()

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.store == nil {
		p.store = make(map[uintptr]*objRef)
	}

	ref, ok := p.store[id]
	if ok {
		ref.cnt++
	} else {
		p.store[id] = &objRef{obj, 1}
	}

	return id
}

func (p *objPool) UnRef(id uintptr) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	ref, ok := p.store[id]
	if !ok {
		return
	}

	ref.cnt--
	if ref.cnt > 0 {
		return
	}

	delete(p.store, id)

	type Closer interface {
		Close()
	}
	if closer, ok := ref.obj.(Closer); ok {
		closer.Close()
	}
}

func (p *objPool) Get(id uintptr) (interface{}, bool) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	ref, ok := p.store[id]
	if !ok {
		return nil, false
	}

	return ref.obj, true
}
