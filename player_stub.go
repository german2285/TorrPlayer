// +build !windows

package main

import "fmt"

func playVideoWithMPV(streamURL string) error {
	return fmt.Errorf("MPV playback is only supported on Windows")
}
