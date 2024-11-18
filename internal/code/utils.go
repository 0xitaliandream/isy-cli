package code

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

func GenerateHash() (string, error) {
	bytes := make([]byte, 5) // 5 bytes for a 10-character hex string
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("could not generate random hash: %v", err)
	}
	return hex.EncodeToString(bytes), nil
}

func CopyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Determine the relative path to construct the target path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)

		// Skip the .isy directory
		if info.IsDir() && info.Name() == ".isy" {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		} else {
			// If the path is a file, copy its contents
			srcFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			dstFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			defer dstFile.Close()

			if _, err := io.Copy(dstFile, srcFile); err != nil {
				return err
			}
			return os.Chmod(dstFile.Name(), info.Mode())
		}
	})
}

func ComputeDirectoryHash(dir string) (string, error) {
	hasher := sha256.New()
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Salta la cartella .isy
		if info.IsDir() && info.Name() == ".isy" {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := io.Copy(hasher, file); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func ListBranchesSortedByDate(dir string) ([]string, error) {
	var branches []string
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("could not read branches directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			branches = append(branches, entry.Name())
		}
	}

	// Sort branches by modified time descending
	sort.Slice(branches, func(i, j int) bool {
		iInfo, _ := os.Stat(filepath.Join(dir, branches[i]))
		jInfo, _ := os.Stat(filepath.Join(dir, branches[j]))
		return iInfo.ModTime().After(jInfo.ModTime())
	})

	return branches, nil
}
