package main

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type reader struct{}

func (r reader) Read(p []byte) (int, error) {
	return 0, errors.New("Reading error")
}

func TestPipe(t *testing.T) {
	t.Run("Error from io.Reader", func(t *testing.T) {
		var buf bytes.Buffer
		r := bufio.NewReader(reader{})

		z := zapper{r, &buf}
		z.pipe()

		assert.Equal(t, "Error reading from pipe: Reading error\n", buf.String(), "should be equal")
	})

	t.Run("Bad formatted JSON", func(t *testing.T) {
		var buf bytes.Buffer
		r := bufio.NewReader(bytes.NewBuffer([]byte(`{"level":"info""ts":1490617664.5938616,"caller":"test/generator.go:23","msg":"New info message"}\n`)))

		z := zapper{r, &buf}
		z.pipe()

		assert.Equal(t, "Error: Value looks like object, but can't find closing '}' symbol\n", buf.String(), "should be equal")
	})

	t.Run("Bad timestamp", func(t *testing.T) {
		var buf bytes.Buffer
		r := bufio.NewReader(bytes.NewBuffer([]byte(`{"level":"info", "ts":"abcdefg","caller":"test/generator.go:23","msg":"New info message"}\n`)))

		z := zapper{r, &buf}
		z.pipe()

		assert.Equal(t, "Error: strconv.ParseFloat: parsing \"abcdefg\": invalid syntax\n", buf.String(), "should be equal")
	})

	t.Run("Good formatted JSON", func(t *testing.T) {
		var buf bytes.Buffer
		r := bufio.NewReader(bytes.NewBuffer([]byte(`{"level":"error","ts":1490617661.5927453,"caller":"test/generator.go:31","msg":"Generating new random stuff","name":"Pavel","surname":"Makarenko","backoff":1,"stacktrace":"go.uber.org/zap.Stack\n\t/Users/w1n2k/Work/Golang/src/go.uber.org/zap/field.go:209\ngo.uber.org/zap.(*Logger).check\n\t/Users/w1n2k/Work/Golang/src/go.uber.org/zap/logger.go:273\ngo.uber.org/zap.(*Logger).Check\n\t/Users/w1n2k/Work/Golang/src/go.uber.org/zap/logger.go:146\ngo.uber.org/zap.(*SugaredLogger).log\n\t/Users/w1n2k/Work/Golang/src/go.uber.org/zap/sugar.go:223\ngo.uber.org/zap.(*SugaredLogger).Errorw\n\t/Users/w1n2k/Work/Golang/src/go.uber.org/zap/sugar.go:181\nmain.main\n\t/Users/w1n2k/Work/Golang/src/github.com/m1ome/zapper/test/generator.go:31"}`)))

		z := zapper{r, &buf}
		z.pipe()

		expected := "2017-03-27T15:27:41Z [ERRO]: Generating new random stuff @ test/generator.go:31\n\tFields:\n\t\tbackoff:1\n\t\tname:Pavel\n\t\tsurname:Makarenko\n\tStacktrace:\n\t\tgo.uber.org/zap.Stack\n\t\t\\t/Users/w1n2k/Work/Golang/src/go.uber.org/zap/field.go:209\n\t\tgo.uber.org/zap.(*Logger).check\n\t\t\\t/Users/w1n2k/Work/Golang/src/go.uber.org/zap/logger.go:273\n\t\tgo.uber.org/zap.(*Logger).Check\n\t\t\\t/Users/w1n2k/Work/Golang/src/go.uber.org/zap/logger.go:146\n\t\tgo.uber.org/zap.(*SugaredLogger).log\n\t\t\\t/Users/w1n2k/Work/Golang/src/go.uber.org/zap/sugar.go:223\n\t\tgo.uber.org/zap.(*SugaredLogger).Errorw\n\t\t\\t/Users/w1n2k/Work/Golang/src/go.uber.org/zap/sugar.go:181\n\t\tmain.main\n\t\t\\t/Users/w1n2k/Work/Golang/src/github.com/m1ome/zapper/test/generator.go:31\n"
		actual := buf.String()

		assert.Equal(t, expected, actual, "Not equal")
	})

	t.Run("Check formatters levels", func(t *testing.T) {
		tests := []struct {
			preset   string
			expected string
		}{
			{
				preset:   `{"level":"info","ts":1490617658.5809216,"caller":"test/generator.go:23","msg":"New info message"}`,
				expected: "2017-03-27T15:27:38Z [INFO]: New info message @ test/generator.go:23\n",
			},
			{
				preset:   `{"level":"debug","ts":1490617658.5809216,"caller":"test/generator.go:23","msg":"New info message"}`,
				expected: "2017-03-27T15:27:38Z [DEBG]: New info message @ test/generator.go:23\n",
			},
			{
				preset:   `{"level":"warn","ts":1490617658.5809216,"caller":"test/generator.go:23","msg":"New info message"}`,
				expected: "2017-03-27T15:27:38Z [WARN]: New info message @ test/generator.go:23\n",
			},
			{
				preset:   `{"level":"error","ts":1490617658.5809216,"caller":"test/generator.go:23","msg":"New info message"}`,
				expected: "2017-03-27T15:27:38Z [ERRO]: New info message @ test/generator.go:23\n",
			},
			{
				preset:   `{"level":"dpanic","ts":1490617658.5809216,"caller":"test/generator.go:23","msg":"New info message"}`,
				expected: "2017-03-27T15:27:38Z [DPAN]: New info message @ test/generator.go:23\n",
			},
			{
				preset:   `{"level":"panic","ts":1490617658.5809216,"caller":"test/generator.go:23","msg":"New info message"}`,
				expected: "2017-03-27T15:27:38Z [PANC]: New info message @ test/generator.go:23\n",
			},
			{
				preset:   `{"level":"fatal","ts":1490617658.5809216,"caller":"test/generator.go:23","msg":"New info message"}`,
				expected: "2017-03-27T15:27:38Z [FATL]: New info message @ test/generator.go:23\n",
			},
		}

		for _, test := range tests {
			r := bufio.NewReader(bytes.NewBufferString(test.preset))
			buf := &bytes.Buffer{}

			z := zapper{r, buf}
			z.pipe()

			assert.Equal(t, test.expected, buf.String(), "Not equal")
		}

	})
}
