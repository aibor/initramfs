package files

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fileEntry = Entry{Type: TypeRegular}
var dirEntry = Entry{Type: TypeDirectory}
var linkEntry = Entry{Type: TypeLink}

func TestIsRegular(t *testing.T) {
	assert.True(t, fileEntry.IsRegular())
	assert.False(t, dirEntry.IsRegular())
	assert.False(t, linkEntry.IsRegular())
}

func TestIsDir(t *testing.T) {
	assert.False(t, fileEntry.IsDir())
	assert.True(t, dirEntry.IsDir())
	assert.False(t, linkEntry.IsDir())
}

func TestIsLink(t *testing.T) {
	assert.False(t, fileEntry.IsLink())
	assert.False(t, dirEntry.IsLink())
	assert.True(t, linkEntry.IsLink())
}

func TestAddFile(t *testing.T) {
	p := dirEntry
	e, err := p.AddFile("file", "source")
	require.NoError(t, err)
	assert.Equal(t, TypeRegular, e.Type)
	assert.Equal(t, "source", e.RelatedPath)
	assert.Empty(t, e.children)
}

func TestAddDirectory(t *testing.T) {
	p := dirEntry
	e, err := p.AddDirectory("dir")
	require.NoError(t, err)
	assert.Equal(t, TypeDirectory, e.Type)
	assert.Equal(t, "", e.RelatedPath)
	assert.Empty(t, e.children)
}

func TestAddLink(t *testing.T) {
	p := dirEntry
	e, err := p.AddLink("link", "target")
	require.NoError(t, err)
	assert.Equal(t, TypeLink, e.Type)
	assert.Equal(t, "target", e.RelatedPath)
	assert.Empty(t, e.children)
}

func TestAddEntry(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		p := dirEntry
		n := Entry{}
		e, err := p.AddEntry("new", &n)
		require.NoError(t, err)
		assert.Equal(t, &n, e)
	})

	t.Run("exists", func(t *testing.T) {
		p := dirEntry
		n := Entry{}
		_, err := p.AddEntry("new", &n)
		require.NoError(t, err)
		e, err := p.AddEntry("new", &n)
		assert.ErrorIs(t, err, ErrEntryExists)
		assert.Equal(t, &n, e)
	})

	t.Run("not dir", func(t *testing.T) {
		p := fileEntry
		n := Entry{}
		_, err := p.AddEntry("new", &n)
		require.ErrorIs(t, err, ErrEntryNotDir)
	})
}

func TestGetEntry(t *testing.T) {
	entry := Entry{
		Type:        TypeRegular,
		RelatedPath: "source",
	}
	p := Entry{
		Type: TypeDirectory,
		children: map[string]*Entry{
			"file": &entry,
		},
	}

	t.Run("exists", func(t *testing.T) {
		e, err := p.GetEntry("file")
		require.NoError(t, err)
		assert.Equal(t, &entry, e)
	})

	t.Run("does not exist", func(t *testing.T) {
		_, err := p.GetEntry("404")
		assert.ErrorIs(t, err, ErrEntryNotExists)
	})

	t.Run("not dir", func(t *testing.T) {
		p := fileEntry
		_, err := p.GetEntry("file")
		require.ErrorIs(t, err, ErrEntryNotDir)
	})
}
