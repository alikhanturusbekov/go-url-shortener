package example

import (
	"log"
	"os"
)

func badPanic() {
	panic("panic") // want "panic is forbidden"
}

func badLogFatal() {
	log.Fatal("error") // want "log.Fatal is forbidden outside main"
}

func badOSExit() {
	os.Exit(1) // want "os.Exit is forbidden outside main"
}
