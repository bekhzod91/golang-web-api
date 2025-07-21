package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

func LoadEnv() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		panic(err)
	}

	envDir, err := findEnvDir(cwd)
	if err != nil {
		fmt.Println("Error finding .env directory:", err)
		panic(err)
	}

	err = godotenv.Load(filepath.Join(envDir, ".env"))
	if err != nil {
		panic(err)
	}
}

// Function to find the .env directory
func findEnvDir(start string) (string, error) {
	// Walk up the directory tree
	for {
		// Check for the .env file in the current directory
		if _, err := os.Stat(filepath.Join(start, ".env")); err == nil {
			return start, nil // Return the root directory
		}

		// Move up one directory
		parent := filepath.Dir(start)
		if parent == start { // We've reached the root
			return "", fmt.Errorf(".env not found in any parent directory")
		}
		start = parent
	}
}
