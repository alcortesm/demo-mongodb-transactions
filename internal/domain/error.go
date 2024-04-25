package domain

import "errors"

var (
	ErrGroupFull            = errors.New("group is full")
	ErrTransientTransaction = errors.New("transient transaction failure")
	ErrNotFound             = errors.New("not found")
)
