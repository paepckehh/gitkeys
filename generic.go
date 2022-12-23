// package gitkeys ...
package gitkeys

// import
import (
	"encoding/base64"
	"errors"
	"os"
	"runtime"
	"strconv"
	"time"
)

//
// DISPLAY I/O
//

// out ...
func out(msg string) { outPlain(msg + _linefeed) }

// outPlain ...
func outPlain(msg string) { os.Stdout.Write([]byte(msg)) }

// outINFO ...
func outINFO(msg string) { outPlain(_app + msg + _linefeed) }

// outERR ...
func outERR(msg string) { outPlain(_app + _err + msg + _linefeed) }

//
// FILE I/O
//

// isFile ...
func isFile(filename string) bool {
	fi, err := os.Lstat(filename)
	if err != nil {
		return false
	}
	return fi.Mode().IsRegular()
}

// saveWriteFile ...
func saveWriteFile(filename string, data *[]byte) error {
	if isFile(filename) {
		if err := os.Rename(filename, filename+".old"); err != nil {
			return err
		}
	}
	if err := os.WriteFile(filename, *data, 0o644); err != nil {
		return err
	}
	return nil
}

//
// CONVERTER
//

// itoa ...
func itoa(in int) string { return strconv.Itoa(in) }

// itoa64 ...
func itoa64(in int64) string { return strconv.FormatInt(in, 10) }

// base64CheckSum ...
func base64CheckSum(h [64]byte) string {
	return base64.StdEncoding.EncodeToString(h[:32])
}

// parseUnixTS
func parseUnixTS(s string) (time.Time, error) {
	t, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return _tsZero, errors.New("unable to parse Unix timestamp: " + err.Error())
	}
	return time.Unix(t, 0).UTC(), nil
}

//
// TOKEN
//

// const
const (
	// ui
	_app         = "[gitkeys] "
	_err         = "[error] "
	_keyFile     = "[KeyFile] [Read] "
	_urlFile     = "[URLFile] [Read] "
	_logFile     = "[LOGFile] [Read] "
	_walkerFile  = "[ConfigFile] [Walker] "
	_checksum    = "[CheckSum] "
	_compromised = "\n[WARNING] [CHECKSUM] [INVALID] consider data as compromised/modified "
	// shared parser & writer token
	_space             = " "
	_linefeed          = "\n"
	_doublelinefeed    = "\n\n"
	_sep               = " : "
	_field             = ", "
	_empty             = ""
	_slashfwd          = "/"
	_dot               = "."
	_log               = ".log"
	_urls              = ".urls"
	_keys              = ".keys"
	_hashmark          = '#'
	_comment           = "# "
	_https             = "https://"
	_commentline       = _comment + _linefeed
	_prefix            = _comment + _urlprefix
	_urlprefix         = "https://"
	_tsUnix            = "# timestamp unix     : "
	_tsZulu            = "# timestamp zulu     : "
	_integ             = "# CheckSum : "
	_globalCheckSum    = "# Global CheckSum    : "
	_checkSumKeyset    = "# CheckSum KeySet    : "
	_checkSumKeysetLen = len(_checkSumKeyset)
	_globalCheckSumLen = len(_globalCheckSum)
	_tsUnixLen         = len(_tsUnix)
	_tsZuluLen         = len(_tsZulu)
	_integLen          = len(_integ)
	_prefixLen         = len(_prefix)
	_minimumIDLen      = 2
	_prefixMinimum     = _prefixLen + _minimumIDLen
	_gitConf           = "/.git/config"
	_gitConfLen        = len(_gitConf)
)

// var
var (
	_worker     = runtime.NumCPU()
	_tsZero     = time.Time{}
	_emptySlice = []byte{}
)
