// package gitkeys ...
package gitkeys

// import
import (
	"encoding/base64"
	"errors"
	"os"
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
