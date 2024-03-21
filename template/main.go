package main

import (
	"os"
	"text/template"
)

func main() {
	// 定义模板
	const tmpl = `
{{- range .Farms}}
	{{- with $farm := . }}
	{{- range $.Namespaces}}
		Namespace: {{.}} {{$farm}}
	{{- end }}
	{{- end }}
{{- end}}
`

	data := struct {
		Farms      []string
		Namespaces []string
	}{
		[]string{"dae", "adp"},
		[]string{"dae", "dae-pre", "dae-stage"},
	}
	// 创建模板并执行
	t := template.Must(template.New("tmpl").Parse(tmpl))
	err := t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
