package fmtatt

import (
	"io"
)

type GitProvider interface {
	Clone(remote SimpleRemoteT) (dirPath string, err error) // Path to repo (name), this will be used by the FmtConfig.Commands and FmtConfig.Target
	IsDirty(dirPath string) bool
	CreateBranch(dirPath string, branchName string) error
	AddRemote(dirPath string, remoteName string) error
	AddAndCommitAll(dirPath string, commit SimpleCommitConfig) (diff io.Reader, err error)
	PushAll(dirPath string, remote string) error
}