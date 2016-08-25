package models

import "testing"
// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("GooseDBVersions", testGooseDBVersions)
	t.Run("Files", testFiles)
	t.Run("Users", testUsers)
	t.Run("Chunks", testChunks)
	t.Run("Thumbnails", testThumbnails)
}

func TestDelete(t *testing.T) {
	t.Run("GooseDBVersions", testGooseDBVersionsDelete)
	t.Run("Files", testFilesDelete)
	t.Run("Users", testUsersDelete)
	t.Run("Chunks", testChunksDelete)
	t.Run("Thumbnails", testThumbnailsDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("GooseDBVersions", testGooseDBVersionsQueryDeleteAll)
	t.Run("Files", testFilesQueryDeleteAll)
	t.Run("Users", testUsersQueryDeleteAll)
	t.Run("Chunks", testChunksQueryDeleteAll)
	t.Run("Thumbnails", testThumbnailsQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("GooseDBVersions", testGooseDBVersionsSliceDeleteAll)
	t.Run("Files", testFilesSliceDeleteAll)
	t.Run("Users", testUsersSliceDeleteAll)
	t.Run("Chunks", testChunksSliceDeleteAll)
	t.Run("Thumbnails", testThumbnailsSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("GooseDBVersions", testGooseDBVersionsExists)
	t.Run("Files", testFilesExists)
	t.Run("Users", testUsersExists)
	t.Run("Chunks", testChunksExists)
	t.Run("Thumbnails", testThumbnailsExists)
}

func TestFind(t *testing.T) {
	t.Run("GooseDBVersions", testGooseDBVersionsFind)
	t.Run("Files", testFilesFind)
	t.Run("Users", testUsersFind)
	t.Run("Chunks", testChunksFind)
	t.Run("Thumbnails", testThumbnailsFind)
}

func TestBind(t *testing.T) {
	t.Run("GooseDBVersions", testGooseDBVersionsBind)
	t.Run("Files", testFilesBind)
	t.Run("Users", testUsersBind)
	t.Run("Chunks", testChunksBind)
	t.Run("Thumbnails", testThumbnailsBind)
}

func TestOne(t *testing.T) {
	t.Run("GooseDBVersions", testGooseDBVersionsOne)
	t.Run("Files", testFilesOne)
	t.Run("Users", testUsersOne)
	t.Run("Chunks", testChunksOne)
	t.Run("Thumbnails", testThumbnailsOne)
}

func TestAll(t *testing.T) {
	t.Run("GooseDBVersions", testGooseDBVersionsAll)
	t.Run("Files", testFilesAll)
	t.Run("Users", testUsersAll)
	t.Run("Chunks", testChunksAll)
	t.Run("Thumbnails", testThumbnailsAll)
}

func TestCount(t *testing.T) {
	t.Run("GooseDBVersions", testGooseDBVersionsCount)
	t.Run("Files", testFilesCount)
	t.Run("Users", testUsersCount)
	t.Run("Chunks", testChunksCount)
	t.Run("Thumbnails", testThumbnailsCount)
}

func TestHelpers(t *testing.T) {
	t.Run("GooseDBVersions", testGooseDBVersionsInPrimaryKeyArgs)
	t.Run("GooseDBVersions", testGooseDBVersionsSliceInPrimaryKeyArgs)
	t.Run("Files", testFilesInPrimaryKeyArgs)
	t.Run("Files", testFilesSliceInPrimaryKeyArgs)
	t.Run("Users", testUsersInPrimaryKeyArgs)
	t.Run("Users", testUsersSliceInPrimaryKeyArgs)
	t.Run("Chunks", testChunksInPrimaryKeyArgs)
	t.Run("Chunks", testChunksSliceInPrimaryKeyArgs)
	t.Run("Thumbnails", testThumbnailsInPrimaryKeyArgs)
	t.Run("Thumbnails", testThumbnailsSliceInPrimaryKeyArgs)
}

func TestHooks(t *testing.T) {
	t.Run("GooseDBVersions", testGooseDBVersionsHooks)
	t.Run("Files", testFilesHooks)
	t.Run("Users", testUsersHooks)
	t.Run("Chunks", testChunksHooks)
	t.Run("Thumbnails", testThumbnailsHooks)
}

func TestInsert(t *testing.T) {
	t.Run("GooseDBVersions", testGooseDBVersionsInsert)
	t.Run("GooseDBVersions", testGooseDBVersionsInsertWhitelist)
	t.Run("Files", testFilesInsert)
	t.Run("Files", testFilesInsertWhitelist)
	t.Run("Users", testUsersInsert)
	t.Run("Users", testUsersInsertWhitelist)
	t.Run("Chunks", testChunksInsert)
	t.Run("Chunks", testChunksInsertWhitelist)
	t.Run("Thumbnails", testThumbnailsInsert)
	t.Run("Thumbnails", testThumbnailsInsertWhitelist)
}

// TestToMany tests cannot be run in parallel
// or postgres deadlocks will occur.
func TestToMany(t *testing.T) {
	t.Run("FileToManyChunks", testFileToManyChunks)
	t.Run("FileToManyThumbnails", testFileToManyThumbnails)
}

// TestToOne tests cannot be run in parallel
// or postgres deadlocks will occur.
func TestToOne(t *testing.T) {
	t.Run("ChunkToFile_File", testChunkToOneFile_File)
	t.Run("ThumbnailToFile_File", testThumbnailToOneFile_File)
}

func TestReload(t *testing.T) {
	t.Run("GooseDBVersions", testGooseDBVersionsReload)
	t.Run("Files", testFilesReload)
	t.Run("Users", testUsersReload)
	t.Run("Chunks", testChunksReload)
	t.Run("Thumbnails", testThumbnailsReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("GooseDBVersions", testGooseDBVersionsReloadAll)
	t.Run("Files", testFilesReloadAll)
	t.Run("Users", testUsersReloadAll)
	t.Run("Chunks", testChunksReloadAll)
	t.Run("Thumbnails", testThumbnailsReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("GooseDBVersions", testGooseDBVersionsSelect)
	t.Run("Files", testFilesSelect)
	t.Run("Users", testUsersSelect)
	t.Run("Chunks", testChunksSelect)
	t.Run("Thumbnails", testThumbnailsSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("GooseDBVersions", testGooseDBVersionsUpdate)
	t.Run("Files", testFilesUpdate)
	t.Run("Users", testUsersUpdate)
	t.Run("Chunks", testChunksUpdate)
	t.Run("Thumbnails", testThumbnailsUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("GooseDBVersions", testGooseDBVersionsSliceUpdateAll)
	t.Run("Files", testFilesSliceUpdateAll)
	t.Run("Users", testUsersSliceUpdateAll)
	t.Run("Chunks", testChunksSliceUpdateAll)
	t.Run("Thumbnails", testThumbnailsSliceUpdateAll)
}

func TestUpsert(t *testing.T) {
	t.Run("GooseDBVersions", testGooseDBVersionsUpsert)
	t.Run("Files", testFilesUpsert)
	t.Run("Users", testUsersUpsert)
	t.Run("Chunks", testChunksUpsert)
	t.Run("Thumbnails", testThumbnailsUpsert)
}

