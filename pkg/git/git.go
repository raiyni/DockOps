package git

import (
	"context"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pkg/errors"
)

type gitClient struct {
	auth RepoAuth
}

func NewGitClient() *gitClient {
	return &gitClient{}
}

func (c *gitClient) Clone(ctx context.Context, dst string, opt CloneOption) error {
	key, err := getAuth(opt.RepoAuth)
	if err != nil {
		return err
	}

	gitOptions := git.CloneOptions{
		URL:   opt.RepositoryUrl,
		Depth: opt.Depth,
		Auth:  key,
	}

	if opt.ReferenceName != "" {
		gitOptions.ReferenceName = plumbing.ReferenceName(opt.ReferenceName)
	}

	_, err = git.PlainCloneContext(ctx, dst, false, &gitOptions)
	if err != nil {
		return errors.Wrap(err, "failed to clone git repository")
	}

	return nil
}

func (c *gitClient) LatestCommit(ctx context.Context, opt FetchOption) (string, error) {
	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{opt.RepositoryUrl},
	})

	key, err := getAuth(opt.RepoAuth)
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

func getAuth(opt RepoAuth) (transport.AuthMethod, error) {
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
