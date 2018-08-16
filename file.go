package main

import (
	"hash/crc64"
	"net/http"
	"strconv"
)

func fileSize(fileId uint64) uint32 {
	// Derive 8 bytes data
	var extraData [8]byte
	bin.PutUint64(extraData[0:8], fileId)
	bin.PutUint64(extraData[0:8],
		crc64.Checksum(extraData[0:8], crc64Table))

	// Size
	return bin.Uint32(extraData[0:4])
}

func serveFile(res http.ResponseWriter, req *http.Request, fileId uint64) {
	if req.Method == "HEAD" {
		res.Header().Set("Content-Length", strconv.FormatUint(uint64(fileSize(fileId)), 10))
		res.Header().Set("Accept-Ranges", "bytes")
		res.WriteHeader(http.StatusOK)
		return
	}

	res.WriteHeader(http.StatusForbidden)
}
