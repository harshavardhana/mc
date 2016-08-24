// +build windows

package cmd

import "github.com/rjeczalik/notify"

var (
	EventTypePut    []notify.Event = []notify.Event{notify.Create, notify.Write, notify.Rename}
	EventTypeDelete []notify.Event = []notify.Event{notify.Remove}
)

func IsPutEvent(event notify.Event) bool {
	switch event {
	case notify.Create:
		return true
	case notify.Rename:
		return true
	case notify.Write:
		return true
	}

	return false
}

func IsDeleteEvent(event notify.Event) bool {
	switch event {
	case notify.Remove:
		return true
	}
	// rename here as well?

	return false
}
