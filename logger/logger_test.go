package logger

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func mockLogger(t *testing.T) (*bufio.Scanner, *os.File, *os.File) {

	reader, writer, err := os.Pipe()

	if err != nil {
		t.Error("couldn't get os Pipe:", err)
	}

	log.SetOutput(writer)

	return bufio.NewScanner(reader), reader, writer
}

func resetLogger(t *testing.T, reader *os.File, writer *os.File) {

	err := reader.Close()
	if err != nil {
		t.Error("error closing reader was ", err)
	}

	if err = writer.Close(); err != nil {
		t.Error("error closing writer was ", err)
	}
}

func TestLogger(t *testing.T) {

	log.SetFlags(0)

	scanner, reader, writer := mockLogger(t)
	defer resetLogger(t, reader, writer)

	Start(LVL_DEBUG)

	// Some test data
	tests := []sLogLine{
		{LVL_ERROR, "Something bad happened"},
		{LVL_WARN, "We need to warn you that abc != 123"},
		{LVL_APP, "For your information - hello"},
		{LVL_DEBUG, "This is really helpful for devs"},
	}

	for _, v := range tests {

		Log(v.lvl, v.msg)
		scanner.Scan()

		// Compare
		expected := fmt.Sprint(v.msg)
		actual := strings.Replace(strings.TrimSuffix(scanner.Text(), "\n"), "[logger_test.go:57] ", "", -1)

		if expected != actual {
			t.Error("Expected", expected, "Got", actual)
		}
	}
}
func TestVerbosity(t *testing.T) {

	log.SetFlags(0)

	Start(LVL_ERROR)

	// Some test data
	tests := []sLogLine{
		{LVL_ERROR, "Something bad happened"},
		{LVL_WARN, ""},
		{LVL_APP, ""},
		{LVL_DEBUG, ""},
	}

	for _, v := range tests {

		scanner, reader, writer := mockLogger(t)
		defer resetLogger(t, reader, writer)

		Log(v.lvl, v.msg)

		if len(v.msg) > 0 {
			scanner.Scan()
		}

		// Compare
		expected := v.msg
		actual := scanner.Text()

		if len(actual) > 0 {
			actual = strings.TrimLeft(strings.TrimSuffix(scanner.Text(), "\n"), " ")
		}

		if expected != actual {
			t.Error("Expected", expected, "Got", actual)
		}
	}
}

func TestMultiStartStop(t *testing.T) {

	log.SetOutput(ioutil.Discard)

	for i := 0; i <= 5; i++ {
		Start(LVL_ERROR)
		Log(LVL_ERROR, "Something bad happened")
		Stop()
	}
}

func BenchmarkLogger(b *testing.B) {

	log.SetOutput(ioutil.Discard)

	Start(LVL_ERROR)

	for i := 0; i < b.N; i++ {
		Log(LVL_ERROR, "Something bad happened")
	}
}
