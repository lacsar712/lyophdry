package model

import (
	"errors"
	"fmt"
)

func DriftChain(pct float64) error {
	return fmt.Errorf("moisture %.1f drift high: %w", pct, ErrShelfDrift)
}

func TripChain(zone string, celsius float64) error {
	return fmt.Errorf("heat alarm zone %s at %.1f: %w", zone, celsius, ErrShelfTrip)
}

func HoldChain(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("gradient hold: %w", err)
}

func IsDrift(err error) bool { return errors.Is(err, ErrShelfDrift) }
func IsTrip(err error) bool { return errors.Is(err, ErrShelfTrip) }
func IsHold(err error) bool { return errors.Is(err, ErrCondHold) }
