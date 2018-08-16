package main

const templateStr = `
<html>
<head>
  <title>Index of {{.Path}}</title>
</head>
<body>
<h1>Index of {{.Path}}</h1>
<pre>
<a href=".">Name</a>                  <a href=".">Last modified</a>        <a href=".">Size</a>
<hr>
{{- if .ParentDir -}}
<a href="{{ .ParentDir }}">Parent Directory</a>
{{ end -}}
{{- range $key, $value := .Dirs -}}
<a href="{{ $value.Link }}">{{ $value.Link }}</a>     {{ $value.Time }}     -
{{ end -}}
{{- range $key, $value := .Files -}}
<a href="{{ $value.Link }}">{{ $value.Link }}</a>      {{ $value.Time }}     {{ $value.Size }}
{{ end -}}
<hr>
</pre>
</body>
</html>
`

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
