package internal

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"in-memory-db/internal/compute"
	"in-memory-db/internal/storage"
	inmemory "in-memory-db/internal/storage/in-memory"
	"in-memory-db/internal/storage/wal"
	"in-memory-db/internal/testingh"
)

type DatabaseSuite struct {
	testingh.BaseDirSuite
	walInst *wal.Wal
}

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(DatabaseSuite))
}

func (s *DatabaseSuite) SetupTest() {
	s.BaseDir = testingh.GetBaseDir(runtime.Caller(0))
	s.BaseDirSuite.SetupTest()
}

func (s *DatabaseSuite) createDataBaseForTest(segSize, batchSize int, tm time.Duration) *Database {
	e := inmemory.NewEngine()
	p := compute.NewParser()
	segment := wal.NewSegment(segSize, s.BaseDir)
	s.walInst = wal.NewWal(s.Ctx, batchSize, tm, segment, zap.NewNop())
	return NewDatabase(e, p, zap.NewNop(), s.walInst)
}

func (s *DatabaseSuite) TestDatabase_RunQuery_SetGetDel() {
	db := s.createDataBaseForTest(100, 100, 100*time.Millisecond)

	r, err := db.RunQuery("SET key val")
	s.NoError(err)
	s.Equal("[ok]", r)

	val, err := db.RunQuery("GET key")
	s.NoError(err)
	s.Equal("val", val)

	_, err = db.RunQuery("GET notexisting")
	s.ErrorIs(err, storage.ErrNotFound)

	_, err = db.RunQuery("DEL key")
	s.NoError(err)

	_, err = db.RunQuery("GET key")
	s.ErrorIs(err, storage.ErrNotFound)
}

func (s *DatabaseSuite) TestDatabase_RunQuery_UnknownCommand() {
	db := s.createDataBaseForTest(100, 100, 100*time.Millisecond)

	_, err := db.RunQuery("CREATE key")
	s.ErrorIs(err, compute.ErrUnknownCommand)
}

func (s *DatabaseSuite) TestDatabase_RunQuery_WriteToWalByTimeout() {
	db := s.createDataBaseForTest(4096, 4096, 100*time.Millisecond)
	db.Init()

	const queryNumber = 10
	for i := 0; i < queryNumber; i++ {
		r, err := db.RunQuery(fmt.Sprintf("SET key%d val", i))
		s.NoError(err)
		s.Equal("[ok]", r)
	}
	time.Sleep(200 * time.Millisecond)

	fileNames := s.FileNamesInBaseDir()
	s.ElementsMatch([]string{"data_1"}, fileNames)

	fileContent := s.ReadFileToSlice(s.BaseDir + "data_1")
	s.Equal(queryNumber, len(fileContent))
	for i := 0; i < len(fileContent); i++ {
		s.Equal(fmt.Sprintf("SET key%d val", i), fileContent[i])
	}
}

func (s *DatabaseSuite) TestDatabase_RunQuery_WriteToWalByBatchSize() {
	db := s.createDataBaseForTest(4096, 130, 1000*time.Millisecond)
	db.Init()

	const queryNumber = 10
	for i := 0; i < queryNumber; i++ {
		r, err := db.RunQuery(fmt.Sprintf("SET key%d val", i))
		s.NoError(err)
		s.Equal("[ok]", r)
	}
	s.CtxCancelFunc()
	s.walInst.WaitWrite()
	//time.Sleep(10 * time.Millisecond)
	fileNames := s.FileNamesInBaseDir()

	s.ElementsMatch([]string{"data_1"}, fileNames)

	fileContent := s.ReadFileToSlice(s.BaseDir + "data_1")
	s.Equal(queryNumber, len(fileContent))
	for i := 0; i < len(fileContent); i++ {
		s.Equal(fmt.Sprintf("SET key%d val", i), fileContent[i])
	}
}

func (s *DatabaseSuite) TestDatabase_RunQuery_WriteToWalTwoFiles() {
	db := s.createDataBaseForTest(104, 90, 500*time.Millisecond)
	db.Init()

	const queryNumber = 10
	for i := 0; i < queryNumber; i++ {
		r, err := db.RunQuery(fmt.Sprintf("SET key%d val", i))
		s.NoError(err)
		s.Equal("[ok]", r)
	}

	s.CtxCancelFunc()
	s.walInst.WaitWrite()

	files, err := os.ReadDir(s.BaseDir)
	s.NoError(err)

	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	s.ElementsMatch([]string{"data_1", "data_2"}, fileNames)

	fileContent := s.ReadFileToSlice(s.BaseDir + "data_1")
	fileContent = append(fileContent, s.ReadFileToSlice(s.BaseDir+"data_2")...)
	s.Equal(queryNumber, len(fileContent))
	for i := 0; i < len(fileContent); i++ {
		s.Equal(fmt.Sprintf("SET key%d val", i), fileContent[i])
	}
}

func (s *DatabaseSuite) TestDatabase_RunQuery_ReadWalInit() {
	fileInfo := map[string]string{
		"data_1": "SET key1 111\nSET key2 222\nSET key3 333\n",
		"data_2": "DEL key1\nSET key3 444\n",
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
	db := s.createDataBaseForTest(104, 104, 500*time.Millisecond)
	db.Init()

	_, err := db.RunQuery("GET key1")
	s.ErrorIs(err, storage.ErrNotFound)

	val, err := db.RunQuery("GET key2")
	s.Equal("222", val)

	val, err = db.RunQuery("GET key3")
	s.Equal("444", val)

	s.CtxCancelFunc()
	s.walInst.WaitWrite()
}
