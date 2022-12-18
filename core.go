// package gitkeys ...
package gitkeys

// import
import (
	"net/url"
	"runtime"
	"sort"
	"strings"
	"time"
)

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

// addOwner ...
func (r *Repo) addOwner(ownerName, keys string) {
	keySET, _, isEmpty := parseKeys([]byte(keys))
	newOwner := &owner{
		name:   ownerName,
		empty:  isEmpty,
		keySET: keySET,
	}
	newOwner.checkSum()
	r.idMap[ownerName] = newOwner
}

// addEmptyOwner ...
func (r *Repo) addEmptyOwner(ownerName string) {
	newOwner := &owner{
		name:  ownerName,
		empty: true,
	}
	newOwner.checkSum()
	r.idMap[ownerName] = newOwner
}

// getOwnerUrl ...
func getOwnerUrl(owner string) string {
	return _urlprefix + owner
}

// getRepoOwner ...
func getRepoOwner(u *url.URL) string {
	var owner string
	if u.Scheme == "https" && strings.Contains(u.Host, _dot) {
		s := strings.Split(u.Path, _slashfwd)
		if len(s) > 0 {
			owner = u.Host + _slashfwd + s[1]
		}
	}
	return owner
}

// sortedOwnerMapIdx ...
func sortedOwnerMapIdx(ownerMap *map[string]*owner) []string {
	var s []string
	for k := range *ownerMap {
		s = append(s, k)
	}
	sort.Strings(s)
	return s
}

// sortedBoolMapIdx ...
func sortedBoolMapIdx(ids *map[string]bool) []string {
	var s []string
	for k := range *ids {
		s = append(s, k)
	}
	sort.Strings(s)
	return s
}
