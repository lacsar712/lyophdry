package app

import (
	"context"
	"fmt"
	"time"

	"github.com/lacsar712/lyophdry/internal/model"
)

func (a *App) ValidateMoldDrift(ctx context.Context, moistPct float64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	limit := a.cfg.TargetMoistPct + a.cfg.MaxGradientDeltaPct
	if moistPct <= limit {
		return nil
	}
	return model.DriftChain(moistPct)
}

func (a *App) ConfirmBladderHold(ctx context.Context, anchor time.Time) error {
	if a.avgWindow == nil {
		return model.Wrap("app", "window", model.ErrCondHold)
	}
	if err := a.avgWindow.Require(anchor); err != nil {
		// Preserve the ErrCondHold sentinel so downstream alarm
		// classification (model.IsHold) still recognizes a limit-hold
		// after the short handling window has closed. Without %w the
		// sentinel is lost and the fault degrades to "unknown", greying
		// out the recovery entry on the night shift.
		return fmt.Errorf("gradient hold: window not satisfied: %w", err)
	}
	return nil
}
