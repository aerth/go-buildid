package main

import (
	"fmt"
	"github.com/aerth/go-buildid/elfnote"
	"golang.org/x/debug/elf"
	"os"
	"path/filepath"
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
