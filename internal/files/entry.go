package files

import (
	"path/filepath"
)

// Entry is a single file tree entry.
type Entry struct {
	// Type of this entry.
	Type Type
	// Related path depending on the file type. Empty for directories,
	// target path for links, source files for regular files.
	RelatedPath string

	children map[string]*Entry
}

// IsDir returns true if the [Entry] is a directory.
func (e *Entry) IsDir() bool {
	return e.Type == TypeDirectory
}

// IsLink returns true if the [Entry] is a link.
func (e *Entry) IsLink() bool {
	return e.Type == TypeLink
}

// IsRegular returns true if the [Entry] is a regular file.
func (e *Entry) IsRegular() bool {
	return e.Type == TypeRegular
}

// AddFile adds a new regular file [Entry] children.
func (e *Entry) AddFile(name, relatedPath string) (*Entry, error) {
	entry := &Entry{
		Type:        TypeRegular,
		RelatedPath: relatedPath,
	}
	return e.AddEntry(name, entry)
}

// AddDirectory adds a new directory [Entry] children.
func (e *Entry) AddDirectory(name string) (*Entry, error) {
	entry := &Entry{
		Type: TypeDirectory,
	}
	return e.AddEntry(name, entry)
}

// AddLink adds a new link [Entry] children.
func (e *Entry) AddLink(name, relatedPath string) (*Entry, error) {
	entry := &Entry{
		Type:        TypeLink,
		RelatedPath: relatedPath,
	}
	return e.AddEntry(name, entry)
}

// AddEntry adds an arbitrary [Entry] as children. The caller is responsible
// for using only valid [Type]s and according fields.
func (e *Entry) AddEntry(name string, entry *Entry) (*Entry, error) {
	if !e.IsDir() {
		return nil, ErrEntryNotDir
	}
	if ee, exists := e.children[name]; exists {
		return ee, ErrEntryExists
	}
	if e.children == nil {
		e.children = make(map[string]*Entry)
	}
	e.children[name] = entry
	return entry, nil
}

// GetEntry getsan [Entry] for the given name. Return ErrEntryNotExists if it
// doesn't exist.
func (e *Entry) GetEntry(name string) (*Entry, error) {
	if !e.IsDir() {
		return nil, ErrEntryNotDir
	}
	entry, exists := e.children[name]
	if !exists {
		return nil, ErrEntryNotExists
	}
	return entry, nil
}

func (e *Entry) walk(base string, fn WalkFunc) error {
	for name, entry := range e.children {
		path := filepath.Join(base, name)
		if err := fn(path, entry); err != nil {
			return err
		}
		if entry.IsDir() {
			if err := entry.walk(path, fn); err != nil {
				return err
			}
		}
	}
	return nil
}
