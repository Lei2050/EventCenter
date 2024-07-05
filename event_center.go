package EventCenter

import (
	"reflect"
	"time"
)

type IEvent interface{}

type EventHandler[T IEvent] func(event T)

type IEventCenter[T IEvent] interface {
	On(EventHandler[T]) IEventCenter[T]
	Off(EventHandler[T]) IEventCenter[T]
	OnMonitor(handler EventHandler[T], checkCost time.Duration, callback func(event T, elapse time.Duration)) IEventCenter[T]
	OffMonitor(handler EventHandler[T]) IEventCenter[T]
	Fire(T)
}

type EventCenter[T IEvent] struct {
	handlers        []EventHandler[T]
	monitorHandlers []monitorHandlerWrap[T]
}

func (ec *EventCenter[T]) On(f EventHandler[T]) *EventCenter[T] {
	ec.handlers = append(ec.handlers, f)
	return ec
}

func (ec *EventCenter[T]) Off(f EventHandler[T]) *EventCenter[T] {
	var handlerSize = len(ec.handlers)
	if handlerSize <= 0 {
		return ec
	}
	pf := reflect.ValueOf(f).Pointer()
	idx := -1
	for k, v := range ec.handlers {
		if reflect.ValueOf(v).Pointer() == pf {
			idx = k
			break
		}
	}
	if idx != -1 {
		if idx != handlerSize-1 {
			ec.handlers[idx], ec.handlers[handlerSize-1] = ec.handlers[handlerSize-1], nil
		}
		ec.handlers = ec.handlers[:handlerSize-1]
	}
	return ec
}

type monitorHandlerWrap[T IEvent] struct {
	handler   EventHandler[T]
	checkCost time.Duration
	callback  func(event T, elapse time.Duration)
}

func (ec *EventCenter[T]) OnMonitor(handler EventHandler[T], checkCost time.Duration, callback func(event T, elapse time.Duration)) *EventCenter[T] {
	ec.monitorHandlers = append(ec.monitorHandlers, monitorHandlerWrap[T]{
		handler:   handler,
		checkCost: checkCost,
		callback:  callback,
	})
	return ec
}

func (ec *EventCenter[T]) OffMonitor(handler EventHandler[T]) *EventCenter[T] {
	var handlerSize = len(ec.monitorHandlers)
	if handlerSize <= 0 {
		return ec
	}
	pf := reflect.ValueOf(handler).Pointer()
	idx := -1
	for k, v := range ec.monitorHandlers {
		if reflect.ValueOf(v.handler).Pointer() == pf {
			idx = k
			break
		}
	}
	if idx != -1 {
		if idx != handlerSize-1 {
			ec.monitorHandlers[idx], ec.monitorHandlers[handlerSize-1] = ec.monitorHandlers[handlerSize-1], monitorHandlerWrap[T]{}
		}
		ec.monitorHandlers = ec.monitorHandlers[:handlerSize-1]
	}
	return ec
}

func (ec EventCenter[T]) Fire(event T) {
	for _, f := range ec.handlers {
		f(event)
	}
	for _, mh := range ec.monitorHandlers {
		st := time.Now().UnixMilli()
		mh.handler(event)
		elapse := time.Millisecond * time.Duration(time.Now().UnixMilli()-st)
		if elapse >= mh.checkCost {
			mh.callback(event, elapse)
		}
	}
}

type centerMgr struct {
	centers map[string]any
}

var gCenterMgr = centerMgr{centers: make(map[string]any)}

func GetOrCreateEventCenter[T IEvent](centerName string) *EventCenter[T] {
	c, exists := gCenterMgr.centers[centerName]
	if exists {
		return c.(*EventCenter[T])
	}
	center := &EventCenter[T]{}
	gCenterMgr.centers[centerName] = center
	return center
}

func On[T IEvent](centerName string, handler EventHandler[T]) *EventCenter[T] {
	return GetOrCreateEventCenter[T](centerName).On(handler)
}

func Off[T IEvent](centerName string, handler EventHandler[T]) *EventCenter[T] {
	return GetOrCreateEventCenter[T](centerName).Off(handler)
}

func OnMonitor[T IEvent](centerName string, handler EventHandler[T], checkCost time.Duration, callback func(event T, elapse time.Duration)) *EventCenter[T] {
	return GetOrCreateEventCenter[T](centerName).OnMonitor(handler, checkCost, callback)
}

func OffMonitor[T IEvent](centerName string, handler EventHandler[T]) *EventCenter[T] {
	return GetOrCreateEventCenter[T](centerName).OffMonitor(handler)
}

func Fire[T IEvent](centerName string, event T) {
	// reflect.TypeOf(event).Name()
	GetOrCreateEventCenter[T](centerName).Fire(event)
}
