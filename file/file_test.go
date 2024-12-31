package file

import "testing"

func TestIsExistE(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{
			input: "/tmp",
			want:  true,
		},
		{
			input: "/tmp/01",
			want:  false,
		},
	}

	for _, test := range tests {
		if got := IsExist(test.input); got != test.want {
			t.Errorf("IsExist(%s) = %t, want %t", test.input, got, test.want)
		}
	}
}
