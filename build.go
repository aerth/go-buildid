package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aerth/go-buildid/elfnote"
	"golang.org/x/debug/dwarf"
	"golang.org/x/debug/elf"
)

func init() {

	if len(os.Args) <= 1 {
		usage()
		os.Exit(111)
	}

	for _, arg := range os.Args {
		switch arg {
		case "-h", "--help", "-help", "help", "usage":
			usage()
			os.Exit(111)
		default:
		}
	}

}

func usage() {
	println("go-buildid", "(https://github.com/aerth/go-buildid)")
	println("usage:")
	println("\t"+filepath.Base(os.Args[0]), "[elf]")

}

func main() {
	var filename = os.Args[1]
	if filename == "" {
		println("fatal: filename is required")
		os.Exit(1)
	}

	p, err := New(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "info: %s is not a go executable\n", filename)
		fmt.Fprintln(os.Stderr, "fatal:", err.Error(), filename)
		os.Exit(111)
	}

	id := p.Notes()
	if id == "" {
		fmt.Fprintf(os.Stderr, "info: %s is not a go executable\n", filename)
		os.Exit(111)
	}

	fmt.Println(id, filename)
}

func (p *Process) Notes() (note string) {
	s := p.efd.SectionByType(elf.SHT_NOTE)
	if s == nil {
		fmt.Fprintln(os.Stderr, "error searching for note (7) section")
	}

	notes, err := elfnote.ReadNotes(s, p.efd.ByteOrder)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading", s.Name, ":", err)
	}

	// space separated
	for _, n := range notes {

		if n.Type == elfnote.NT_GO_BUILD {
			data := string(n.Data)
			if data != "" {
				note += data
			}
		}

	}

	return note
}

/*
https://github.com/sitano/goelf
*/

type Process struct {
	path string

	efd *elf.File
	dwf *dwarf.Data
}

func New(path string) (*Process, error) {
	var err error

	p := &Process{}
	if p.efd, err = Open(path); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Process) DWARF() (*dwarf.Data, error) {
	var err error

	if p.dwf == nil {
		if p.dwf, err = p.efd.DWARF(); err != nil {
			return nil, err
		}
	}

	return p.dwf, err
}

func Open(path string) (*elf.File, error) {
	fd, err := os.OpenFile(path, 0, os.ModePerm)
	if err != nil {
		return nil, err
	}

	efd, err := elf.NewFile(fd)
	if err != nil {
		return nil, err
	}

	return efd, nil
}
