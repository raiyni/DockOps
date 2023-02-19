package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/raiyni/dockops/v1/pkg/git"
)

func basic(repository string) {
	g := git.NewGitClient()
	opts := git.FetchOption{
		BaseOption: git.BaseOption{
			RepositoryUrl: repository,
		},
	}
	hash, err := g.LatestHashHttp(context.TODO(), opts)

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
		SshAuth: git.SshAuth{
			KeyPath:     keyPath,
			KeyPassword: "",
		},
	}

	hash, err := g.LatestHashSsh(context.TODO(), opts)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Latest hash: %s \n", hash)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ssh(os.Getenv("SSH_REPO"), os.Getenv("SSH_PATH"))
}
