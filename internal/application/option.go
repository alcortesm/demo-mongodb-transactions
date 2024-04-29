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

func mustDelayBeforeUpdating(options ...Option) (time.Duration, bool) {
	for _, o := range options {
		if raw, ok := o.(DelayBeforeUpdating); ok {
			return time.Duration(raw), true
		}
	}

	return time.Duration(0), false
}

type EnableTransactions struct{}

func (EnableTransactions) option() {}

func areTransactionsEnabled(options ...Option) bool {
	for _, o := range options {
		if _, ok := o.(EnableTransactions); ok {
			return true
		}
	}

	return false
}
