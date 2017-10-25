package util

import "log"

func Safe(f func()) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Panicln(rec)
		}
	}()

	f()
}
