package main

import (
	"bytes"
	"fmt"
	"github.com/valyala/fasthttp"
	"hash/crc32"
	"hash/crc64"
	"strings"
	"time"
)

func genIDs(id uint64) (parts []uint64) {
	round := uint32(id) // & 0xFFFFFFFF
	parts = make([]uint64, 16)
	for i := range parts {
		var roundBin [4]byte

		// Round to byte array
		bin.PutUint32(roundBin[:], round)

		// Generate ID
		round = crc32.Checksum(roundBin[:], crcTable)

		// Set last byte of ID (kind)
		round &^= 0xFF
		if i >= len(parts)/2 {
			round |= 1
		}

		parts[i] = uint64(round)

		// Round to byte array
		bin.PutUint32(roundBin[:], round)

		// Generate ID
		round = crc32.Checksum(roundBin[:], crcTable)
		parts[i] |= uint64(round) << 32
	}
	return
}

func genLink(part uint64) string {
	return fmt.Sprintf("%016x", part)
}

func genFullPath(ids []uint64) string {
	var buf bytes.Buffer
	for _, id := range ids {
		buf.WriteString(fmt.Sprintf("/%016x", id))
	}

	switch byte(ids[len(ids)-1]) { // & 0xFF
	case KindFile:
		buf.WriteString(".bin")
	case KindDir:
		buf.WriteString("/")
	}

	return buf.String()
}

func serveRootListing(ctx *fasthttp.RequestCtx) {
	info := templateInfo{Path: "/"}
	genListingInfo(&info, 0)

	ctx.SetStatusCode(200)
	ctx.Response.Header.Set(
		"content-type", "text/html")

	writeListing(ctx, &info)
}

func serveListing(ctx *fasthttp.RequestCtx, ids []uint64) {
	if !bytes.Equal(ctx.Method(), []byte("GET")) {
		ctx.SetStatusCode(404)
		return
	}

	ctx.SetStatusCode(200)
	ctx.Response.Header.Set(
		"content-type", "text/html")

	info := templateInfo{Path: genFullPath(ids)}
	pathNoTrailingSlash := strings.TrimSuffix(info.Path, "/")
	lastSlash := strings.LastIndexByte(pathNoTrailingSlash, '/')
	info.ParentDir = info.Path[:lastSlash+1]
	genListingInfo(&info, ids[len(ids)-1])

	writeListing(ctx, &info)
}

func genListingInfo(info *templateInfo, lastID uint64) {
	ids := genIDs(lastID)
	files := ids[:len(ids)/2]
	dirs := ids[len(ids)/2:]

	info.Dirs = make([]templateDir, len(dirs))
	info.Files = make([]templateFile, len(files))

	for i, dirId := range dirs {
		// Name
		info.Dirs[i].Link = genLink(dirId) + "/"

		// Generate 8 bytes of random data
		var extraData [8]byte
		bin.PutUint64(extraData[:], dirId)
		bin.PutUint64(extraData[:],
			crc64.Checksum(extraData[:], crc64Table))

		// Time
		timeUnix := int64(bin.Uint64(extraData[0:8]))
		timeUnix %= maxTimestamp
		timestamp := time.Unix(timeUnix, 0)
		info.Dirs[i].Time = timestamp.Format("2006-01-02 15:04")
	}

	for i, fileId := range files {
		// Name
		info.Files[i].Link = genLink(fileId)

		// Generate 16 bytes of random data
		var extraData [16]byte
		bin.PutUint64(extraData[0:8], fileId)
		bin.PutUint64(extraData[0:8],
			crc64.Checksum(extraData[0:8], crc64Table))
		bin.PutUint64(extraData[8:16],
			crc64.Checksum(extraData[0:8], crc64Table))

		// Size
		info.Files[i].Size = bin.Uint32(extraData[0:4])

		// Time
		timeUnix := int64(bin.Uint64(extraData[8:16]))
		timeUnix %= maxTimestamp
		timestamp := time.Unix(timeUnix, 0)
		info.Files[i].Time = timestamp.Format("2006-01-02 15:04")
	}
}
