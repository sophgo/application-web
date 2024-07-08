package dist

import "embed"

//go:embed index.html
var IndexHtml embed.FS

//go:embed assets/*
var Assets embed.FS

//go:embed index.html
var IndexByte []byte

//go:embed admin.gif
var AdminGif embed.FS

//go:embed logo.png
var LogoPng embed.FS

//go:embed favicon.ico
var Favicon embed.FS
