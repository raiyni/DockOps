package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/raiyni/compose-ops/pkg/config"
)

func init() {
}

func main() {
	fmt.Println(os.Getenv("ABC"))
	c := config.Store.String("auth.github.username")
	fmt.Println(c)
}
