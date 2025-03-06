package gitkeys

import (
	"bufio"
	"bytes"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

//
// REPOS
//

// type
type scanResult struct {
	url   string
	owner string
}

// var
var (
	configFileChan              = make(chan string, 25)
	urlCollectChan              = make(chan *scanResult, 100)
	feeder, collect, confWorker sync.WaitGroup
)

// scanRepos ...
func (r *Repo) scanRepoStore() error {
	out(_app + "scan local store for urls : " + r.GitStore)
	collect.Add(1)
	go spinUpCollector(r)
	confWorker.Add(_worker)
	go spinUpConfWorker(r.GitStore)
	feeder.Wait()
	close(urlCollectChan)
	collect.Wait()
	return nil
}

// spinUpCollector ...
func spinUpCollector(r *Repo) {
	var newGit, newUrl []string
	for s := range urlCollectChan {
		if s.url != _empty && !r.urlMap[s.url] {
			r.urlMap[s.url] = true
			newUrl = append(newUrl, s.url)
		}
		if s.owner != _empty {
			if _, ok := r.idMap[s.owner]; !ok {
				r.addEmptyOwner(s.owner)
				newGit = append(newGit, s.owner)
			}
		}
	}
	urlLen := len(newUrl)
	if urlLen > 0 {
		var log bytes.Buffer
		log.WriteString("#@ [" + time.Now().Format(time.RFC3339) + "]\n")
		log.WriteString("#@ [NEW URLS FROM LOCAL GIT STORE] [" + r.GitStore + "]\n")
		for _, e := range newUrl {
			log.WriteString(e)
			log.WriteString(_field)
		}
		r.changeLog = append(r.changeLog, log.String())
		out(_app + "added from git store urls : " + itoa(urlLen))
	}
	gitLen := len(newGit)
	if gitLen > 0 {
		var log bytes.Buffer
		log.WriteString("#@ [" + time.Now().Format(time.RFC3339) + "]\n")
		log.WriteString("#@ [NEW OWNER FROM LOCAL GIT STORE] [" + r.GitStore + "]\n")
		for _, e := range newGit {
			log.WriteString(e)
			log.WriteString(_field)
		}
		r.changeLog = append(r.changeLog, log.String())
		out(_app + "added from git store owner: " + itoa(gitLen))
	}
	collect.Done()
}

// spinUpConfWorker ...
func spinUpConfWorker(repo string) {
	go func() {
		filename := filepath.Join(repo, "keys")
		f, err := os.Open(filename)
		if err != nil {
			log.Fatal("failed to open keys file: " + filename)
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "url = http") {
				s := strings.Split(line, "url = ")
				if len(s) > 1 {
					u, err := url.Parse(s[1])
					if err != nil {
						out("faild to parse" + err.Error())
						continue
					}
					urlCollectChan <- &scanResult{u.String(), getRepoOwner(u)}
				}
			}
		}
		f.Close()
		confWorker.Done()
	}()
}
