package models

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("Downloads", testDownloads)
	t.Run("Files", testFiles)
	t.Run("Users", testUsers)
	t.Run("Chunks", testChunks)
	t.Run("Thumbnails", testThumbnails)
}

func TestDelete(t *testing.T) {
	t.Run("Downloads", testDownloadsDelete)
	t.Run("Files", testFilesDelete)
	t.Run("Users", testUsersDelete)
	t.Run("Chunks", testChunksDelete)
	t.Run("Thumbnails", testThumbnailsDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("Downloads", testDownloadsQueryDeleteAll)
	t.Run("Files", testFilesQueryDeleteAll)
	t.Run("Users", testUsersQueryDeleteAll)
	t.Run("Chunks", testChunksQueryDeleteAll)
	t.Run("Thumbnails", testThumbnailsQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("Downloads", testDownloadsSliceDeleteAll)
	t.Run("Files", testFilesSliceDeleteAll)
	t.Run("Users", testUsersSliceDeleteAll)
	t.Run("Chunks", testChunksSliceDeleteAll)
	t.Run("Thumbnails", testThumbnailsSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("Downloads", testDownloadsExists)
	t.Run("Files", testFilesExists)
	t.Run("Users", testUsersExists)
	t.Run("Chunks", testChunksExists)
	t.Run("Thumbnails", testThumbnailsExists)
}

func TestFind(t *testing.T) {
	t.Run("Downloads", testDownloadsFind)
	t.Run("Files", testFilesFind)
	t.Run("Users", testUsersFind)
	t.Run("Chunks", testChunksFind)
	t.Run("Thumbnails", testThumbnailsFind)
}

func TestBind(t *testing.T) {
	t.Run("Downloads", testDownloadsBind)
	t.Run("Files", testFilesBind)
	t.Run("Users", testUsersBind)
	t.Run("Chunks", testChunksBind)
	t.Run("Thumbnails", testThumbnailsBind)
}

func TestOne(t *testing.T) {
	t.Run("Downloads", testDownloadsOne)
	t.Run("Files", testFilesOne)
	t.Run("Users", testUsersOne)
	t.Run("Chunks", testChunksOne)
	t.Run("Thumbnails", testThumbnailsOne)
}

func TestAll(t *testing.T) {
	t.Run("Downloads", testDownloadsAll)
	t.Run("Files", testFilesAll)
	t.Run("Users", testUsersAll)
	t.Run("Chunks", testChunksAll)
	t.Run("Thumbnails", testThumbnailsAll)
}

func TestCount(t *testing.T) {
	t.Run("Downloads", testDownloadsCount)
	t.Run("Files", testFilesCount)
	t.Run("Users", testUsersCount)
	t.Run("Chunks", testChunksCount)
	t.Run("Thumbnails", testThumbnailsCount)
}

func TestHooks(t *testing.T) {
	t.Run("Downloads", testDownloadsHooks)
	t.Run("Files", testFilesHooks)
	t.Run("Users", testUsersHooks)
	t.Run("Chunks", testChunksHooks)
	t.Run("Thumbnails", testThumbnailsHooks)
}

func TestInsert(t *testing.T) {
	t.Run("Downloads", testDownloadsInsert)
	t.Run("Downloads", testDownloadsInsertWhitelist)
	t.Run("Files", testFilesInsert)
	t.Run("Files", testFilesInsertWhitelist)
	t.Run("Users", testUsersInsert)
	t.Run("Users", testUsersInsertWhitelist)
	t.Run("Chunks", testChunksInsert)
	t.Run("Chunks", testChunksInsertWhitelist)
	t.Run("Thumbnails", testThumbnailsInsert)
	t.Run("Thumbnails", testThumbnailsInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("DownloadToFileUsingFile", testDownloadToOneFileUsingFile)
	t.Run("ChunkToFileUsingFile", testChunkToOneFileUsingFile)
	t.Run("ThumbnailToFileUsingFile", testThumbnailToOneFileUsingFile)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("FileToDownloads", testFileToManyDownloads)
	t.Run("FileToChunks", testFileToManyChunks)
	t.Run("FileToThumbnails", testFileToManyThumbnails)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("DownloadToFileUsingFile", testDownloadToOneSetOpFileUsingFile)
	t.Run("ChunkToFileUsingFile", testChunkToOneSetOpFileUsingFile)
	t.Run("ThumbnailToFileUsingFile", testThumbnailToOneSetOpFileUsingFile)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {
	t.Run("DownloadToFileUsingFile", testDownloadToOneRemoveOpFileUsingFile)
}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {
	t.Run("FileToDownloads", testFileToManyAddOpDownloads)
	t.Run("FileToChunks", testFileToManyAddOpChunks)
	t.Run("FileToThumbnails", testFileToManyAddOpThumbnails)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {
	t.Run("FileToDownloads", testFileToManySetOpDownloads)
}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {
	t.Run("FileToDownloads", testFileToManyRemoveOpDownloads)
}

func TestReload(t *testing.T) {
	t.Run("Downloads", testDownloadsReload)
	t.Run("Files", testFilesReload)
	t.Run("Users", testUsersReload)
	t.Run("Chunks", testChunksReload)
	t.Run("Thumbnails", testThumbnailsReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("Downloads", testDownloadsReloadAll)
	t.Run("Files", testFilesReloadAll)
	t.Run("Users", testUsersReloadAll)
	t.Run("Chunks", testChunksReloadAll)
	t.Run("Thumbnails", testThumbnailsReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("Downloads", testDownloadsSelect)
	t.Run("Files", testFilesSelect)
	t.Run("Users", testUsersSelect)
	t.Run("Chunks", testChunksSelect)
	t.Run("Thumbnails", testThumbnailsSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("Downloads", testDownloadsUpdate)
	t.Run("Files", testFilesUpdate)
	t.Run("Users", testUsersUpdate)
	t.Run("Chunks", testChunksUpdate)
	t.Run("Thumbnails", testThumbnailsUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("Downloads", testDownloadsSliceUpdateAll)
	t.Run("Files", testFilesSliceUpdateAll)
	t.Run("Users", testUsersSliceUpdateAll)
	t.Run("Chunks", testChunksSliceUpdateAll)
	t.Run("Thumbnails", testThumbnailsSliceUpdateAll)
}

func TestUpsert(t *testing.T) {
	t.Run("Downloads", testDownloadsUpsert)
	t.Run("Files", testFilesUpsert)
	t.Run("Users", testUsersUpsert)
	t.Run("Chunks", testChunksUpsert)
	t.Run("Thumbnails", testThumbnailsUpsert)
}
