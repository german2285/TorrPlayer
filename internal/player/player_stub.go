// +build !windows

package player

import "fmt"

func PlayVideoWithMPV(streamURL string) error {
	return fmt.Errorf("MPV playback is only supported on Windows")
}
