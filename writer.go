// package gitkeys ...
package gitkeys

// import
import (
	"bytes"
	"strings"
	"time"
)

//
// WRITER
//

// writeFiles ...
func (r *Repo) writeFiles() error {
	ts := time.Now()
	if err := r.writeKeyFile(ts); err != nil {
		if err != nil {
			out(_app + err.Error())
		}
	}
	if err := r.writeLogFile(ts); err != nil {
		if err != nil {
			out(_app + err.Error())
		}
	}
	if err := r.writeUrlFile(ts); err != nil {
		if err != nil {
			out(_app + err.Error())
		}
	}
	return nil
}

// writeLogFile ...
func (r *Repo) writeLogFile(ts time.Time) error {
	if len(r.changeLog) == 0 {
		return nil
	}
	out(_app + "writing clean log file    : " + r.KeyFile + _log)
	var buf bytes.Buffer
	writeHead(len(r.changeLog), ts, &buf)
	buf.WriteString(_doublelinefeed)
	// todo:
	// buf.WriteString("# Global CheckSum  : ")
	// buf.WriteString(...)
	// buf.WriteString(_linefeed)
	for _, u := range r.changeLog {
		buf.WriteString(u)
		buf.WriteString(_linefeed)
	}
	old := r.readLogFile()
	buf.WriteString(strings.Join(old, _linefeed))
	data := buf.Bytes()
	if err := saveWriteFile(r.KeyFile+_log, &data); err != nil {
		return err
	}
	return nil
}

// writeUrlFile ...
func (r *Repo) writeUrlFile(ts time.Time) error {
	out(_app + "writing clean url file    : " + r.KeyFile + _urls)
	var buf bytes.Buffer
	writeHead(len(r.urlMap), ts, &buf)
	buf.WriteString(_doublelinefeed)
	// todo:
	// buf.WriteString("# Global CheckSum  : ")
	// buf.WriteString(...)
	// buf.WriteString(_linefeed)
	if len(r.urlMap) > 0 {
		urls := sortedBoolMapIdx(&r.urlMap)
		for _, u := range urls {
			buf.WriteString(u)
			buf.WriteString(_linefeed)
		}
	}
	data := buf.Bytes()
	if err := saveWriteFile(r.KeyFile+_urls, &data); err != nil {
		return err
	}
	return nil
}

// writeKeyFile ...
func (r *Repo) writeKeyFile(ts time.Time) error {
	out(_app + "writing clean key file    : " + r.KeyFile)
	var buf bytes.Buffer
	writeHead(len(r.idMap), ts, &buf)
	buf.WriteString(_globalCheckSum)
	buf.WriteString(r.globalCheckSum(ts) + _linefeed)
	buf.WriteString(_doublelinefeed)
	if len(r.idMap) > 0 {
		ownerList := sortedOwnerMapIdx(&r.idMap)
		for _, o := range ownerList {
			owner := r.idMap[o]
			owner.checkSum()
			buf.WriteString(_prefix + o + _linefeed)
			buf.WriteString(_checkSumKeyset + base64CheckSum(owner.checksum) + _linefeed)
			if !owner.empty {
				buf.WriteString(owner.keySET.String())
			}
			buf.WriteString(_doublelinefeed)
		}
	}
	data := buf.Bytes()
	if err := saveWriteFile(r.KeyFile, &data); err != nil {
		return err
	}
	return nil
}

// newHead ...
func writeHead(items int, ts time.Time, buf *bytes.Buffer) {
	buf.WriteString("# [paepcke.de/gitkeys] DATABASE FILE -= DO NOT EDIT MANUALLY =-\n")
	buf.WriteString(_commentline)
	buf.WriteString(_tsUnix)
	buf.WriteString(itoa64(ts.Unix()) + _linefeed)
	buf.WriteString(_tsZulu)
	buf.WriteString(ts.Format(time.RFC3339) + _linefeed)
	buf.WriteString("# number of targets  : " + itoa(items) + _linefeed)
	buf.WriteString(_commentline)
	buf.WriteString("# WARNING The Global CheckSum covers only keys, urls, logs, timestamps!\n")
	buf.WriteString("# WARNING To get a clean and verified file, run 'gitkeys check.\n")
	buf.WriteString(_commentline)
}

// writeSignature ... todo as soon [ec or pq signature] formats race settled (signify,minify,age,ssh-plain,sphincs ...)
// func writeSignature() {
// 		buf.WriteString("# WARNING: The Global Signature covers only keys, urls, logs, timestamps!")
// 		buf.WriteString(_linefeed)
// 		buf.WriteString("# Global Signature  : ")
// 		buf.WriteString(_linefeed)
// }
