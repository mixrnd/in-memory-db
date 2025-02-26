package testingh

import (
	"bufio"
	"os"
	"path/filepath"
)

type BaseDirSuite struct {
	ContextSuite
	BaseDir string
}

func (s *BaseDirSuite) SetupTest() {
	s.ContextSuite.SetupTest()
	if _, err := os.Stat(s.BaseDir); os.IsNotExist(err) {
		if err := os.MkdirAll(s.BaseDir, os.ModePerm); err != nil {
			s.NoError(err)
		}
	}
}

func (s *BaseDirSuite) TearDownTest() {
	s.ContextSuite.TearDownTest()
	err := os.RemoveAll(s.BaseDir)
	s.NoError(err)
}

func (s *BaseDirSuite) ReadFileToSlice(filePath string) []string {
	readFile, err := os.Open(filePath)
	s.NoError(err)
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	err = readFile.Close()
	s.NoError(err)

	return fileLines
}

func (s *BaseDirSuite) FileNamesInBaseDir() []string {
	files, err := os.ReadDir(s.BaseDir)
	s.NoError(err)

	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	return fileNames
}

func GetBaseDir(_ uintptr, filename string, _ int, _ bool) string {
	return filepath.Dir(filename) + "/test_files/"
}
