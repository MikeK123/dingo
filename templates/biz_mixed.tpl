package {{.PackageName}}
{{if .HasImports}}
{{range .ImportPackages}}import "{{.}}"
{{end}}{{end}}
{{range .BizTypes}}
// {{.TypeName}} is a business object for {{.Model.TypeName}} entities.
type {{.TypeName}} struct {
	{{range .Fields}}{{.FieldName}} {{.FieldType}}{{end}}
}
//New{{.TypeName}} creates a {{.TypeName}}
func New{{.TypeName}}() *{{.TypeName}} {
	return &{{.TypeName}}{ Dao:&{{.Dao.PackageName}}.{{.Dao.TypeName}}{} }
}

// ToViewModel converts a model entity in a view-model
func (b *{{.TypeName}}) ToViewModel(m *{{.Model.PackageName}}.{{.Model.TypeName}}) *{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}{
	v := &{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}{}
	{{range .Model.Fields}}v.{{.FieldName}} = *New{{.FieldType}}Biz().ToViewModel(&m.{{.FieldName}})
	{{end}}
	return v
}

// ToModel converts a view-model in a model entity
func (b *{{.TypeName}}) ToModel(v *{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}) *{{.Model.PackageName}}.{{.Model.TypeName}}{
	m := &{{.Model.PackageName}}.{{.Model.TypeName}}{}
	{{range .Model.Fields}}m.{{.FieldName}} = *New{{.FieldType}}Biz().ToModel(&v.{{.FieldName}})
	{{end}}
	return m
}

// List the {{.Model.TypeName}} entities.
func (b *{{.TypeName}}) List(take int, skip int) (list []*{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}, err error) {
	mlist, err := b.Dao.List(dao.Connection, take, skip)
	if err != nil {
		return nil, err
	}
	for _, m := range mlist {
		list = append(list, b.ToViewModel(m))
	}
	return list, nil
}
{{end}}