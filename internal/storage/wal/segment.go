package wal

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

// имя файла будет иметь вид data_123 - где 123 будет возрастающей последовательностью

const fileNameTemplate = "data_%d"

type Segment struct {
	MaxSegmentSizeBytes int64
	DataDirectory       string

	currentFile       *os.File
	currentFileNumber int
}

func NewSegment(maxSegmentSizeBytes int, dataDirectory string) *Segment {
	return &Segment{
		MaxSegmentSizeBytes: int64(maxSegmentSizeBytes),
		DataDirectory:       dataDirectory,
	}
}

func (s *Segment) Init(fileHandler func(data []byte) error) error {
	files, err := os.ReadDir(s.DataDirectory)
	if err != nil {
		return err
	}

	fileNums := make([]int, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fn := s.fileNumber(file.Name())
		if fn == 0 {
			continue
		}
		if fn > s.currentFileNumber {
			s.currentFileNumber = fn
		}
		fileNums = append(fileNums, fn)
	}
	// т.к. у нас нет гарантии порядка файлов, то на надо отсортировать
	slices.Sort(fileNums)

	fileFullPathTemplate := s.DataDirectory + fileNameTemplate
	for idx, fn := range fileNums {
		data, err := os.ReadFile(fmt.Sprintf(fileFullPathTemplate, fn))
		if err != nil {
			return err
		}

		if err := fileHandler(data); err != nil {
			return err
		}

		if idx == len(fileNums)-1 {
			s.currentFileNumber = fn
			if err = s.setAndOpenFile(); err != nil {
				return err
			}
		}
	}

	if len(fileNums) == 0 {
		s.currentFileNumber = 1
		if err = s.setAndOpenFile(); err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (s *Segment) Write(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	stat, err := s.currentFile.Stat()
	if err != nil {
		return err
	}

	if stat.Size()+int64(len(data)) >= s.MaxSegmentSizeBytes && stat.Size() > 0 {
		if err = s.currentFile.Close(); err != nil {
			return err
		}
		s.currentFileNumber++
		if err = s.setAndOpenFile(); err != nil {
			return err
		}
	}

	if _, err := s.currentFile.Write(data); err != nil {
		return err
	}

	if err := s.currentFile.Sync(); err != nil {
		return err
	}
	return nil
}

func (s *Segment) Close() error {
	if s.currentFile != nil {
		return s.currentFile.Close()
	}
	return nil
}

func (s *Segment) setAndOpenFile() error {
	fileFullPathTemplate := s.DataDirectory + fileNameTemplate
	fileName := fmt.Sprintf(fileFullPathTemplate, s.currentFileNumber)
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	s.currentFile = f
	return nil
}

func (s *Segment) fileNumber(fileName string) int {
	res := strings.Split(fileName, "_")
	if len(res) != 2 {
		return 0
	}

	number, _ := strconv.Atoi(res[1])
	return number
}
