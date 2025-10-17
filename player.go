// +build windows

package main

/*
#cgo windows LDFLAGS: -L. -lmpv-2
#include <mpv/client.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func playVideoWithMPV(streamURL string) error {
	// Create MPV instance
	mpv := C.mpv_create()
	if mpv == nil {
		return fmt.Errorf("failed to create MPV instance")
	}
	defer C.mpv_terminate_destroy(mpv)

	// Helper function to set options
	setOption := func(name, value string) error {
		cName := C.CString(name)
		cValue := C.CString(value)
		defer C.free(unsafe.Pointer(cName))
		defer C.free(unsafe.Pointer(cValue))

		ret := C.mpv_set_option_string(mpv, cName, cValue)
		if ret != 0 {
			return fmt.Errorf("failed to set option %s=%s (error code: %d)", name, value, int(ret))
		}
		return nil
	}

	// Configure MPV player
	if err := setOption("vo", "gpu"); err != nil {
		return err
	}
	if err := setOption("keepaspect", "yes"); err != nil {
		return err
	}
	if err := setOption("keepaspect-window", "no"); err != nil {
		return err
	}
	if err := setOption("osc", "yes"); err != nil {
		return err
	}
	if err := setOption("input-default-bindings", "yes"); err != nil {
		return err
	}
	if err := setOption("input-vo-keyboard", "yes"); err != nil {
		return err
	}

	// Cache settings for streaming
	if err := setOption("cache", "yes"); err != nil {
		return err
	}
	if err := setOption("demuxer-max-bytes", "512M"); err != nil {
		return err
	}
	if err := setOption("demuxer-max-back-bytes", "256M"); err != nil {
		return err
	}

	// Initialize MPV
	ret := C.mpv_initialize(mpv)
	if ret != 0 {
		return fmt.Errorf("failed to initialize MPV (error code: %d)", int(ret))
	}

	// Load stream URL
	cStreamURL := C.CString(streamURL)
	cLoadfile := C.CString("loadfile")
	defer C.free(unsafe.Pointer(cStreamURL))
	defer C.free(unsafe.Pointer(cLoadfile))

	cmd := []*C.char{
		cLoadfile,
		cStreamURL,
		nil, // NULL terminator
	}

	ret = C.mpv_command(mpv, &cmd[0])
	if ret != 0 {
		return fmt.Errorf("failed to load stream (error code: %d)", int(ret))
	}

	// Event loop - wait for playback to finish
	for {
		event := C.mpv_wait_event(mpv, -1) // Wait indefinitely
		if event == nil {
			break
		}

		eventID := event.event_id

		if eventID == C.MPV_EVENT_SHUTDOWN {
			break
		}
		if eventID == C.MPV_EVENT_END_FILE {
			break
		}
	}

	return nil
}
