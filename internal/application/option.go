package application

import "time"

type Option interface {
	option()
}

// DelayBeforeUpdating introduces an artificial delay before calling
// Store.Update, which can be used to improve the chance concurrent application
// calls will race against each other during testing.
type DelayBeforeUpdating time.Duration

func (DelayBeforeUpdating) option() {}
