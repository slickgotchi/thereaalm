package utils

import "runtime"

// GetFuncName returns the name of the calling function
func GetFuncName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "UnknownFunction"
	}
	return runtime.FuncForPC(pc).Name()
}