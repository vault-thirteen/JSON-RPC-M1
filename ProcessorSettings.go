package jrm1

import "errors"

const (
	ErrEnableExceptionCaptureToLogThem = "enable exception capture to log them"
	ErrMetaDataFieldNameConflict       = "meta data field name conflict"
)

// ProcessorSettings are settings of the RPC processor (server).
type ProcessorSettings struct {
	// When enabled, RPC processor (server) will catch exceptions.
	CatchExceptions bool

	// When enabled, RPC processor (server) will journal exceptions.
	LogExceptions bool

	// When enabled, RPC processor (server) will count requests.
	CountRequests bool

	// Name of a meta-data field where to store request duration.
	// When enabled, RPC processor (server) will measure time taken to run
	// functions. Duration is shown only for successful function calls.
	// To enable this feature, set the field name as non-null value.
	DurationFieldName *string

	// Name of a meta-data field where to store ID of the current request.
	// When enabled, RPC processor (server) will add ID of the current request
	// as a meta-data field, so that user function will be able to read it.
	// This field is automatically removed when function call finishes.
	// To enable this feature, set the field name as non-null value.
	RequestIdFieldName *string
}

// Check verifies processor's settings.
func (ps *ProcessorSettings) Check() (err error) {
	if ps.LogExceptions && (!ps.CatchExceptions) {
		return errors.New(ErrEnableExceptionCaptureToLogThem)
	}

	if ps.DurationFieldName != nil && ps.RequestIdFieldName != nil {
		if *ps.DurationFieldName == *ps.RequestIdFieldName {
			return errors.New(ErrMetaDataFieldNameConflict)
		}
	}

	return nil
}

// isDurationEnabled tells whether time measurement is enabled.
func (ps *ProcessorSettings) isDurationEnabled() bool {
	return ps.DurationFieldName != nil
}

// isRequestIdShown tells whether request ID is added to the meta-data set.
// Note that request ID is added to the meta-data set only for the duration of
// the function call. When the requested function returns, the ID is removed
// from the meta-data set.
func (ps *ProcessorSettings) isRequestIdShown() bool {
	return ps.RequestIdFieldName != nil
}
