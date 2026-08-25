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
	return WindowElapsed(w.clk, anchor, w.duration)
}

func (w *SublimWindow) Require(anchor time.Time) error {
	if w.Active(anchor) {
		return nil
	}
	return model.ErrCondHold
}
