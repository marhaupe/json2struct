package cmd

import (
	"go/format"
	"io/ioutil"
	"path"
	"reflect"
	"strings"
	"testing"
)

const dirName = "__testdata"
const expectedSuffix = "_expected"

func TestFiles(t *testing.T) {
	inputFiles, err := listValidInputFiles()
	if err != nil {
		t.Fatal("Error reading input files", err)
	}

	inputExpected, err := setupInputExpected(inputFiles)
	if err != nil {
		t.Fatal("Error setting test cases up", err)
	}

	for input, expected := range inputExpected {
		actual := generate(input)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Test failed.\nActual: %v\nActualLen: %v\nExpected: %v\nExpectedLen: %v\n", actual, len(actual), expected, len(expected))
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

func setupInputExpected(inputFiles []string) (map[string]string, error) {
	inputExpected := make(map[string]string, len(inputFiles))
	for _, f := range inputFiles {
		input, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, err
		}
		expected, err := ioutil.ReadFile(f + expectedSuffix)
		if err != nil {
			return nil, err
		}
		formattedExpected, _ := format.Source(expected)
		inputExpected[string(input)] = string(formattedExpected)
	}
	return inputExpected, nil
}
