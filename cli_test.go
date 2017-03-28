package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("TESTING", "1")
	m.Run()
}

func TestCli(t *testing.T) {
	t.Run("os.Stdin reading error", func(t *testing.T) {
		tmpOut := os.Stdout
		tmpIn := os.Stdin
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		os.Stdout = r
		os.Stdin = w

		os.Stdin, _ = os.Open("/path/to/unknown/file")
		main()

		os.Stdin = tmpIn
		os.Stdout = tmpOut
	})

}
