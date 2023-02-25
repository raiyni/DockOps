package git

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pkg/errors"
	"github.com/raiyni/compose-ops/pkg/config"
	"github.com/rs/zerolog/log"
)

type gitClient struct {
	service config.Service
	data    string
}

func NewGitClient(opts config.Service, data string) *gitClient {
	return &gitClient{
		service: opts,
		data:    data,
	}
}

func (c *gitClient) PullMostRecent(ctx context.Context, savedHash string) (string, error) {
	if _, err := os.Stat(c.getPath()); errors.Is(err, os.ErrNotExist) {
		return c.Clone(ctx)
	}

	return c.Pull(ctx)
}

func (c *gitClient) getPath() string {
	path := fmt.Sprintf("%s/%s", c.data, c.service.Name)
	if c.service.Path != "" {
		path = c.service.Path
	}

	return path
}

func (c *gitClient) Pull(ctx context.Context) (string, error) {
	log.Debug().Msgf("attempting to pull: %s", c.service.Name)
	key, err := getAuth(c.service.AuthObj)
	if err != nil {
		return "", err
	}

	gitOptions := &git.PullOptions{
		Auth:              key,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		RemoteName:        "origin",
	}

	if c.service.Ref != "" {
		gitOptions.ReferenceName = plumbing.ReferenceName(c.service.Ref)
	}

	path := c.getPath()

	r, err := git.PlainOpen(path)
	if err != nil {
		return "", err
	}

	w, err := r.Worktree()
	if err != nil {
		return "", err
	}

	err = w.Pull(gitOptions)
	if err != nil && err.Error() != "already up-to-date" {
		return "", err
	}

	return fetchCommit(r)
}

func (c *gitClient) Clone(ctx context.Context) (string, error) {
	log.Debug().Msgf("attempting to clone: %s", c.service.Name)

	key, err := getAuth(c.service.AuthObj)
	if err != nil {
		return "", err
	}

	gitOptions := git.CloneOptions{
		URL:               c.service.Url,
		Auth:              key,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}

	if c.service.Ref != "" {
		gitOptions.ReferenceName = plumbing.ReferenceName(c.service.Ref)
	}

	path := c.getPath()

	r, err := git.PlainCloneContext(ctx, path, false, &gitOptions)
	if err != nil {
		return "", err
	}

	return fetchCommit(r)
}

func (c *gitClient) LatestCommit(ctx context.Context) (string, error) {
	remote := git.NewRemote(memory.NewStorage(), &gitconfig.RemoteConfig{
		Name: "origin",
		URLs: []string{c.service.Url},
	})

	key, err := getAuth(c.service.AuthObj)
	if err != nil {
		return "", err
	}

	listOptions := &git.ListOptions{
		Auth: key,
	}

	refs, err := remote.List(listOptions)
	if err != nil {
		return "", errors.Wrap(err, "failed to list repository refs")
	}

	referenceName := c.service.Ref
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

	return "", errors.Errorf("could not find ref %q in the repository", c.service.Ref)
}

func getAuth(opt config.Auth) (transport.AuthMethod, error) {
	if opt.KeyPath != "" {
		return getSshAuth(opt.KeyPath, opt.KeyPassword)
	}

	return getBasicAuth(opt.Username, opt.Password)
}

func getBasicAuth(username, password string) (*githttp.BasicAuth, error) {
	if password != "" {
		if username == "" {
			username = "token"
		}

		return &githttp.BasicAuth{
			Username: username,
			Password: password,
		}, nil
	}

	return nil, nil
}

func getSshAuth(path, password string) (*ssh.PublicKeys, error) {
	return ssh.NewPublicKeysFromFile("git", path, password)
}

func fetchCommit(r *git.Repository) (string, error) {
	ref, err := r.Head()
	if err != nil {
		return "", err
	}

	return ref.Hash().String(), nil
}
