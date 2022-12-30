// package gitkeys ...
package gitkeys

// import
import (
	"bufio"
	"bytes"
	"errors"
	"net/url"
	"os"
	"time"
)

// readFiles ...
func (r *Repo) readFiles() error {
	if r.GitStore == _empty || r.KeyFile == _empty {
		return errors.New("Need RepoStore and KeyFile")
	}
	_ = r.readKeyFile()
	r.readUrlFile()
	return nil
}

// readKeyFile ...
func (r *Repo) readKeyFile() string {
	var err error
	if !isFile(r.KeyFile) {
		r.idMap = make(map[string]*owner)
		return _empty
	}
	f, err := os.Open(r.KeyFile)
	if err != nil {
		outERR(_keyFile + r.KeyFile + _sep + err.Error())
		return _empty
	}
	var counter int
	scanner := bufio.NewScanner(f)
	ownerName, keys, globalCheckSum, activeSET, tsUnix, tsZulu := "", "", "", false, time.Time{}, time.Time{}
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case len(line) < _prefixMinimum:
			continue
		case activeSET:
			switch {
			case line[:_prefixLen] == _prefix:
				counter++
				if keys == _empty {
					r.addEmptyOwner(ownerName)
				} else {
					r.addOwner(ownerName, keys)
					keys = _empty
				}
				ownerName, activeSET = line[_prefixLen:], true
			default:
				if line != _empty && line[0] != _hashmark {
					keys += line + _linefeed
				}
			}
		case line[:_prefixLen] == _prefix:
			ownerName, keys, activeSET = line[_prefixLen:], _empty, true
		case line[:_globalCheckSumLen] == _globalCheckSum:
			if globalCheckSum != "" {
				outERR(_compromised)
				outERR(_keyFile + _checksum + "duplicate integ line")
			}
			globalCheckSum = line[_globalCheckSumLen:]
		case line[:_tsUnixLen] == _tsUnix:
			if !tsUnix.Equal(_tsZero) {
				outERR(_compromised)
				outERR(_keyFile + _checksum + "duplicate unix timestamp")
			}
			if tsUnix, err = parseUnixTS(line[_tsUnixLen:]); err != nil {
				outERR(_compromised)
				outERR(_keyFile + _checksum + "parse timestamp unix : " + err.Error())
			}
		case line[:_tsZuluLen] == _tsZulu:
			if !tsZulu.Equal(_tsZero) {
				outERR(_compromised)
				outERR(_keyFile + _checksum + "duplicate zulu timestamp")
			}
			if tsZulu, err = time.Parse(time.RFC3339, line[_tsZuluLen:]); err != nil {
				outERR(_compromised)
				outERR(_keyFile + _checksum + "parse timestamp zulu: " + err.Error())
			}
		}
	}
	if ownerName != _empty {
		counter++
		if keys == _empty {
			r.addEmptyOwner(ownerName)
		} else {
			r.addOwner(ownerName, keys)
			keys = ""
		}
	}
	if !tsUnix.Equal(tsZulu) {
		outERR(_compromised)
		outERR(_keyFile + _checksum + "timestamps [unix|zulu] missmatch")
	}
	if globalCheckSum == _empty {
		outERR(_compromised)
		outERR(_keyFile + _checksum + "global checksum fail: no checksum found in keyfile")
	}
	if !r.verifyGlobalInteg(tsUnix, globalCheckSum) {
		outERR(_compromised)
		outERR(_keyFile + _checksum + "global checksum fail")
	}
	outINFO("key file stats            : " + r.KeyFile + " => owner total : " + itoa(counter))
	return globalCheckSum
}

// readUrlFile ...
func (r *Repo) readUrlFile() {
	filename := r.KeyFile + _urls
	if !isFile(filename) {
		return
	}
	f, err := os.Open(filename)
	if err != nil {
		outERR(_urlFile + filename + _sep + err.Error())
		return
	}
	var counter int
	var newUrl []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 1 {
			continue
		}
		if line[0] == '#' {
			continue
		}
		if u, err := url.Parse(line); err == nil {
			counter++
			r.urlMap[u.String()] = true
			ownerName := getRepoOwner(u)
			if ownerName != _empty {
				if _, ok := r.idMap[ownerName]; !ok {
					newUrl = append(newUrl, ownerName)
					r.addEmptyOwner(ownerName)
				}
			}
		}
	}
	l := len(newUrl)
	if l > 0 {
		var log bytes.Buffer
		log.WriteString("#@ [" + time.Now().Format(time.RFC3339) + "]\n")
		log.WriteString("#@ [NEW URLS FROM URL FILE] [" + r.GitStore + "]\n")
		for _, e := range newUrl {
			log.WriteString(e)
			log.WriteString(_field)
		}
		r.changeLog = append(r.changeLog, log.String())
		out(_app + "urls added from url file  : " + itoa(l))
	}
	outINFO("url file stats            : " + r.KeyFile + _urls + " => urls total : " + itoa(counter))
}

// readLogFile ...
func (r *Repo) readLogFile() []string {
	filename := r.KeyFile + _log
	if !isFile(filename) {
		return []string{}
	}
	f, err := os.Open(filename)
	if err != nil {
		outERR(_logFile + filename + _sep + err.Error())
		return []string{}
	}
	var counter int
	var log []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case len(line) < 2:
			continue
		case line[0] == '#' && line[1] != '@':
			continue
		}
		counter++
		log = append(log, line)
	}
	outINFO("reading log file          : " + r.KeyFile + _log + " => existing logfile lines : " + itoa(counter))
	return log
}
