package gitutil

import "fmt"

// RefRevParseRules are a set of rules to parse references into short names, or expand into a full reference.
// These are the same rules as used by git in shorten_unambiguous_ref and expand_ref.
// See: https://github.com/git/git/blob/e0aaa1b6532cfce93d87af9bc813fb2e7a7ce9d7/refs.c#L417
var RefRevParseRules = []string{
	"%s",
	"refs/%s",
	"refs/tags/%s",
	"refs/heads/%s",
	"refs/remotes/%s",
	"refs/remotes/%s/HEAD",
}

// ShortName returns the short name of a ReferenceName
func ShortName(name string) string {
	if name == "" {
		return ""
	}
	res := name
	for _, format := range RefRevParseRules[1:] {
		_, err := fmt.Sscanf(name, format, &res)
		if err == nil {
			continue
		}
	}

	return res
}
