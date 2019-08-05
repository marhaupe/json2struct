package cmd

import (
	"go/format"
	"io/ioutil"
	"path"
	"strings"
	"testing"
)

const dirName = "testdata"
const expectedSuffix = "_expected"

func TestFiles(t *testing.T) {
	inputFiles, err := listValidInputFiles()
	if err != nil {
		t.Fatal("Error reading input files", err)
	}

	for _, filename := range inputFiles {

		input := readFile(filename)
		expected := readFile(filename + expectedSuffix)
		actual := generate(input)

		formatExpectedBytes, _ := format.Source([]byte(expected))
		formatActualBytes, _ := format.Source([]byte(actual))

		expected = string(formatExpectedBytes)
		actual = string(formatActualBytes)

		if actual != expected {
			t.Errorf("Test failed. Filename: %v\nActual: %v\nActualLen: %v\nExpected: %v\nExpectedLen: %v\n",
				filename, actual, len(actual), expected, len(expected))
		}
	}
}

func listValidInputFiles() ([]string, error) {
	dirFiles, err := ioutil.ReadDir(dirName)
	if err != nil {
		return nil, err
	}

	var inputFiles []string
	inputFileHasExpectedFile := make(map[string]bool, 2)

	for _, f := range dirFiles {
		fileName := f.Name()
		isExpectedFile := strings.HasSuffix(fileName, expectedSuffix)
		if isExpectedFile {
			inputFileNameLength := strings.Index(fileName, expectedSuffix)
			inputFileName := fileName[:inputFileNameLength]
			inputFileHasExpectedFile[inputFileName] = true
		} else {
			inputFiles = append(inputFiles, fileName)
		}
	}

	var validInputFiles []string
	for _, f := range inputFiles {
		if inputFileHasExpectedFile[f] {
			validInputFiles = append(validInputFiles, path.Join(dirName, f))
		}
	}
	return validInputFiles, nil
}

func readFile(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("error reading file " + filename)
	}
	return string(content)
}
