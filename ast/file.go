package ast

import "os"
import "io/ioutil"

type File struct {
	Scope    *Scope
	FileName string
	FileInfo os.FileInfo
}

// NewFile stat the file to get the file info
func NewFile(filename string) (*File, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	return &File{FileName: filename, FileInfo: fi}, nil
}

func (f *File) ReadFile() ([]byte, error) {
	return ioutil.ReadFile(f.FileName)
}

func (f *File) String() string {
	return f.FileName
}
