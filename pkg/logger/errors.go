package logger

import "errors"

var ErrInvalidLogFormat = errors.New("provided invalid log format (must be json or text)")
