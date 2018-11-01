package main

import (
	"bytes"
	"github.com/valyala/fasthttp"
	"strconv"
)

// Ugly af :c
func genListing(b *bytes.Buffer, info *templateInfo) {
	b.Write([]byte("<html>\n<head>\n  <title>Index of "))
	b.WriteString(info.Path)
	b.Write([]byte("</title>\n</head>\n<body>\n<h1>Index of "))
	b.WriteString(info.Path)
	b.Write([]byte("</h1>\n<pre>\n" +
		`<a href=".">Name</a>                  <a href=".">Last modified</a>        <a href=".">Size</a>` +
		"\n<hr>\n"))

	if info.ParentDir != "" {
		b.Write([]byte(`<a href="`))
		b.WriteString(info.ParentDir)
		b.Write([]byte(`">Parent Directory</a>` + "\n"))
	}

	for _, value := range info.Dirs {
		b.Write([]byte(`<a href="`))
		b.WriteString(value.Link)
		b.Write([]byte(`">`))
		b.WriteString(value.Link)
		b.Write([]byte(`</a>     `))
		fmtTime := formatTime(value.Time)
		b.Write(fmtTime[:])
		b.Write([]byte("     -\n"))
	}

	for _, value := range info.Files {
		b.Write([]byte(`<a href="`))
		b.WriteString(value.Link)
		b.Write([]byte(`">`))
		b.WriteString(value.Link)
		b.Write([]byte(`</a>      `))
		fmtTime := formatTime(value.Time)
		b.Write(fmtTime[:])
		b.Write([]byte("     "))
		b.WriteString(strconv.FormatUint(uint64(value.Size), 10))
		b.WriteByte('\n')
	}

	b.WriteString("<hr>\n</pre>\n</body>\n</html>")
}

func writeListing(ctx *fasthttp.RequestCtx, info *templateInfo) (err error) {
	var b bytes.Buffer
	genListing(&b, info)
	_, err = ctx.Write(b.Bytes())
	return
}

type templateDir struct {
	Link string
	Time uint64
}

type templateFile struct {
	Link string
	Time uint64
	Size uint32
}

type templateInfo struct {
	Path      string
	ParentDir string
	Dirs      []templateDir
	Files     []templateFile
}

func formatTime(unix uint64) (b [16]byte) {
	year := unix / 32140800
	month := (unix / 2678400) % 12
	day := (unix / 86400) % 28
	hour := (unix / 3600) % 24
	min := (unix / 60) % 60
	year += 1900
	month++
	day++

	// "2006-01-02 15:04"
	b[0] = 0x30 + byte(year/1000)
	b[1] = 0x30 + byte((year%1000)/100)
	b[2] = 0x30 + byte((year%100)/10)
	b[3] = 0x30 + byte(year%10)
	b[4] = '-'
	b[5] = 0x30 + byte((month%100)/10)
	b[6] = 0x30 + byte(month%10)
	b[7] = '-'
	b[8] = 0x30 + byte((day%100)/10)
	b[9] = 0x30 + byte(day%10)
	b[10] = ' '
	b[11] = 0x30 + byte((hour%100)/10)
	b[12] = 0x30 + byte(hour%10)
	b[13] = ':'
	b[14] = 0x30 + byte((min%100)/10)
	b[15] = 0x30 + byte(min%10)

	return
}
