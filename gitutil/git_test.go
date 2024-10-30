package gitutil

import "testing"

func TestShort(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// zero values
		{
			"",
			"",
		},
		{
			"refs/",
			"refs/",
		},
		{
			"refs/tags/",
			"tags/",
		},
		{
			"refs/heads/",
			"heads/",
		},
		{
			"refs/remotes/",
			"remotes/",
		},
		// simple non-zero values
		{
			"main",
			"main",
		},
		{
			"refs/meta/config",
			"meta/config",
		},
		{
			"refs/tags/v0.0.1",
			"v0.0.1",
		},
		{
			"refs/heads/main",
			"main",
		},
		{
			"refs/remotes/origin/feature/AllowSlashes",
			"origin/feature/AllowSlashes",
		},
	}

	for _, test := range tests {
		if got := ShortName(test.input); got != test.want {
			t.Errorf("Short(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}
