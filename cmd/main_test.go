package main

import (
	"context"
	"fmt"

	"github.com/raiyni/compose-ops/pkg/git"
)

func basic(repository string) {
	g := git.NewGitClient()
	opts := git.FetchOption{
		BaseOption: git.BaseOption{
			RepositoryUrl: repository,
		},
	}
	hash, err := g.LatestCommit(context.TODO(), opts)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Latest hash: %s \n", hash)
}

func ssh(repository, keyPath string) {
	g := git.NewGitClient()
	opts := git.FetchOption{
		BaseOption: git.BaseOption{
			RepositoryUrl: repository,
		},
		RepoAuth: git.RepoAuth{
			SshAuth: git.SshAuth{
				KeyPath:     keyPath,
				KeyPassword: "",
			},
		},
	}

	hash, err := g.LatestCommit(context.TODO(), opts)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Latest hash: %s \n", hash)
}
