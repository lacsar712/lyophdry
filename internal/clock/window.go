package clock

import (
	"time"

	"github.com/lacsar712/lyophdry/internal/model"
)

type SublimWindow struct {
	clk      Clock
	duration time.Duration
}

func NewSublimWindow(clk Clock, duration time.Duration) *SublimWindow {
	if duration <= 0 {
		duration = 2 * time.Minute
	}
	return &SublimWindow{clk: clk, duration: duration}
}

func (w *SublimWindow) Active(anchor time.Time) bool {
	return ProcessWindowOpen(w.clk, anchor, w.duration)
}

func (w *SublimWindow) Require(anchor time.Time) error {
	// The sublimation requirement is met only once the window has fully
	// elapsed. While the window is still open (升华进行中, remaining minutes
	// still ticking) the requirement is NOT satisfied, so downstream
	// coordination must not advance — otherwise the coordination page pops
	// an empty gray step area before sublimation is done.
	if ProcessWindowOpen(w.clk, anchor, w.duration) {
		return model.ErrCondHold
	}
	if ProcessWindowClosed(w.clk, anchor, w.duration) {
		return nil
	}
	return model.ErrCondHold
}
