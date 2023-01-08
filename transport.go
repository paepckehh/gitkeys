package gitkeys

import (
	"compress/gzip"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/url"
)

const (
	_userAgent = "curl"
	// _uaApp        = "Mozilla/5.0 (X11; CrOS aarch64 13597.84.0) "
	// _uaFramework  = "AppleWebKit/537.36 (KHTML, like Gecko) "
	// _uaOS         = "Chrome/104.0.5112.105 Safari/537.36"
	// _userAgent    = _uaApp + _uaFramework + _uaOS
)

// uncomment to enable & manage keypin (every 180 days)
// const _githubKeyPin = "/3ftdeWqIAONye/CeEQuLGvtlw4MPnQmKgyPLugFbK8="
// var client = getClient(getTransport(tlsconf(_githubKeyPin)))
var client = getClient(getTransport(getTlsConf(_empty)))

// getTlsConf ...
func getTlsConf(keyPin string) *tls.Config {
	tlsConfig := &tls.Config{
		InsecureSkipVerify:     false,
		SessionTicketsDisabled: true,
		Renegotiation:          0,
		MinVersion:             tls.VersionTLS13,
		MaxVersion:             tls.VersionTLS13,
		CipherSuites:           []uint16{tls.TLS_CHACHA20_POLY1305_SHA256},
		CurvePreferences:       []tls.CurveID{tls.X25519},
	}
	if keyPin != _empty {
		tlsConfig.VerifyConnection = func(state tls.ConnectionState) error {
			if !pinVerifyState(keyPin, &state) {
				return errors.New("[error] keypin invalid: " + keyPin)
			}
			return nil
		}
	}
	return tlsConfig
}

// pinVerifyState ...
func pinVerifyState(keyPin string, state *tls.ConnectionState) bool {
	if len(state.PeerCertificates) > 0 {
		if keyPin == keyPinBase64(state.PeerCertificates[0]) {
			return true
		}
	}
	return false
}

// keyPinBase64 ...
func keyPinBase64(cert *x509.Certificate) string {
	h := sha256.Sum256(cert.RawSubjectPublicKeyInfo)
	return base64.StdEncoding.EncodeToString(h[:])
}

// getClient
func getClient(transport *http.Transport) *http.Client {
	return &http.Client{
		CheckRedirect: nil,
		Jar:           nil,
		Transport:     transport,
	}
}

// getTransport ...
func getTransport(tlsconf *tls.Config) *http.Transport {
	return &http.Transport{
		Proxy:              http.ProxyFromEnvironment,
		TLSClientConfig:    tlsconf,
		DisableCompression: true,
		ForceAttemptHTTP2:  false,
	}
}

// getRequest ...
func getRequest(plainURL string) (*http.Request, error) {
	targetURL, err := url.Parse(plainURL)
	if err != nil {
		return &http.Request{}, errors.New("invalid url")
	}
	return &http.Request{
		URL:        targetURL,
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header: http.Header{
			"User-Agent":      []string{_userAgent},
			"Accept-Encoding": []string{"gzip"},
		},
	}, nil
}

// decodeResponse ...
func decodeResponse(resp *http.Response) ([]byte, error) {
	var err error
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
	default:
		reader = resp.Body
	}
	if err != nil {
		return nil, errors.New("decode: " + err.Error())
	}
	return io.ReadAll(reader)
}
