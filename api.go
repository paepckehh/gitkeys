// package gitkeys ...
package gitkeys

// import
import ()

// Repo ...
type Repo struct {
	KeyFile   string
	GitStore  string
	urlMap    map[string]bool
	idMap     map[string]*owner
	changeLog []string
	Verbose   bool
}

// NewRepo ...
func NewRepo() *Repo {
	return &Repo{
		urlMap: make(map[string]bool),
		idMap:  make(map[string]*owner),
	}
}

// Pinonly re-calculates the Pin (checksum) of all keys within KeyFile.
// If checksum is valid, return a verified checksum-pin of KeyFile.
func (r *Repo) Pinonly() string {
	if r.KeyFile == _empty {
		outERR("Need KeyFile")
		return _empty
	}
	return r.readKeyFile()
}

// Check reads and verifies (when exists) KeyFiles (keys and urls) and scans (when provided) the git repo (read-only) for urls.
func (r *Repo) Check() error {
	outINFO("SSH Key Transparency Log  : Mode Check [CheckInteg] [AddLocal] [CleanRewrite]")
	if err := r.readFiles(); err != nil {
		return err
	}
	if err := r.scanRepoStore(); err != nil {
		return err
	}
	if err := r.writeFiles(); err != nil {
		return err
	}
	return nil
}

// Fetch performes the same actions as Update, but does an online keyfetch download/diff as well.
func (r *Repo) Fetch() error {
	outINFO("SSH Key Transparency Log  : Mode Fetch => [CheckInteg] [AddLocal] [AddOnline] [CleanRewrite]")
	if err := r.readFiles(); err != nil {
		return err
	}
	if err := r.scanRepoStore(); err != nil {
		return err
	}
	if err := r.scanOnline(); err != nil {
		return err
	}
	if err := r.writeFiles(); err != nil {
		return err
	}
	return nil
}
