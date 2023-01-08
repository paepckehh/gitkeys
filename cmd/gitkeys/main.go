// package main ...
package main

// import
import (
	"os"
	"syscall"

	"paepcke.de/gitkeys"
)

// const
const (
	_app      = "[gitkeys] "
	_err      = "[error] "
	_verbose  = "VERBOSE"
	_gitstore = "GITSTORE"
	_pinonly  = "PINONLY"
	_keyfile  = ".keys"
	_slashfwd = "/"
	_empty    = ""
	_erropt   = "commandline option unknown, example for online update: GITSTORE=/usr/store/git gitkeys fetch\n"
)

// main ...
func main() {
	r := gitkeys.NewRepo()
	if !isEnv(_gitstore) {
		errExit("env variable GITSTORE mandatory, example: REPO=/usr/store/git gitkeys")
	}
	r.GitStore, _ = getEnv(_gitstore)
	r.KeyFile = r.GitStore + _slashfwd + _keyfile
	switch {
	case isEnv(_pinonly):
		pin := r.Pinonly()
		if pin == _empty {
			os.Exit(1)
		}
		out(pin)
		os.Exit(0)
	case isOsArgs():
		switch os.Args[1] {
		case "check":
		case "fetch":
			if err := r.Fetch(); err != nil {
				errExit(err.Error())
			}
			os.Exit(0)
		default:
			errExit(_erropt)
		}
	}
	if err := r.Check(); err != nil {
		errExit(err.Error())
	}
	os.Exit(0)
}

//
// LITTLE GENERIC HELPER SECTION
//

// out ...
func out(msg string) {
	os.Stdout.Write([]byte(msg))
}

// errExit
func errExit(msg string) {
	out(_app + _err + msg)
	os.Exit(1)
}

// isEnv
func isEnv(in string) bool {
	if _, ok := syscall.Getenv(in); ok {
		return true
	}
	return false
}

// getEnv ...
func getEnv(in string) (string, bool) {
	return syscall.Getenv(in)
}

// isOsArgs ...
func isOsArgs() bool {
	return len(os.Args) > 1
}
