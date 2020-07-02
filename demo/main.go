package main

import (
	"fmt"
	"os"

	"github.com/haroldadmin/pathfix"
)

func main() {
	fmt.Printf("Before fixing: %s\n\n", os.Getenv("PATH"))
	pathfix.Fix()
	fmt.Printf("After fixing: %s\n", os.Getenv("PATH"))
}
