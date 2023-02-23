package git

import (
	"context"
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
)

type gitClient struct {
	service config.Service
}

func NewGitClient(opts config.Service) *gitClient {
	return &gitClient{
		service: opts,
	}
}

func (c *gitClient) Clone(ctx context.Context, dst string) error {
	key, err := getAuth(c.service.AuthObj)
	if err != nil {
		return err
	}

	gitOptions := git.CloneOptions{
		URL:  c.service.Url,
		Auth: key,
	}

	if c.service.Ref != "" {
		gitOptions.ReferenceName = plumbing.ReferenceName(c.service.Ref)
	}

	_, err = git.PlainCloneContext(ctx, dst, false, &gitOptions)
	if err != nil {
		return errors.Wrap(err, "failed to clone git repository")
	}

	return nil
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
