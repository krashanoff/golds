package server

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
)

func OpenBrowser(url string) error {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	return exec.Command(cmd, append(args, url)...).Start()
}

func writeAutoRefreshHTML(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<meta http-equiv="refresh" content="1.5; url=%[1]v">
	
		<title>Analyzing ...</title>
	</head>

	<body>
	Analyzing ... (<a href="%[1]v">refresh</a>)
	</body>
</html>`, r.URL.String())
}

type htmlPage struct {
	bytes.Buffer
	theme *Theme
	trans Translation
}

func NewHtmlPage(title, themeName string, inGenModeRootPages bool) *htmlPage {
	var page htmlPage
	page.Grow(4 * 1024 * 1024)

	fmt.Fprintf(&page, `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%s</title>
<link href="%s" rel="stylesheet">
<script src="%s"></script>
<body><div>
`,
		title,
		buildPageHref("css", themeName, inGenModeRootPages, "", nil),
		buildPageHref("jvs", "gold", inGenModeRootPages, "", nil),
	)

	return &page
}

func (page *htmlPage) Done() []byte {
	writePageGenerationInfo(page)

	page.WriteString(`</div></body></html>`)
	return append([]byte(nil), page.Bytes()...)
}
