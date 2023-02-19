package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/go-git/go-git/v5"
	. "github.com/go-git/go-git/v5/_examples"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/go-co-op/gocron"
)

func main() {
	cron()
}

func task() {
	fmt.Println("1")
}

func cron() {
	s := gocron.NewScheduler(time.UTC)
	s.TagsUnique()

	s.Every(1).Minute().Tag("foo").Do(task)

	s.StartBlocking()
}

func compose() {
	// Define the Docker Compose command and arguments.
	_, err := exec.Command("docker", "compose", "up", "-d").Output()
	if err != nil {
		fmt.Printf("Error running Docker Compose: %v\n", err)
		os.Exit(1)
	}
}

func clone() {
	CheckArgs("<url>", "<directory>", "<commit>")
	url, directory, commit := os.Args[1], os.Args[2], os.Args[3]

	// Clone the given repository to the given directory
	Info("git clone %s %s", url, directory)
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL: url,
	})

	CheckIfError(err)

	// ... retrieving the commit being pointed by HEAD
	Info("git show-ref --head HEAD")
	ref, err := r.Head()
	CheckIfError(err)
	fmt.Println(ref.Hash())

	w, err := r.Worktree()
	CheckIfError(err)

	// ... checking out to commit
	Info("git checkout %s", commit)
	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(commit),
	})
	CheckIfError(err)

	// ... retrieving the commit being pointed by HEAD, it shows that the
	// repository is pointing to the giving commit in detached mode
	Info("git show-ref --head HEAD")
	ref, err = r.Head()
	CheckIfError(err)
	fmt.Println(ref.Hash())
}
