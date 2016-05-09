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
{{if .IsSimplePK}}
// ConvertPK will convert string value to primary key native type value 
func (vm *{{.TypeName}}) ConvertPK(value string) {{.PKType}} {
	{{.PKStringConv}}
	return ret
}
{{end}}
{{end}}