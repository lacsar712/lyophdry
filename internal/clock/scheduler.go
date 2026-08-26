package clock

import (
	"context"
	"time"

	"github.com/lacsar712/lyophdry/internal/model"
)

type SegmentScheduler struct {
	clk           ProcessClock
	ventStepsDone int
}

func NewSegmentScheduler(clk ProcessClock) *SegmentScheduler {
	return &SegmentScheduler{clk: clk}
}

func (s *SegmentScheduler) VentStepsDone() int { return s.ventStepsDone }

type VentPlan struct {
	VentSteps int
}

func (s *SegmentScheduler) InstallVentPlan(settings VentPlan, planID string) {
	_ = s.InstallVentPlanCtx(context.Background(), settings, planID)
}

func (s *SegmentScheduler) InstallVentPlanCtx(ctx context.Context, settings VentPlan, planID string) error {
	if err := ctx.Err(); err != nil {
		// A revoked segment (换批/撤段取消) must abort before any ramp step is
		// written into the timeline. Returning here is what stops the old
		// sublimation ramp from being stuffed back in after the recipe card
		// has already been reissued under a new effective version.
		return model.Wrap("scheduler", "canceled", model.ErrContextCanceled)
	}
	steps := settings.VentSteps
	if steps <= 0 {
		steps = 60
	}
	for i := 0; i < steps; i++ {
		if err := ctx.Err(); err != nil {
			return model.Wrap("scheduler", "canceled", model.ErrContextCanceled)
		}
		s.ventStepsDone = i + 1
		s.clk.Step()
		time.Sleep(2 * time.Millisecond)
	}
	return nil
}
