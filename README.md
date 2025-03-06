# OVERVIEW
[![Go Reference](https://pkg.go.dev/badge/paepcke.de/gitkeys.svg)](https://pkg.go.dev/paepcke.de/gitkeys) [![Go Report Card](https://goreportcard.com/badge/paepcke.de/gitkeys)](https://goreportcard.com/report/paepcke.de/gitkeys) [![Go Build](https://github.com/paepckehh/gitkeys/actions/workflows/golang.yml/badge.svg)](https://github.com/paepckehh/gitkeys/actions/workflows/golang.yml)

[paepcke.de/gitkeys](https://paepcke.de/gitkeys/)
log-store : [paepcke.de/keys](https://paepcke.de/keys/)

git ssh keys logging, stupid simple, fast, local-first
 
- perfect companion for age-encryption (have always up-to-date trusted keys)
- easy to use & review (hash/checksum/protected) clear text database files
- verify all [ssh-key] signed commits, tags, files - yourself, locally, offline 
- all files are add/append only: we never remove any entries from keys, keys.urls or keys.log
- all key sets in the keyfile (and the keyfile itself) is protected by (chained) sha512 hash checksums (wip:signatures)
- 100 % pure go, 100 % pure stdlib only, no external dependencies

# RUN

```
go run paepcke.de/gitkeys/cmd/gitkeys@latest
```

# INSTALL

```
go install paepcke.de/gitkeys/cmd/gitkeys@latest
```

### DOWNLOAD (prebuild)

[github.com/paepckehh/gitkeys/releases](https://github.com/paepckehh/gitkeys/releases)

# SHOWTIME

## Want to use and verify my example store? 

``` Shell
git clone https://github.com/paepckehh/keys
cd keys
GITSTORE="." go run paepcke.de/gitkeys/cmd/gitkeys@latest fetch
[gitkeys] SSH Key Transparency Log  : Mode Check [CheckInteg] [AddLocal] [CleanRewrite]
[gitkeys] key file stats            : /usr/store/git/.keys => owner total : 1117
[gitkeys] url file stats            : /usr/store/git/.keys.urls => urls total : 1926
[gitkeys] scan local store for urls : /usr/store/git
[gitkeys] writing clean key file    : /usr/store/git/.keys
[gitkeys] writing clean url file    : /usr/store/git/.keys.urls
```

## Want your own local store? 

``` Shell
GITSTORE="/tmp/keystore"
echo "https://github.com/klauspost" > $GITSTORE/.keys.urls
go run paepcke.de/gitkeys/cmd/gitkeys@latest fetch
```

# DOCS

[pkg.go.dev/paepcke.de/gitkeys](https://pkg.go.dev/paepcke.de/gitkeys)

# CONTRIBUTION

Yes, Please! PRs Welcome! 
