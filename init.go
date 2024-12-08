package main

import (
	"github.com/gin-gonic/gin"
	"io/fs"
	"path/filepath"
	"strings"
)

func InitTemplate(r *gin.Engine) {
	var tmpFiles []string
	filepath.WalkDir("templates", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".html") {
			tmpFiles = append(tmpFiles, path)
		}
		return nil
	})
	r.LoadHTMLFiles(tmpFiles...)

	r.Static("/static", "./static")
}
