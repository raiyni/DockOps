package git

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pkg/errors"
)

type gitClient struct {
}

type BaseOption struct {
	RepositoryUrl string
}

type BasicAuth struct {
	Username string
	Password string
}

type SshAuth struct {
	KeyPath     string
	KeyPassword string
}

type FetchOption struct {
	BaseOption
	BasicAuth
	SshAuth
	ReferenceName string
}

func NewGitClient() *gitClient {
	return &gitClient{}
}

func (c *gitClient) LatestHashSsh(ctx context.Context, opt FetchOption) (string, error) {
	return c.latestCommitSsh(ctx, opt)
}

func (c *gitClient) LatestHashHttp(ctx context.Context, opt FetchOption) (string, error) {
	return c.latestCommitHttp(ctx, opt)
}

func (c *gitClient) latestCommitSsh(ctx context.Context, opt FetchOption) (string, error) {
	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{opt.RepositoryUrl},
	})

	key, err := getSshAuth(opt.KeyPath, opt.KeyPassword)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	listOptions := &git.ListOptions{
		Auth: key,
	}

	refs, err := remote.List(listOptions)
	if err != nil {
		return "", errors.Wrap(err, "failed to list repository refs")
	}

	referenceName := opt.ReferenceName
	if referenceName == "" {
		for _, ref := range refs {
			if strings.EqualFold(ref.Name().String(), "HEAD") {
				referenceName = ref.Target().String()
			}
		}
	}

	for _, ref := range refs {
		if strings.EqualFold(ref.Name().String(), referenceName) {
			return ref.Hash().String(), nil
		}
	}

	return "", errors.Errorf("could not find ref %q in the repository", opt.ReferenceName)
}

func (c *gitClient) latestCommitHttp(ctx context.Context, opt FetchOption) (string, error) {
	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{opt.RepositoryUrl},
	})

	listOptions := &git.ListOptions{
		Auth: getBasicAuth(opt.Username, opt.Password),
	}

	refs, err := remote.List(listOptions)
	if err != nil {
		if err.Error() == "authentication required" {
			return "", err
		}
		return "", errors.Wrap(err, "failed to list repository refs")
	}

	referenceName := opt.ReferenceName
	if referenceName == "" {
		for _, ref := range refs {
			if strings.EqualFold(ref.Name().String(), "HEAD") {
				referenceName = ref.Target().String()
			}
		}
	}

	for _, ref := range refs {
		if strings.EqualFold(ref.Name().String(), referenceName) {
			return ref.Hash().String(), nil
		}
	}

	return "", errors.Errorf("could not find ref %q in the repository", opt.ReferenceName)
}

func getBasicAuth(username, password string) *githttp.BasicAuth {
	if password != "" {
		if username == "" {
			username = "token"
		}

		return &githttp.BasicAuth{
			Username: username,
			Password: password,
		}
	}
	return nil
}

func getSshAuth(path, password string) (*ssh.PublicKeys, error) {
	return ssh.NewPublicKeysFromFile("git", path, password)
}
