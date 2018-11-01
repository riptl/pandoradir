package main

import (
	"fmt"
	"hash/crc32"
	"hash/crc64"
	"log"
	"net/http"
	"path"
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

func serveRootListing(res http.ResponseWriter, dirPath string) {
	info := templateInfo{
		Path: dirPath,
	}
	genListingInfo(&info, 0)
	if dirPath != "/" {
		info.ParentDir = path.Dir(dirPath)
		if !strings.HasSuffix(info.ParentDir, "/") {
			info.ParentDir += "/"
		}
	}

	res.Header().Set("content-type", "text/html")
	res.WriteHeader(http.StatusOK)
	err := listingTemplate.Execute(res, info)
	if err != nil {
		log.Fatal(err)
	}
}

func serveListing(res http.ResponseWriter, req *http.Request, ids []uint64) {
	if req.Method != "GET" {
		notFound(res)
		return
	}

	info := templateInfo{
		Path: path.Clean(req.URL.Path),
	}
	pathNoTrailingSlash := strings.TrimSuffix(info.Path, "/")
	lastSlash := strings.LastIndexByte(pathNoTrailingSlash, '/')
	info.ParentDir = info.Path[:lastSlash+1]
	genListingInfo(&info, ids[len(ids)-1])

	res.WriteHeader(http.StatusOK)
	err := listingTemplate.Execute(res, info)
	if err != nil {
		log.Fatal(err)
	}
}

func genListingInfo(info *templateInfo, lastID uint64) {
	ids := genIDs(lastID)
	files := ids[:len(ids)/2]
	dirs := ids[len(ids)/2:]

	info.Files = make([]templateFile, len(files))

	if recursive != "" {
		info.Dirs = append(info.Dirs, templateDir{
			Link: "./" + recursive + "/",
			Show: recursive + "/",
			Pad:  recursivePad,
			Time: "1912-04-15 05:18",
		})
	}

	for _, dirId := range dirs {
		dir := templateDir{}

		// Name
		dir.Link = genLink(dirId) + "/"
		dir.Show = dir.Link
		dir.Pad = "    "

		// Generate 8 bytes of random data
		var extraData [8]byte
		bin.PutUint64(extraData[:], dirId)
		bin.PutUint64(extraData[:],
			crc64.Checksum(extraData[:], crc64Table))

		// Time
		timeUnix := int64(bin.Uint64(extraData[0:8]))
		timeUnix %= maxTimestamp
		timestamp := time.Unix(timeUnix, 0)
		dir.Time = timestamp.Format("2006-01-02 15:04")

		info.Dirs = append(info.Dirs, dir)
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
