package gitkeys

import (
	"bufio"
	"bytes"
	"io/fs"
	"net/url"
	"os"
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
	feeder.Add(1)
	go spinUpFeeder(r.GitStore)
	confWorker.Add(_worker)
	go spinUpConfWorker(r.GitStore)
	feeder.Wait()
	close(configFileChan)
	confWorker.Wait()
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

// spinUpFeeder ...
func spinUpFeeder(repos string) {
	fsRoot := os.DirFS(repos)
	err := fs.WalkDir(fsRoot, _dot, feederWalkFn)
	if err != nil {
		out(_walkerFile + "unable to walk: " + repos + _sep + err.Error())
	}
	feeder.Done()
}

// spinUpConfWorker ...
func spinUpConfWorker(repos string) {
	for i := 0; i < _worker; i++ {
		go func() {
			for filename := range configFileChan {
				filename = repos + _slashfwd + filename
				f, err := os.Open(filename)
				if err != nil {
					out("faild to open: " + filename)
					continue
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
			}
			confWorker.Done()
		}()
	}
}

// feederWalkFn ...
func feederWalkFn(path string, d fs.DirEntry, err error) error {
	_ = d
	if err != nil {
		out("unable to walk:" + path)
	}
	l := len(path)
	if l > _gitConfLen {
		if path[l-_gitConfLen:] == _gitConf {
			configFileChan <- path
		}
	}
	return nil
}
