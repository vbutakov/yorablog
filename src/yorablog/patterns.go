package main

import "regexp"

var (
	// IndexURLPattern - URL pattern for index page
	IndexURLPattern *regexp.Regexp

	// PostURLPattern - URL pattern for index page
	PostURLPattern *regexp.Regexp
)

// InitURLPatterns initialize URL patterns
func InitURLPatterns() {
	IndexURLPattern = regexp.MustCompile(`^/([\d]*)/?$`)
	PostURLPattern = regexp.MustCompile(`^/post/([\d]*)/?$`)
}
