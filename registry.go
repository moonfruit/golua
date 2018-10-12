package golua

type Registry interface {
	Ref(interface{}) uintptr
	UnRef(uintptr)
	Get(uintptr) interface{}
}

type sliceRegistry struct {
	store []interface{}
	frees []uintptr
}

func (r *sliceRegistry) Ref(val interface{}) uintptr {
	if val == nil {
		panic("cannot register nil")
	}

	freeLen := len(r.frees)
	if freeLen > 0 {
		last := freeLen - 1
		id := r.frees[last]
		r.frees = r.frees[:last]
		r.store[id] = val
		return id + 1
	}

	id := uintptr(len(r.store))
	r.store = append(r.store, val)
	return id + 1
}

func (r *sliceRegistry) UnRef(id uintptr) {
	id--
	if r.get(id) == nil {
		return
	}

	r.store[id] = nil
	r.frees = append(r.frees, id)
}

func (r *sliceRegistry) Get(id uintptr) interface{} {
	id--
	return r.get(id)
}

func (r *sliceRegistry) get(id uintptr) interface{} {
	if id >= uintptr(len(r.store)) {
		return nil
	}
	return r.store[id]
}

type mapRegistry struct {
	store map[uintptr]interface{}
	next  uintptr
}

func (r *mapRegistry) Ref(val interface{}) uintptr {
	id := r.next

	if r.store == nil {
		r.store = make(map[uintptr]interface{})
	}

	r.store[id] = val
	r.next++

	return id + 1
}

func (r *mapRegistry) UnRef(id uintptr) {
	id--
	delete(r.store, id)
}

func (r *mapRegistry) Get(id uintptr) interface{} {
	id--
	return r.store[id]
}
