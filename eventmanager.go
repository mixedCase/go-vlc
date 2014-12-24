// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0

package vlc

//#include <stdlib.h>
//#include <vlc/vlc.h>
// extern void  goEventCB(const struct libvlc_event_t*, void*);
// static int goAttach(libvlc_event_manager_t* em, libvlc_event_type_t et, void* userdata) {
// 	return libvlc_event_attach(em, et, goEventCB, userdata);
// }
//
// static void goDetach(libvlc_event_manager_t* em, libvlc_event_type_t et, void* userdata) {
// 	libvlc_event_detach(em, et, goEventCB, userdata);
// }
import "C"
import (
	"sync"
	"syscall"
	"unsafe"
)

var (
	// because the go callback handler is not aware of EventManagers, and must still be ----
	events     = make(map[int]*eventData)
	eventsInc  IncrementalInt
	eventsLock sync.RWMutex
)

// A libvlc instance has an event manager which can be used to hook event callbacks,
type EventManager struct {
	ptr *C.libvlc_event_manager_t
	m   *sync.Mutex
}

func NewEventManager(p *C.libvlc_event_manager_t) *EventManager {
	ev := &EventManager{
		ptr: p,
	}
	return ev
}

// Attach registers the given event handler and returns a unique id
// we can use to detach the event at a later point.
func (this *EventManager) Attach(et EventType, cb EventHandler, userdata interface{}) (int, error) {
	if this.ptr == nil {
		return 0, syscall.EINVAL
	}
	if cb == nil {
		panic("cannot attach nil EventHandler")
	}

	ed := &eventData{
		id: eventsInc.Next(),
		t:  C.libvlc_event_type_t(et),
		f:  cb,
		d:  userdata,
	}

	eventsLock.Lock()
	events[ed.id] = ed
	eventsLock.Unlock()

	if C.goAttach(this.ptr, ed.t, unsafe.Pointer(&ed.id)) != 0 {
		err := checkError()
		if err != nil {
			return 0, err
		}
	}

	return ed.id, nil
}

// Detach unregisters the given event id.
func (this *EventManager) Detach(id int) (err error) {
	if this.ptr == nil {
		return syscall.EINVAL
	}

	var ed *eventData
	var ok bool

	eventsLock.Lock()
	if ed, ok = events[id]; !ok {
		eventsLock.Unlock()
		return syscall.EINVAL
	}

	delete(events, id)
	eventsLock.Unlock()

	C.goDetach(this.ptr, ed.t, unsafe.Pointer(&ed.id))
	return
}
