package main

import (
	"github.com/valyala/fasthttp"
	"hash/crc64"
	"net"
	"strconv"
	"time"
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

func serveFile(ctx *fasthttp.RequestCtx, fileId uint64) {
	switch string(ctx.Method()) {
	case "HEAD":
		ctx.SetStatusCode(200)
		ctx.Response.Header.Set(
			"Content-Length", strconv.FormatUint(uint64(fileSize(fileId)), 10))
		ctx.Response.Header.Set(
			"Accept-Ranges", "bytes")
		return
	case "GET":
		ctx.Hijack(func(c net.Conn) {
			time.Sleep(30 * time.Second)
		})
	}
}
