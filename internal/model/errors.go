package model

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidID       = errors.New("lyophdry: invalid identifier")
	ErrNotFound        = errors.New("lyophdry: entity not found")
	ErrConflict        = errors.New("lyophdry: state conflict")
	ErrInterlock       = errors.New("lyophdry: interlock denied")
	ErrMoistureHold    = errors.New("lyophdry: moisture hold active")
	ErrAirflowSetpoint = errors.New("lyophdry: airflow setpoint violation")
	ErrFanFault        = errors.New("lyophdry: fan fault")
	ErrScheduleEmpty   = errors.New("lyophdry: schedule empty")
	ErrGradient        = errors.New("lyophdry: moisture gradient violation")
	ErrShelfDrift   = errors.New("lyophdry: moisture drift exceeded")
	ErrShelfTrip    = errors.New("lyophdry: heat overtemperature")
	ErrCondHold    = errors.New("lyophdry: gradient hold not satisfied")
	ErrContextCanceled = errors.New("lyophdry: operation canceled")
)

type DomainError struct {
	Op   string
	Code string
	Err  error
}

func (e *DomainError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Err != nil {
		return fmt.Sprintf("lyophdry %s [%s]: %v", e.Op, e.Code, e.Err)
	}
	return fmt.Sprintf("lyophdry %s [%s]", e.Op, e.Code)
}

func (e *DomainError) Unwrap() error { return e.Err }

func Wrap(op, code string, err error) error {
	if err == nil {
		return nil
	}
	return &DomainError{Op: op, Code: code, Err: err}
}

func Is(err, target error) bool   { return errors.Is(err, target) }
func As(err error, target any) bool { return errors.As(err, target) }
