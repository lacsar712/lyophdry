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
		return fmt.Errorf("gradient hold: window not satisfied")
	}
	return nil
}
