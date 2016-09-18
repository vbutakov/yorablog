package main

import "regexp"

var (
	// IndexURLPattern - URL pattern for index page
	IndexURLPattern *regexp.Regexp

	// PostURLPattern - URL pattern for post page
	PostURLPattern *regexp.Regexp

	// EditURLPattern - URL pattern for edit page
	EditURLPattern *regexp.Regexp

	// CreateURLPattern - URL pattern for create page
	CreateURLPattern *regexp.Regexp

	TestPostURLPattern *regexp.Regexp
)

// InitURLPatterns initialize URL patterns
func InitURLPatterns() {
	IndexURLPattern = regexp.MustCompile(`^/([\d]*)/?$`)
	PostURLPattern = regexp.MustCompile(`^/post/([\d]+)/?$`)
	EditURLPattern = regexp.MustCompile(`^/edit/([\d]+)/?$`)
	CreateURLPattern = regexp.MustCompile(`^/create/?$`)

	TestPostURLPattern = regexp.MustCompile(`^/test_post/([\d]+)/?$`)
}
