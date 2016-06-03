package {{.PackageName}}
{{if .HasImports}}
{{range .ImportPackages}}import "{{.}}"
{{end}}{{end}}
{{range .ModelTypes}}
// {{.TypeName}} data transfer object
type {{.TypeName}} struct {
	{{range .Fields}}{{.FieldName}} {{.FieldType}} {{if .FieldMetadata}}`{{.FieldMetadata}}`{{end}}
	{{end}}
}
{{end}}
{{range .ViewModelTypes}}
// {{.TypeName}} data transfer object
type {{.TypeName}} struct {
	{{range .Fields}}{{.FieldName}} {{.FieldType}} {{if .FieldMetadata}}`{{.FieldMetadata}}`{{end}}
	{{end}}
}

{{end}}