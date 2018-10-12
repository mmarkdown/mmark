package markdown

import "runtime"

func caller() string {
	fpcs := make([]uintptr, 1)

	// skip 3 levels to get to the caller of whoever called Caller()
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return "n/a"
	}

	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "n/a"
	}

	return fun.Name()
}
