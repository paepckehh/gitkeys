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
