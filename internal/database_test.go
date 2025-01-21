package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"in-memory-db/internal/compute"
	"in-memory-db/internal/storage"
	inmemory "in-memory-db/internal/storage/in-memory"
)

func TestDatabase_RunQuery_SetGetDel(t *testing.T) {
	db := createDataBaseForTest()

	r, err := db.RunQuery("SET key val")
	assert.NoError(t, err)
	assert.Equal(t, "[ok]", r)

	val, err := db.RunQuery("GET key")
	assert.NoError(t, err)
	assert.Equal(t, "val", val)

	_, err = db.RunQuery("GET notexisting")
	assert.ErrorIs(t, err, storage.ErrNotFound)

	_, err = db.RunQuery("DEL key")
	assert.NoError(t, err)

	_, err = db.RunQuery("GET key")
	assert.ErrorIs(t, err, storage.ErrNotFound)
}

func TestDatabase_RunQuery_UnknownCommand(t *testing.T) {
	db := createDataBaseForTest()

	_, err := db.RunQuery("CREATE key")
	assert.ErrorIs(t, err, compute.ErrUnknownCommand)
}

func createDataBaseForTest() *Database {
	e := inmemory.NewEngine()
	p := compute.NewParser()
	return NewDatabase(e, p, zap.NewNop())
}
