package main

import "regexp"

var (
	// IndexURLPattern - URL pattern for index page
	IndexURLPattern *regexp.Regexp
)

// InitURLPatterns initialize URL patterns
func InitURLPatterns() {
	IndexURLPattern = regexp.MustCompile(`^/(?P<pageNum>[\d]*)/?$`)
}
