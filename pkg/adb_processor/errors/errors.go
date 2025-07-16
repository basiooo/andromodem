// Package errors defines custom error types used throughout the ADB processing modules.
package errors

import _errors "errors"

var (
	ErrorNeedRoot                        = _errors.New("cannot perform action without root")
	ErrorNeedShellSuperUserPermission    = _errors.New("cannot perform action without allow super-user permission for 'com.android.shell'")
	ErrorMinimumAndroidVersionNotSupport = _errors.New("cannot perform action minimum android version not supported")
	ErrorDeviceIsNil                     = _errors.New("device is nil")
)
