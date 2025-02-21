package wal

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/suite"
	"in-memory-db/internal/testingh"
)

type SegmentSuite struct {
	testingh.BaseDirSuite
}

func TestSegmentSuite(t *testing.T) {
	suite.Run(t, new(SegmentSuite))
}

func (s *SegmentSuite) SetupTest() {
	s.BaseDir = testingh.GetBaseDir(runtime.Caller(0))
	s.BaseDirSuite.SetupTest()
}

func (s *SegmentSuite) TestInitRead_BasedirWithData() {
	fileInfo := map[string]string{
		fmt.Sprintf(fileNameTemplate, 1): "1111111111",
		fmt.Sprintf(fileNameTemplate, 2): "2222222222",
		fmt.Sprintf(fileNameTemplate, 3): "3333333333",
		fmt.Sprintf(fileNameTemplate, 4): "4444444444",
	}
	for name, content := range fileInfo {
		f, err := os.OpenFile(s.BaseDir+name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		s.NoError(err)
		_, err = f.WriteString(content)
		s.NoError(err)
		err = f.Sync()
		s.NoError(err)
		err = f.Close()
		s.NoError(err)
	}

	segment := NewSegment(120, s.BaseDir)
	counter := 1
	err := segment.InitRead(func(data []byte) error {
		content := fileInfo[fmt.Sprintf(fileNameTemplate, counter)]
		s.Equal(string(data), content)
		counter++
		return nil
	})
	s.NoError(err)

	err = segment.Close()
	s.NoError(err)
}

func (s *SegmentSuite) TestInitReadAndWrite_NoFilesExisted() {
	segment := NewSegment(120, s.BaseDir)
	counter := 1
	err := segment.InitRead(func(data []byte) error {
		counter++
		return nil
	})
	s.NoError(err)
	s.Equal(1, counter)

	insertData := []byte("1111111111111111")
	err = segment.Write(insertData)
	s.NoError(err)
	err = segment.Close()
	s.NoError(err)

	f, err := os.OpenFile(s.BaseDir+fmt.Sprintf(fileNameTemplate, 1), os.O_RDONLY, 0644)
	s.NoError(err)
	buffer := make([]byte, len(insertData))
	_, err = f.Read(buffer)
	s.NoError(err)
	s.Equal(buffer, insertData)
	err = f.Close()
	s.NoError(err)
}

func (s *SegmentSuite) TestFileRotation() {
	segment := NewSegment(32, s.BaseDir)
	counter := 1
	err := segment.InitRead(func(data []byte) error {
		counter++
		return nil
	})
	s.NoError(err)
	s.Equal(1, counter)

	data15bytes := []byte("123456789abcdfg")
	err = segment.Write(data15bytes)
	s.NoError(err)

	err = segment.Write(data15bytes)
	s.NoError(err)

	err = segment.Write(data15bytes)
	s.NoError(err)
	err = segment.Close()
	s.NoError(err)

	fileNames := s.FileNamesInBaseDir()

	s.ElementsMatch(fileNames, []string{fmt.Sprintf(fileNameTemplate, 1), fmt.Sprintf(fileNameTemplate, 2)})
}

func (s *SegmentSuite) TestFileRotation_FirstDataMoreThenMaxSegmentSizeBytes() {
	segment := NewSegment(12, s.BaseDir)
	counter := 1
	err := segment.InitRead(func(data []byte) error {
		counter++
		return nil
	})
	s.NoError(err)
	s.Equal(1, counter)

	data15bytes := []byte("123456789abcdfg")
	err = segment.Write(data15bytes)
	s.NoError(err)
	err = segment.Close()
	s.NoError(err)

	fileNames := s.FileNamesInBaseDir()

	s.ElementsMatch(fileNames, []string{fmt.Sprintf(fileNameTemplate, 1)})

	fileContent := s.ReadFileToSlice(s.BaseDir + fmt.Sprintf(fileNameTemplate, 1))
	s.Equal(1, len(fileContent))
	s.Equal(string([]byte("123456789abcdfg")), fileContent[0])
}
