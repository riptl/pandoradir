package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"hash/crc32"
	"hash/crc64"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

var crcTable = crc32.IEEETable
var crc64Table = crc64.MakeTable(crc64.ECMA)
var bin = binary.BigEndian
var listingTemplate *template.Template
var recursive, recursivePad string

const maxTimestamp = int64(1257894000)

const (
	KindFile = byte(iota)
	KindDir
)

func main() {
	var nologging bool
	flag.BoolVar(&nologging, "quiet", false, "Suppress logging output")
	flag.StringVar(&recursive, "recursive", "", "Add a folder to each listing that symlinks to the root directory")
	flag.Parse()

	recursive = path.Clean(recursive)
	recursive = strings.TrimSuffix(recursive, "/")
	recursive = url.PathEscape(recursive)
	if recursive != "" {
		for i := len(recursive); i < 20; i++ {
			recursivePad += " "
		}
	}
	if recursive == "." {
		recursive = ""
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <bind>\n",
			filepath.Base(os.Args[0]))
		return
	}
	bind := args[0]

	t := template.New("listing")
	var err error
	listingTemplate, err = t.Parse(templateStr)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", logHandler(baseHandler))

	log.Print("Serving at ", bind)
	log.Fatal(http.ListenAndServe(bind, nil))
}

func baseHandler(res http.ResponseWriter, req *http.Request) {
	p := req.URL.Path
	p = path.Clean(p)
	p = strings.TrimPrefix(p, "/")

	var parts []uint64

	for _, partHex := range strings.Split(p, "/") {
		if recursive != "" && partHex == recursive {
			parts = nil
			continue
		}

		switch len(partHex) {
		case 16:
			part, err := strconv.ParseUint(partHex, 16, 64)
			if err != nil {
				notFound(res)
				return
			}
			parts = append(parts, part)
		case 0:
			continue
		default:
			notFound(res)
			return
		}
	}

	// Base dir
	if len(parts) == 0 && req.Method == "GET" {
		serveRootListing(res, "/"+p)
		return
	}

	for i, part := range parts {
		id := uint32(part)
		var idBytes [4]byte
		bin.PutUint32(idBytes[:], id)

		// Verify checksum
		crc := uint32(part >> 32)
		if crc32.Checksum(idBytes[:], crcTable) != crc {
			notFound(res)
			return
		}

		// Check if is "(dir)*(file|dir)"
		kind := byte(id) // & 0xFF
		if i < len(parts)-1 {
			if kind != KindDir {
				notFound(res)
				return
			}
		} else {
			switch kind {
			case KindFile:
				serveFile(res, req, parts[len(parts)-1])
			case KindDir:
				serveListing(res, req, parts)
			default:
				notFound(res)
			}
		}
	}
}

func notFound(res http.ResponseWriter) {
	res.WriteHeader(404)
	res.Write([]byte("404 Not Found."))
}
