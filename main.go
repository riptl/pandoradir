package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/valyala/fasthttp"
	"hash/crc32"
	"hash/crc64"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var crcTable = crc32.IEEETable
var crc64Table = crc64.MakeTable(crc64.ECMA)
var bin = binary.BigEndian

const maxTimestamp = int64(1257894000)

const (
	KindFile = byte(iota)
	KindDir
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <bind>\n",
			filepath.Base(os.Args[0]))
		return
	}
	bind := os.Args[1]

	fmt.Println("Serving at", bind)
	err := fasthttp.ListenAndServe(bind, baseHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func baseHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, ".bin")

	ctx.Response.Header.Set("Server", "pandoradir")

	// Base dir
	if bytes.Equal(ctx.Method(), []byte("GET")) && path == "" {
		serveRootListing(ctx)
		return
	}

	var parts []uint64

	for _, partHex := range strings.Split(path, "/") {
		switch len(partHex) {
		case 16:
			part, err := strconv.ParseUint(partHex, 16, 64)
			if err != nil {
				ctx.SetStatusCode(404)
				return
			}
			parts = append(parts, part)
		case 0:
			continue
		default:
			ctx.SetStatusCode(404)
			return
		}
	}

	for i, part := range parts {
		id := uint32(part)
		var idBytes [4]byte
		bin.PutUint32(idBytes[:], id)

		// Verify checksum
		crc := uint32(part >> 32)
		if crc32.Checksum(idBytes[:], crcTable) != crc {
			ctx.SetStatusCode(404)
			return
		}

		// Check if is "(dir)*(file|dir)"
		kind := byte(id) // & 0xFF
		if i < len(parts)-1 {
			if kind != KindDir {
				ctx.SetStatusCode(404)
				return
			}
		} else {
			switch kind {
			case KindFile:
				serveFile(ctx, parts[len(parts)-1])
			case KindDir:
				serveListing(ctx, parts)
			default:
				ctx.SetStatusCode(404)
			}
		}
	}
}
