package testfiles

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

// TestFile creates a temporary file with the given filename, data and with the default permissions.
// The directory is automatically removed when the test and all its subtests complete.
// Each subsequent call to t.TempDir returns a unique directory.
func TestFile(t *testing.T, filename string, data []byte) string {
	t.Helper()
	dir, err := os.MkdirTemp(t.TempDir(), "")
	require.NoError(t, err)

	fp := filepath.Join(dir, filename)

	err = os.WriteFile(fp, data, 0o600)
	require.NoError(t, err)
	return fp
}

// TestFileWithCustomPermissions creates a temporary file with the given filename and permissions.
// The directory is automatically removed when the test and all its subtests complete.
// Each subsequent call to t.TempDir returns a unique directory.
func TestFileWithCustomPermissions(t *testing.T, filename string, data []byte, perms fs.FileMode) string {
	t.Helper()
	path := TestFile(t, filename, data)

	oldMask := unix.Umask(0o000)
	defer unix.Umask(oldMask)

	err := os.Chmod(path, perms)
	require.NoError(t, err)

	return path
}
