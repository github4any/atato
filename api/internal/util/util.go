package util

import (
	"runtime/debug"

	log "github.com/sirupsen/logrus"
)

// PanicOnError calls panic() in case the error is not nil
func PanicOnError(e error) {
	if e != nil {
		panic(e)
	}
}

// RecoverAtStartup recovers in case of panic at start-up
func RecoverAtStartup() {
	if r := recover(); r != nil {
		log.Error("Startup error: ", r)
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Error(string(debug.Stack()))
		}
		log.Fatal("App stop...")
	}
}
