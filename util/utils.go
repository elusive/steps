package util 

import (
    "io/fs"
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
      if e != nil { return e }
      if filepath.Ext(d.Name()) == ext {
         a = append(a, s)
      }
      return nil
   })
   return a
}


