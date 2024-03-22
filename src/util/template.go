package util

import (
	"bytes"
	"text/template"
)

// EncodeTemplate 模版渲染
func EncodeTemplate(name string, tpl string, data interface{}) (re string, err error) {
	tmpl, err := template.New(name).Parse(tpl)
	if err != nil {
		panic(err)
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, data)
	if err != nil {
		panic(err)
	}
	re = buffer.String()
	return
}
