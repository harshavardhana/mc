// +build linux

package cmd

var (
	EventTypePut    []notify.Event = []notify.Event{notify.InCloseWrite | notify.InMovedTo}
	EventTypeDelete []notify.Event = []notify.Event{notify.InDelete | notify.InDeleteSelf | notify.InMovedFrom}
)

func IsPutEvent(event notify.Event) bool {
	switch event {
	case notify.InCloseWrite:
		return true
	case notify.InMovedTo:
		return true
	}

	return false
}

func IsDeleteEvent(event notify.Event) bool {
	switch event {
	case notify.InDelete:
		return true
	case notify.InDeleteSelf:
		return true
	case notify.InMovedFrom:
		return true
	}

	return false
}
