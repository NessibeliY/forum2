package temp

import (
	"fmt"
	"path/filepath"
	"text/template"

	"forum/pkg/logger"
)

func NewTemplateCache(l *logger.Logger) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join("./templates/*.html"))
	if err != nil {
		l.Error("glob filepath .html", err)
		return nil, fmt.Errorf("glob filepath: ", err)
	}

	for _, page := range pages {
		name := filepath.Base(page)

		files := []string{
			"./templates/header.html",
			"./templates/posts.html",
			"./templates/user_info.html",
			"./templates/create_post.html",
			"./templates/post.html",
			page,
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			l.Error("parse html page", page, err)
			return nil, fmt.Errorf("parse html: ", err)
		}

		cache[name] = ts
	}

	l.Info("html templates cached")
	return cache, nil
}
