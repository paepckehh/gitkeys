package gitkeys

import (
	"bufio"
	"bytes"
	"crypto/sha512"
	"strings"
	"time"
)

// keySET ...
type keySET struct {
	kRSA        []string
	kDSA        []string
	kECDSA256   []string
	kECDSA384   []string
	kECDSA521   []string
	kED25519    []string
	kSKED25519  []string
	kSKECDSA256 []string
}

// const algo
const (
	aRSA        = "ssh-rsa"
	aDSA        = "ssh-dss"
	aECDSA256   = "ecdsa-sha2-nistp256"
	aECDSA384   = "ecdsa-sha2-nistp384"
	aECDSA521   = "ecdsa-sha2-nistp521"
	aED25519    = "ssh-ed25519"
	aSKED25519  = "sk-ssh-ed25519@openssh.com"
	aSKECDSA256 = "sk-ecdsa-sha2-nistp256@openssh.com"
)

// newKeySET ...
func newKeySET() *keySET {
	return &keySET{
		kRSA:        make([]string, 0),
		kDSA:        make([]string, 0),
		kECDSA256:   make([]string, 0),
		kECDSA384:   make([]string, 0),
		kECDSA521:   make([]string, 0),
		kED25519:    make([]string, 0),
		kSKED25519:  make([]string, 0),
		kSKECDSA256: make([]string, 0),
	}
}

// verifyGlobalInteg ...
func (r *Repo) verifyGlobalInteg(ts time.Time, integ string) bool {
	// out("DEBUG INTEG TARGET:" + integ)
	// out("DEBUG INTEG CALC  :" + r.globalCheckSum(ts))
	return integ == r.globalCheckSum(ts)
}

// globalCheckSum ...
func (r *Repo) globalCheckSum(ts time.Time) string {
	var buf bytes.Buffer
	buf.WriteString(itoa64(ts.Unix()))
	s := sortedOwnerMapIdx(&r.idMap)
	for _, v := range s {
		owner := r.idMap[v]
		owner.checkSum()
		buf.WriteString(base64CheckSum(owner.checksum))
	}
	// out("### CALC DEBUG DATA : " + buf.String())
	// out("### CALC DEBUG CSUM : " + base64CheckSum(sha512.Sum512(buf.Bytes())))
	return base64CheckSum(sha512.Sum512(buf.Bytes()))
}

// parseKeys ...
func parseKeys(body []byte) (*keySET, string, bool) {
	var short strings.Builder
	k, isEmpty := newKeySET(), true
	scanner := bufio.NewScanner(bytes.NewBuffer(body))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 10 {
			continue
		}
		s := strings.Split(line, _space)
		if len(s) != 2 {
			continue
		}
		switch s[0] {
		case aRSA:
			short.WriteString("r")
			k.kRSA, isEmpty = append(k.kRSA, s[1]), false
		case aDSA:
			short.WriteString("d")
			k.kDSA, isEmpty = append(k.kDSA, s[1]), false
		case aECDSA256:
			short.WriteString("e")
			k.kECDSA256, isEmpty = append(k.kECDSA256, s[1]), false
		case aSKECDSA256:
			short.WriteString("e")
			k.kSKECDSA256, isEmpty = append(k.kSKECDSA256, s[1]), false
		case aECDSA384:
			short.WriteString("e")
			k.kECDSA384, isEmpty = append(k.kECDSA384, s[1]), false
		case aECDSA521:
			short.WriteString("e")
			k.kECDSA521, isEmpty = append(k.kECDSA521, s[1]), false
		case aED25519:
			short.WriteString("ED")
			k.kED25519, isEmpty = append(k.kED25519, s[1]), false
		case aSKED25519:
			short.WriteString("ED")
			k.kSKED25519, isEmpty = append(k.kSKED25519, s[1]), false
		default:
			continue // skip garbage
		}
	}
	return k, short.String(), isEmpty
}

// checkSum ...
func (k *keySET) String() string {
	var buf bytes.Buffer
	if len(k.kRSA) > 0 {
		for _, key := range k.kRSA {
			buf.WriteString(aRSA)
			buf.WriteString(_space)
			buf.WriteString(key)
			buf.WriteString(_linefeed)
		}
	}
	if len(k.kDSA) > 0 {
		for _, key := range k.kDSA {
			buf.WriteString(aDSA)
			buf.WriteString(_space)
			buf.WriteString(key)
			buf.WriteString(_linefeed)
		}
	}
	if len(k.kECDSA256) > 0 {
		for _, key := range k.kECDSA256 {
			buf.WriteString(aECDSA256)
			buf.WriteString(_space)
			buf.WriteString(key)
			buf.WriteString(_linefeed)
		}
	}
	if len(k.kECDSA384) > 0 {
		for _, key := range k.kECDSA384 {
			buf.WriteString(aECDSA384)
			buf.WriteString(_space)
			buf.WriteString(key)
			buf.WriteString(_linefeed)
		}
	}
	if len(k.kECDSA521) > 0 {
		for _, key := range k.kECDSA521 {
			buf.WriteString(aECDSA521)
			buf.WriteString(_space)
			buf.WriteString(key)
			buf.WriteString(_linefeed)
		}
	}
	if len(k.kED25519) > 0 {
		for _, key := range k.kED25519 {
			buf.WriteString(aED25519)
			buf.WriteString(_space)
			buf.WriteString(key)
			buf.WriteString(_linefeed)
		}
	}
	if len(k.kSKED25519) > 0 {
		for _, key := range k.kSKED25519 {
			buf.WriteString(aSKED25519)
			buf.WriteString(_space)
			buf.WriteString(key)
			buf.WriteString(_linefeed)
		}
	}
	if len(k.kSKECDSA256) > 0 {
		for _, key := range k.kSKECDSA256 {
			buf.WriteString(aSKECDSA256)
			buf.WriteString(_space)
			buf.WriteString(key)
			buf.WriteString(_linefeed)
		}
	}
	return buf.String()
}
