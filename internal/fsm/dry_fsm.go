package fsm

import (
	"context"
	"fmt"

	"github.com/lacsar712/lyophdry/internal/model"
)

var ErrIllegalDryTransition = fmt.Errorf("illegal dry transition")

type LyoFSM struct {
	id    model.TowerID
	state model.DryState
	hooks *DryHookChain
}

func NewLyoFSM(id model.TowerID, effect func(context.Context, model.TowerID, model.DryState, model.DryState) error) *LyoFSM {
	_ = effect
	return &LyoFSM{id: id, state: model.DryIdle, hooks: NewDryHookChain()}
}

func (f *LyoFSM) Hooks() *DryHookChain { return f.hooks }

func (f *LyoFSM) State() model.DryState { return f.state }

func (f *LyoFSM) Dispatch(ctx context.Context, event string) (model.DryState, error) {
	next, ok := allowedDry(f.state, event)
	if !ok {
		// Illegal transition: no state change and no after-hooks. Running
		// RunAfter here would fire drive side effects (e.g. the shelf-heat
		// pulse registered via RegisterDryHeatHook) even though the FSM never
		// left its current state — observed in production as an isolated
		// heater pulse while the tower marker stayed in standby.
		return f.state, fmt.Errorf("%s from %s: %w", event, f.state, ErrIllegalDryTransition)
	}
	from := f.state
	if f.hooks != nil {
		if err := f.hooks.RunBefore(ctx, from, next, event); err != nil {
			return f.state, err
		}
	}
	f.state = next
	if f.hooks != nil {
		if err := f.hooks.RunAfter(ctx, from, next, event); err != nil {
			return f.state, err
		}
	}
	return f.state, nil
}

func allowedDry(from model.DryState, event string) (model.DryState, bool) {
	switch from {
	case model.DryIdle:
		if event == "arm_heat" {
			return model.DryHeating, true
		}
	case model.DryHeating:
		if event == "hold" {
			return model.DryHold, true
		}
	case model.DryHold:
		if event == "cool" {
			return model.DryCool, true
		}
	case model.DryCool:
		if event == "done" {
			return model.DryIdle, true
		}
	}
	return from, false
}
