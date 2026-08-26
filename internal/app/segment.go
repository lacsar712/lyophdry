package app

import (
	"context"

	"github.com/lacsar712/lyophdry/internal/clock"
)

// SegmentPlan describes inter-segment vent staging for drying batches.

type SegmentPlan struct {
	VentSteps int
}

func (a *App) ExecutePlan(ctx context.Context, plan SegmentPlan) error {
	if a.scheduler == nil {
		return nil
	}
	// Propagate the caller's context (batch/HMI scope) into the coordination
	// layer so a segment-revocation cancel reaches the scheduler. Substituting
	// context.Background() here would sever the cancellation signal, leaving
	// revoked ramp steps installed in the timeline.
	return a.scheduler.InstallVentPlanCtx(ctx, clock.VentPlan{VentSteps: plan.VentSteps}, "segment-plan")
}

func (a *App) SegmentVentStepsDone() int {
	if a.scheduler == nil {
		return 0
	}
	return a.scheduler.VentStepsDone()
}
