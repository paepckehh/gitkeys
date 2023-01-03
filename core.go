// package gitkeys ...
package gitkeys

// import
import (
	"net/url"
	"sort"
	"strings"
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
