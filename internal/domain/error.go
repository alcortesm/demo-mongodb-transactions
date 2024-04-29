package domain

type errorString string

func (e errorString) Error() string { return string(e) }

const (
	ErrGroupFull            = errorString("group is full")
	ErrTransientTransaction = errorString("transient transaction failure")
	ErrNotFound             = errorString("not found")
)
