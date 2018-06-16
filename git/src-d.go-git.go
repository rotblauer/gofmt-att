package git

import (
	"fmt"
	"github.com/rotblauer/gofmt-att/remote"
	"gopkg.in/src-d/go-git.v4"
	gitconfig "gopkg.in/src-d/go-git.v4/config"
	gitobject "gopkg.in/src-d/go-git.v4/plumbing/object"
	gitplumbing "gopkg.in/src-d/go-git.v4/plumbing"
	"path/filepath"
	"time"
	"os"
	"errors"
	"net/http"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/client"
githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"os/exec"
	"io"
	"bytes"
	"strings"
)

type GoGit struct {
	BasePath string
}

var defaultRemoteName = "fmtatt"

func NewGoGit() GoGit {
	return GoGit{}
}

func (g GoGit) SetBase(ppath string) (err error) {
	p := filepath.Clean(ppath)
	p, err = filepath.Abs(p)
	g.BasePath = p
	return
}

func (g GoGit) SetClient(c *http.Client) error {
	// Override http(s) default protocol to use our custom client
	client.InstallProtocol("https", githttp.NewClient(c))
	return nil
}

var ErrDirExists = errors.New("dir exists")

func (g GoGit) Clone(repo *remote.RepoT) (err error) {
	opts := &git.CloneOptions{
		URL:               repo.GitUrl,
		Auth:              nil,
		// RemoteName:        "origin",
		// ReferenceName:     "master",
		SingleBranch:      false,
		NoCheckout:        false,
		Depth:             0,
		RecurseSubmodules: 0,
		Progress:          nil,
		Tags:              0,
	}
	ps, _ := filepath.Split(repo.Target)
	if len(ps) < 3 {
		panic("short repo path")
	}
	if mkerr := os.MkdirAll(filepath.Dir(repo.Target), os.ModePerm); mkerr != nil {
		err = mkerr
		return
	}

	_, existsErr := os.Stat(repo.Target)
	if existsErr == nil {
		err = os.RemoveAll(repo.Target)
		if err != nil {
			return
		}
	}
	_, err = git.PlainClone(repo.Target, false, opts)
	return
}

func (g GoGit) Status(dirPath string) (dirty bool, statusString string, err error) {
	r, err := git.PlainOpen(dirPath)
	if err != nil {
		return
	}
	wt, err := r.Worktree()
	if err != nil {
		return
	}
	status, err := wt.Status()
	if err != nil {
		return
	}
	statusString = status.String()
	dirty = !status.IsClean()
	return
}

func (g GoGit) Add(dirPath, filePath string) (err error) {
	r, err := git.PlainOpen(dirPath)
	if err != nil {
		return
	}
	wt, err := r.Worktree()
	if err != nil {
		return
	}
	_, err = wt.Add(filePath)
	if err != nil {
		return
	}
	return
}

func (g GoGit) CommitWithBranch(dirPath string, commit *GitCommitConfig) (hash, status string, err error) {
	r, err := git.PlainOpen(dirPath)
	if err != nil {
		return
	}
	wt, err := r.Worktree()
	if err != nil {
		return
	}
	// get status
	stat, err := wt.Status()
	if err != nil {
		return
	}
	// set status for returnable (this will show the staged/unstaged changes)
	status = stat.String()

	// make commit
	c, err := wt.Commit(fmt.Sprintf(`%s

%s
`, commit.Title, commit.Body), &git.CommitOptions{
		Author: &gitobject.Signature{
			Name:  commit.AuthorName,
			Email: commit.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return
	}
	hash = c.String()

	// set up branch ref name
	refName := gitplumbing.ReferenceName("refs/heads/"+commit.BranchName)

	// create a branch
	b := &gitconfig.Branch{
		Name:   commit.BranchName,
		Remote: "origin",
		Merge:  refName,
	}
	err = b.Validate()
	if err != nil {
		return
	}
	err = r.CreateBranch(b)
	if err != nil {
		return
	}
	if !refName.IsBranch() {
		err = errors.New("failed to create branch "+ commit.BranchName)
		return
	}
	err = wt.Checkout(&git.CheckoutOptions{
		// Hash:   c,
		Branch: refName,
		Create: true, // useless, tho, i think
		Force:  true, // gotta throw away any unstaged changes
	})
	return
}

func (g GoGit) PushAll(dirPath, remote, branchName string) (err error) {
	// r, err := git.PlainOpen(dirPath)
	// if err != nil {
	// 	return
	// }
	// opts := &git.PushOptions{
	// 	RemoteName: defaultRemoteName,
	// 	RefSpecs:   []gitconfig.RefSpec{"+refs/heads/*:refs/remotes/origin/*"},
	// 	Auth:       nil,
	// 	Progress:   nil,
	// }
	// err = opts.Validate()
	// if err != nil {
	// 	return
	// }
	//
	// return r.Set(opts)

	cmdRemote := []string{"git", "remote", "add", defaultRemoteName, remote}
	fmt.Println("running: $", strings.Join(cmdRemote, " "))
	cmdR := exec.Command(cmdRemote[0], cmdRemote[1:]...)
	cmdR.Dir = dirPath
	outR, err := cmdR.CombinedOutput()
	io.Copy(os.Stderr, bytes.NewBuffer(outR))

	command := []string{"git", "push", "-u", defaultRemoteName, branchName}
	fmt.Println("running: $", strings.Join(command, " "))
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = dirPath
	out, err := cmd.CombinedOutput()
	io.Copy(os.Stderr, bytes.NewBuffer(out))
	return
}
