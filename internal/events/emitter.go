package events

import "sync"

type (
	// EventName is a string alias for event names.
	EventName string

	// Listener is a function that is called when an event is fired.
	Listener func(EventName, interface{}) error

	// Emitter is an interface that allows to fire events.
	Emitter interface {
		// Emit fires an event with the given name and payload.
		Emit(EventName, interface{})
		// On registers a listener for the given event name.
		On(EventName, ...Listener)
		// OnMany registers a listener for the given event names.
		ListenEvents(Listener, ...EventName)
	}

	// Logger is an interface that allows to log events.
	Logger interface {
		Debugf(format string, args ...interface{})
		Infof(format string, args ...interface{})
		Errorf(format string, args ...interface{})
	}

	emitter struct {
		sync.RWMutex
		listeners map[EventName][]Listener
		log       Logger
	}
)

// NewEmitter creates a new Emitter.
func NewEmitter(log Logger) Emitter {
	return &emitter{
		listeners: make(map[EventName][]Listener),
		log:       log,
	}
}

// Emit fires an event with the given name and payload.
func (e *emitter) Emit(name EventName, payload interface{}) {
	e.RLock()
	defer e.RUnlock()

	for _, listener := range e.listeners[name] {
		if listener != nil {
			go func(fn Listener, i interface{}) {
				if err := fn(name, payload); err != nil {
					e.log.Errorf("failed to handle event %s: %s", name, err.Error())
				}
			}(listener, payload)
		}
	}

	return
}

// On registers a listener for the given event name.
func (e *emitter) On(name EventName, listeners ...Listener) {
	e.Lock()
	defer e.Unlock()

	e.listeners[name] = append(e.listeners[name], listeners...)
}

// ListenEvents registers a listener for the given event names.
func (e *emitter) ListenEvents(listener Listener, names ...EventName) {
	e.Lock()
	defer e.Unlock()

	for _, name := range names {
		e.listeners[name] = append(e.listeners[name], listener)
	}
}
