package docker

import (
	"fmt"
	"os"
	"os/exec"
)

func Compose() {
	// Define the Docker Compose command and arguments.
	_, err := exec.Command("docker", "compose", "up", "-d").Output()
	if err != nil {
		fmt.Printf("Error running Docker Compose: %v\n", err)
		os.Exit(1)
	}
}
