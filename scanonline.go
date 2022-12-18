// package gitkeys ...
package gitkeys

// import
import (
	"bytes"
	"errors"
	"sync"
	"time"
)

// const
const scaleConnect = 1 // number of parallel sessions per cpu core

// var
var (
	onlineOwnerChan                            = make(chan string, 100)
	onlineResultChan                           = make(chan *owner, 100)
	onlineCollect, onlineFeeder, onlineWorkers sync.WaitGroup
)

// scanGithub
func (r *Repo) scanOnline() error {
	ownerList := sortedOwnerMapIdx(&r.idMap)
	onlineCollect.Add(1)
	go onlineSpinUpCollect(r)
	onlineFeeder.Add(1)
	go onlineSpinUpFeeder(ownerList)
	onlineWorkers.Add(_worker * scaleConnect)
	go onlineSpinUpWorker(_worker * scaleConnect)
	onlineFeeder.Wait()
	close(onlineOwnerChan)
	onlineWorkers.Wait()
	close(onlineResultChan)
	onlineCollect.Wait()
	return nil
}

// onlineSpinUpFeeder
func onlineSpinUpFeeder(ownerList []string) {
	out(_app + "online repos to scan      : " + itoa(len(ownerList)) + _linefeed)
	for _, owner := range ownerList {
		onlineOwnerChan <- owner
	}
	onlineFeeder.Done()
}

// onlineSpinUpCollect
func onlineSpinUpCollect(r *Repo) {
	var changed, added int
	var buf bytes.Buffer
	for result := range onlineResultChan {
		if oldOwner, ok := r.idMap[result.name]; ok {
			if oldOwner.empty && result.empty {
				continue
			}
			result.checkSum()
			if bytes.Equal(oldOwner.checksum[:], result.checksum[:]) {
				continue
			}
			if !oldOwner.empty {
				changed++
				buf.WriteString("CHANGE KEYS  : " + result.name + _linefeed)
				buf.WriteString("OLD KEYS     : ")
				buf.WriteString(oldOwner.keySET.String() + _linefeed)
			} else {
				added++
			}
		}
		if !result.empty {
			buf.WriteString("NEW KEYS     : " + result.name + _linefeed)
			buf.WriteString(result.keySET.String() + _linefeed)
		}
		buf.WriteString(_linefeed)
		r.idMap[result.name] = result
	}
	out(_linefeed)
	s := buf.String()
	if len(s) > 0 {
		var log bytes.Buffer
		log.WriteString("#@ [" + time.Now().Format(time.RFC3339) + "]\n")
		log.WriteString("#@ [NEW/CHANGED KEYS FROM ONLINE REPO]\n")
		log.WriteString(s)
		r.changeLog = append(r.changeLog, log.String())
		out(_app + "github keys added         : " + itoa(added))
		out(_app + "github keys changed       : " + itoa(changed))
	}
	onlineCollect.Done()
}

// onlineSpinUpWorker
func onlineSpinUpWorker(worker int) {
	for i := 0; i < worker; i++ {
		go onlineWorker()
	}
}

// onlineWorker
func onlineWorker() {
	for o := range onlineOwnerChan {
		rawkeys, err := fetchRawKeys(o)
		if err != nil {
			continue
		}
		keySET, _, isEmpty := parseKeys(rawkeys)
		newOwner := owner{
			name:   o,
			empty:  isEmpty,
			keySET: keySET,
		}
		newOwner.checkSum()
		onlineResultChan <- &newOwner
		outPlain(_dot)
	}
	onlineWorkers.Done()
}

// fetchRawKeys
func fetchRawKeys(o string) ([]byte, error) {
	request, err := getRequest(_https + o + _keys)
	if err != nil {
		return _emptySlice, err
	}
	response, err := client.Do(request)
	if err != nil {
		return _emptySlice, err
	}
	if response.StatusCode != 200 {
		return _emptySlice, errors.New("fetch failed")
	}
	return decodeResponse(response)
}
