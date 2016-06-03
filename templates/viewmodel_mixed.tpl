package {{.PackageName}}
{{if .HasImports}}
{{range .ImportPackages}}import "{{.}}"
{{end}}{{end}}
{{range .ViewModelTypes}}
// {{.TypeName}} view-Model object
type {{.TypeName}} struct {
	{{range .Fields}}{{.FieldName}} {{.FieldType}}
	{{end}}
}

{{end}}