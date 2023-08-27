package util

import (
	"io/fs"
	"os"
	"path/filepath"
)

/**
 * function to find files by ext.
 * e.g.
 *  for _, s := range find("/root", ".md") {
 *     println(s)
 *  }
 */
func Find(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	return a
}

/**
 * Returns true if file exists, otherwise false.
 */
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
