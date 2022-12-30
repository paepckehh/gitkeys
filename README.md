# OVERVIEW 

[paepche.de/gitkeys](https://paepcke.de/gitkeys)

git ssh keys logging , stupid simple, fast, local
 
- perfect companion for age-encryption (have always up-to-date trusted keys)
- easy to use & review (hash/checksum/protected) clear text database files
- verify all [ssh-key] signed commits, tags, files - yourself, locally, offline 
- all files are add/append only: we never remove any entries from keys, keys.urls or keys.log
- all key sets in the keyfile (and the keyfile itself) is protected by (chained) sha512 hash checksums (wip:signatures)
- 100 % pure go, 100 % pure stdlib only, no external dependencies

# INSTALL

```
go install paepcke.de/gitkeys/cmd/gitkeys@latest
```

### DOWNLOAD (prebuild)

[github.com/paepckehh/gitkeys/releases](https://github.com/paepckehh/gitkeys/releases)

# SHOWTIME

## Do you have a store of local git (mirrors)? 

``` Shell
GITSTORE="/usr/store/git" gitkeys fetch
[gitkeys] SSH Key Transparency Log  : Mode Check [CheckInteg] [AddLocal] [CleanRewrite]
[gitkeys] key file stats            : /usr/store/git/.keys => owner total : 1117
[gitkeys] url file stats            : /usr/store/git/.keys.urls => urls total : 1926
[gitkeys] scan local store for urls : /usr/store/git
[gitkeys] writing clean key file    : /usr/store/git/.keys
[gitkeys] writing clean url file    : /usr/store/git/.keys.urls
```

## Do you have a list of git repo urls?

``` Shell
echo "https://github.com/klauspost" > /usr/store/git/.keys.urls
GITSTORE="/usr/store/git" gitkeys fetch
[...]
```

## Do you have a existing keys file that you want to update?

``` Shell
GITSTORE="/usr/store/git" gitkeys fetch
[...]
```

## Do you have a existing keys file that you want to integ check, verify, clean-rewrite only?

``` Shell
GITSTORE="/usr/store/git" gitkeys
[...]
```

## Behind a (corp|security) proxy?

``` Shell
HTTPS_PROXY="proxy.bigCorp.local" SSL_CERT_FILE="/etc/ssl/bigCorpProxy.pem" GITSTORE="/usr/store/git" gitkeys fetch
 [...]
```

# CONTRIBUTION

Yes, Please! PRs Welcome! 
