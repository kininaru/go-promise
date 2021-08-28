package promise

type Promise struct {
	State     int
	data      interface{}
	fn        func(func(interface{}), func(interface{}))
	resolveFn func(interface{})
	rejectFn  func(interface{})
}

func NewPromise(fn func(resolve, reject func(interface{}))) *Promise {
	return &Promise{
		State: PENDING,
		fn: fn,
	}
}

func (promise *Promise) Start() *Promise {
	if promise.State != PENDING {
		return promise
	}

	go promise.fn(promise.resolve, promise.reject)
	promise.State = RUNNING
	return promise
}

func NewPromiseAndStart(fn func(resolve, reject func(interface{}))) *Promise {
	return NewPromise(fn).Start()
}

func (promise *Promise) resolve(data interface{}) {
	if promise.State != RUNNING {
		return
	}

	promise.State = RESOLVED
	promise.data = data

	if promise.resolveFn != nil {
		promise.resolveFn(data)
		promise.State = DISCARD
	}
}

func (promise *Promise) reject(data interface{}) {
	if promise.State != RUNNING {
		return
	}

	promise.State = REJECTED
	promise.data = data

	if promise.rejectFn != nil {
		promise.rejectFn(data)
		promise.State = DISCARD
	}
}

func (promise *Promise) Then(resolveFn func(interface{}), fns ...func(interface{})) {
	var rejectFn func(interface{})
	switch len(fns) {
	case 1:
		rejectFn = fns[0]
	case 0:
		rejectFn = func(interface{}) {}
	default:
		panic("Too many functions")
	}

	if promise.State == DISCARD {
		return
	}

	promise.resolveFn = resolveFn
	promise.rejectFn = rejectFn

	if promise.State == RESOLVED {
		promise.State = DISCARD
		go resolveFn(promise.data)
	} else if promise.State == REJECTED {
		promise.State = DISCARD
		go rejectFn(promise.data)
	}
}
