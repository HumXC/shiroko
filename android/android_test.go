package android_test

import (
	"bytes"
	"testing"

	"github.com/HumXC/shiroko/android"
)

func TestTrimEnd(t *testing.T) {
	testsString := []string{
		"Hello\nWorld\n",
		"Hello\rWorld\r",
		"Hello\nWorld\r",
		"Hello\rWorld\n",
		"Hello\r",
		"Hello\n",
	}
	wantsString := []string{
		"Hello\nWorld",
		"Hello\rWorld",
		"Hello\nWorld",
		"Hello\rWorld",
		"Hello",
		"Hello",
	}
	for i, v := range testsString {
		want := wantsString[i]
		got := android.TrimEnd(v)
		if got != want {
			t.Errorf("TrimEnd(%s) = %s, want %s", v, got, want)
		}
	}
	testsByte := make([][]byte, len(testsString))
	wantsByte := make([][]byte, len(wantsString))
	for i, v := range testsString {
		testsByte[i] = []byte(v)
		wantsByte[i] = []byte(wantsString[i])
	}
	for i, v := range testsByte {
		want := wantsByte[i]
		got := android.TrimEnd(v)

		if !bytes.Equal(got, want) {
			t.Errorf("TrimEnd(%s) = %s, want %s", v, got, want)
		}
	}
}
