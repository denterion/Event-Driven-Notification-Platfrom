package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// OverloadDotEnv searches upward from the current working directory for a
// `.env` file and overloads environment variables from it.
//
// This avoids "it works only when I run from repo root" issues on Windows.
func OverloadDotEnv() {
	if path, ok := findUp(".env"); ok {
		_ = godotenv.Overload(path)
		return
	}
	// Fallback: try default behavior (current directory).
	_ = godotenv.Overload()
}

func findUp(filename string) (string, bool) {
	dir, err := os.Getwd()
	if err != nil {
		return "", false
	}

	for {
		candidate := filepath.Join(dir, filename)
		if st, err := os.Stat(candidate); err == nil && !st.IsDir() {
			return candidate, true
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}
