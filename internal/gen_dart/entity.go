package gen_dart

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/inclee/gokit/container/maps"
	"github.com/inclee/protogen/internal/entity"
	"github.com/inclee/protogen/internal/parse"
	"github.com/inclee/protogen/internal/utool"
)

// 生成代码的模板
const entiryTemplate = `

/*GENERATE BY protogen(https://github.com/inclee/protogen); PLEASE DON'T EDIT IT */

{{range .Imports}}
import './{{.}}';
{{- end }}
{{range .Messages}}
{{- $Method := .Method }}
{{- $Parent := .Parent }}
{{- if .Parent }}
	class {{.Name}} extends {{.Parent.Name}}{
{{- else }}
	class {{.Name}} {
{{- end }}
	{{- range .Fields }} 
		{{- $Validate := "" }}
		{{- if .Validate }}
			{{- $Validate = printf "validate:\"%s\"" .Validate }}
		{{- end }}
		{{- if .Type }}
			final {{ .Type | typeMap }} {{ .Name | lowerFirst }};
		{{- else }}
			{{ .Name | lowerFirst }} 
		{{- end }}
	{{- end }}
	{{- $len := len .Fields }}
	{{ .Name }}(
		{{- if $Parent }}
	 		{{- if gt $len 0 }}
			{
			{{- end}}
			{{- $flen := len $Parent.Fields }}
	 		{{- if eq $len 0 }}
	 			{{- if gt $flen 0 }}
					{
				{{- end}}
			{{- end}}
			{{- range .Fields }} 
				required this.{{ .Name | lowerFirst }},
			{{- end }}
			{{- range $pindex, $pfield := $Parent.Fields }} 
				required super.{{ $pfield.Name | lowerFirst }} {{ if lt $pindex (sub $flen 1) }},{{ end }}
			{{- end}}
			{{- if gt $len 0 }}
			}	
			{{- end}}
	 		{{- if eq $len 0 }}
	 			{{- if gt $flen 0 }}
					}	
				{{- end}}
			{{- end}}
			);
		{{- else}}
	 		{{- if gt $len 0 }}
			{	
				{{- range $index, $field := .Fields }} 
					required this.{{ $field.Name | lowerFirst }}{{ if lt $index (sub $len 1) }},{{ end }}
				{{- end }}
			}
			{{- end}}
			);
		{{- end}}

	factory {{.Name}}.fromJson(Map<String, dynamic> json) =>
      _${{ .Name }}FromJson(json);
	{{- if $Parent }}
	@override
	{{- end}}
  	Map<String, dynamic> toJson() => _${{ .Name }}ToJson(this);
}

{{ .Name }} _${{.Name}}FromJson(Map<String, dynamic> json) =>{{ .Name }}(
	 {{- range .Fields }} 
		{{ .Name | lowerFirst }}:json['{{ .Name | camelToSnakeCase }}'],
	 {{- end }}
	 {{- if $Parent }}
	 	{{- range $Parent.Fields }} 
			{{ .Name | lowerFirst }}:json['{{ .Name | camelToSnakeCase }}'],
	 	{{- end }}
	{{- end }}
);

Map<String, dynamic> _${{.Name}}ToJson({{ .Name }} instance) =><String, dynamic>{
	{{- range .Fields }} 
      	'{{ .Name | camelToSnakeCase }}': instance.{{ .Name | lowerFirst }},
	{{- end }}
	{{- if $Parent }}
		{{- range $Parent.Fields }} 
 	    	'{{ .Name | camelToSnakeCase }}': instance.{{ .Name | lowerFirst }},
		{{- end }}
	{{- end }}
};
	
{{end}}

`

func parseEntity(workdir string, filepath string) (string, error) {
	compnent, err := parse.Parsefile(workdir, filepath)
	if err != nil {
		return "", err
	}
	//
	// 准备模板数据
	data := struct {
		Package  string
		Imports  []string
		Messages []*entity.Message
	}{
		Package:  "",
		Imports:  compnent.Imports,
		Messages: maps.ValuesWith(compnent.Messages, func(msg *entity.Message) string { return msg.Name }, true),
	}

	// 创建模板并解析
	tmpl := template.Must(template.New("code").Funcs(template.FuncMap{"typeMap": typeMap, "sub": sub, "lowerFirst": utool.LowerFirst, "camelToSnakeCase": utool.CamelToSnakeCase}).Parse(entiryTemplate))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	// 将生成的代码输出为字符串
	return buf.String(), nil
}

func typeMap(typ string) string {
	typ = strings.TrimSpace(typ)
	typmap := map[string]string{
		"string": "String",
		"int":    "int",
	}
	if toTyp, ok := typmap[typ]; ok {
		return toTyp
	}
	if strings.HasPrefix(typ, "[]") {
		typ = strings.ReplaceAll(typ, "[]", "")
		typ = fmt.Sprintf("List<%s>", typ)
	}
	return typ
}
func sub(a, b int) int {
	return a - b
}
