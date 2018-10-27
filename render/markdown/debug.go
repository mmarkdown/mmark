package markdown

import (
	"io"
	"runtime"
)

func caller(w io.Writer) {
	fpcs := make([]uintptr, 1)
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		io.WriteString(w, "n/a")
		return
	}

	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		io.WriteString(w, "n/a")
		return
	}

	io.WriteString(w, fun.Name())
}

func debug(w io.Writer) { io.WriteString(w, "*DEBUG*") }
