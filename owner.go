package gitkeys

import (
	"bytes"
	"crypto/sha512"
)

// owner ...
type owner struct {
	name     string
	empty    bool
	keySET   *keySET
	checksum [sha512.Size]byte
}

// checkSum ...
func (o *owner) checkSum() {
	var buf bytes.Buffer
	if o.name == _empty {
		out("[internal-error] [checkSum] owner empty")
	}
	buf.WriteString(o.name)
	if !o.empty {
		buf.WriteString(o.keySET.String())
	}
	o.checksum = sha512.Sum512(buf.Bytes())
}

// getOwnerUrl ...
// func getOwnerUrl(owner string) string {
// 	return _urlprefix + owner
// }

// newOwner ..
// func newOwner() *owner {
// 	return &owner{
// 		empty:    true,
// 		keySET:   newKeySET(),
// 		checksum: [sha512.Size]byte{},
// 	}
// }

// valid ...
// func (o *owner) valid(newCheckSum [sha512.Size]byte) bool {
// 	o.checkSum()
// 	return bytes.Equal(o.checksum[:], newCheckSum[:])
// }
