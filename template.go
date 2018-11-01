package main

import (
	"bytes"
	"github.com/valyala/fasthttp"
	"strconv"
)

// Ugly af :c
func writeListing(ctx *fasthttp.RequestCtx, info *templateInfo) (err error) {
	var b bytes.Buffer
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
		b.WriteString(value.Time)
		b.Write([]byte("     -\n"))
	}

	for _, value := range info.Files {
		b.Write([]byte(`<a href="`))
		b.WriteString(value.Link)
		b.Write([]byte(`">`))
		b.WriteString(value.Link)
		b.Write([]byte(`</a>      `))
		b.WriteString(value.Time)
		b.Write([]byte("     "))
		b.WriteString(strconv.FormatUint(uint64(value.Size), 10))
		b.WriteByte('\n')
	}

	b.WriteString("<hr>\n</pre>\n</body>\n</html>")

	_, err = ctx.Write(b.Bytes())
	return
}

type templateDir struct {
	Link string
	Time string
}

type templateFile struct {
	Link string
	Time string
	Size uint32
}

type templateInfo struct {
	Path      string
	ParentDir string
	Dirs      []templateDir
	Files     []templateFile
}
