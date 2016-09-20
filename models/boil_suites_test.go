package models

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("Users", testUsers)
	t.Run("Files", testFiles)
	t.Run("Chunks", testChunks)
	t.Run("Thumbnails", testThumbnails)
}

func TestDelete(t *testing.T) {
	t.Run("Users", testUsersDelete)
	t.Run("Files", testFilesDelete)
	t.Run("Chunks", testChunksDelete)
	t.Run("Thumbnails", testThumbnailsDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("Users", testUsersQueryDeleteAll)
	t.Run("Files", testFilesQueryDeleteAll)
	t.Run("Chunks", testChunksQueryDeleteAll)
	t.Run("Thumbnails", testThumbnailsQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("Users", testUsersSliceDeleteAll)
	t.Run("Files", testFilesSliceDeleteAll)
	t.Run("Chunks", testChunksSliceDeleteAll)
	t.Run("Thumbnails", testThumbnailsSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("Users", testUsersExists)
	t.Run("Files", testFilesExists)
	t.Run("Chunks", testChunksExists)
	t.Run("Thumbnails", testThumbnailsExists)
}

func TestFind(t *testing.T) {
	t.Run("Users", testUsersFind)
	t.Run("Files", testFilesFind)
	t.Run("Chunks", testChunksFind)
	t.Run("Thumbnails", testThumbnailsFind)
}

func TestBind(t *testing.T) {
	t.Run("Users", testUsersBind)
	t.Run("Files", testFilesBind)
	t.Run("Chunks", testChunksBind)
	t.Run("Thumbnails", testThumbnailsBind)
}

func TestOne(t *testing.T) {
	t.Run("Users", testUsersOne)
	t.Run("Files", testFilesOne)
	t.Run("Chunks", testChunksOne)
	t.Run("Thumbnails", testThumbnailsOne)
}

func TestAll(t *testing.T) {
	t.Run("Users", testUsersAll)
	t.Run("Files", testFilesAll)
	t.Run("Chunks", testChunksAll)
	t.Run("Thumbnails", testThumbnailsAll)
}

func TestCount(t *testing.T) {
	t.Run("Users", testUsersCount)
	t.Run("Files", testFilesCount)
	t.Run("Chunks", testChunksCount)
	t.Run("Thumbnails", testThumbnailsCount)
}

func TestHooks(t *testing.T) {
	t.Run("Users", testUsersHooks)
	t.Run("Files", testFilesHooks)
	t.Run("Chunks", testChunksHooks)
	t.Run("Thumbnails", testThumbnailsHooks)
}

func TestInsert(t *testing.T) {
	t.Run("Users", testUsersInsert)
	t.Run("Users", testUsersInsertWhitelist)
	t.Run("Files", testFilesInsert)
	t.Run("Files", testFilesInsertWhitelist)
	t.Run("Chunks", testChunksInsert)
	t.Run("Chunks", testChunksInsertWhitelist)
	t.Run("Thumbnails", testThumbnailsInsert)
	t.Run("Thumbnails", testThumbnailsInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("ChunkToFile_File", testChunkToOneFile_File)
	t.Run("ThumbnailToFile_File", testThumbnailToOneFile_File)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("FileToChunks", testFileToManyChunks)
	t.Run("FileToThumbnails", testFileToManyThumbnails)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("ChunkToFile_File", testChunkToOneSetOpFile_File)
	t.Run("ThumbnailToFile_File", testThumbnailToOneSetOpFile_File)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {
	t.Run("FileToChunks", testFileToManyAddOpChunks)
	t.Run("FileToThumbnails", testFileToManyAddOpThumbnails)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {}

func TestReload(t *testing.T) {
	t.Run("Users", testUsersReload)
	t.Run("Files", testFilesReload)
	t.Run("Chunks", testChunksReload)
	t.Run("Thumbnails", testThumbnailsReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("Users", testUsersReloadAll)
	t.Run("Files", testFilesReloadAll)
	t.Run("Chunks", testChunksReloadAll)
	t.Run("Thumbnails", testThumbnailsReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("Users", testUsersSelect)
	t.Run("Files", testFilesSelect)
	t.Run("Chunks", testChunksSelect)
	t.Run("Thumbnails", testThumbnailsSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("Users", testUsersUpdate)
	t.Run("Files", testFilesUpdate)
	t.Run("Chunks", testChunksUpdate)
	t.Run("Thumbnails", testThumbnailsUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("Users", testUsersSliceUpdateAll)
	t.Run("Files", testFilesSliceUpdateAll)
	t.Run("Chunks", testChunksSliceUpdateAll)
	t.Run("Thumbnails", testThumbnailsSliceUpdateAll)
}

func TestUpsert(t *testing.T) {
	t.Run("Users", testUsersUpsert)
	t.Run("Files", testFilesUpsert)
	t.Run("Chunks", testChunksUpsert)
	t.Run("Thumbnails", testThumbnailsUpsert)
}

