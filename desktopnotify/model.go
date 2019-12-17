package desktopnotify

import (
	"bytes"
	"html/template"
)

type notifyService struct {
	tpl    *template.Template
	buffer bytes.Buffer
}
