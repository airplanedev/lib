package runtime

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"regexp"
)

var (
	// commentRegex matches against the string produced by Comment() below.
	//
	// It is used to extract a slug from a comment in a script file.
	commentRegex = regexp.MustCompile(`Linked to (https://.*air.*/t/.*) \[do not edit this line\]`)
	// maxBytesToReadForSlug is the max bytes we should read in a file when looking for a task slug.
	maxBytesToReadForSlug int64 = 4096
)

// Comment generates a linking comment that is used
// to associate a script file with an Airplane task.
//
// This comment can be parsed out of a script file using Slug.
func Comment(r Interface, taskURL string) string {
	return r.FormatComment("Linked to " + taskURL + " [do not edit this line]")
}

// Slug returns the slug from the given file.
//
// Ok is true if the slug was found and isn't empty.
func Slug(filePath string) (slug string, ok bool) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	return slugFromReader(file)
}

func slugFromReader(reader io.Reader) (slug string, ok bool) {
	code := make([]byte, maxBytesToReadForSlug)
	_, err := io.ReadFull(reader, code)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return
	}

	result := commentRegex.FindSubmatch(code)
	if len(result) == 0 {
		return
	}

	u, err := url.Parse(string(result[1]))
	if err != nil {
		return
	}

	_, slug = path.Split(u.Path)
	ok = slug != ""
	return
}

// ErrNotLinked is an error that is raised when a path unexpectedly
// does not contain a slug. It can be used to explain to a user how
// they should link that file with a task.
type ErrNotLinked struct {
	Path string
}

func (e ErrNotLinked) Error() string {
	return fmt.Sprintf(
		"the file %s is not linked to a task",
		e.Path,
	)
}

func (e ErrNotLinked) ExplainError() string {
	return fmt.Sprintf(
		"You can link the file by running:\n  airplane init --slug <slug> %s",
		e.Path,
	)
}
