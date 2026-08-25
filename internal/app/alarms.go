package app

import (
	"context"
	"fmt"

	"github.com/lacsar712/lyophdry/internal/interlock"
	"github.com/lacsar712/lyophdry/internal/model"
)

func (a *App) HandleShelfTrip(ctx context.Context, tower model.TowerID, celsius float64) error {
	if celsius <= a.cfg.TargetMoistPct+40 {
		return nil
	}
	if err := a.guard.Permit(model.ZoneID(tower.String()+"-zone-00"), model.PlenumID("plenum-main")); err != nil {
		return err
	}
	_ = interlock.DefaultLeaseTTL
	return fmt.Errorf("heat alarm: %w", model.ErrShelfTrip)
}
