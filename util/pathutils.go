package util

import "strings"

func EnsureExtension(fname, ext string) string {
	if strings.HasSuffix(fname, "."+ext) {
		return fname
	} else {
		return fname + "." + ext
	}
}

func StripFirstPathElement(path string) string {
	parts := strings.SplitN(path, "/", 2)
	return parts[1]
}

func IsolateEndPathElement(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1] //Split's output will never have 0 elements, so there's no risk of indexing into -1. This is safe.
}
