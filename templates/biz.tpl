package {{.PackageName}}
{{if .HasImports}}
{{range .ImportPackages}}import "{{.}}"
{{end}}{{end}}
{{range .BizTypes}}
// Businness object for {{.Model.TypeName}} entities.
type {{.TypeName}} struct {
	{{range .Fields}}{{.FieldName}} {{.FieldType}}
	{{end}}
}
func New{{.TypeName}}() *{{.TypeName}} {
	return &{{.TypeName}}{ Dao:&{{.Dao.PackageName}}.{{.Dao.TypeName}}{} }
}
// Convert an model entity in a view-model
func (b *{{.TypeName}}) ToViewModel(m *{{.Model.PackageName}}.{{.Model.TypeName}}) *{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}{
	v := &{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}{}
	{{range .Model.Fields}}{{if .IsNullable}}if m.{{.FieldName}}.Valid { v.{{.FieldName}} = m.{{.FieldName}}.{{.NullableFieldType}} }{{else}}v.{{.FieldName}} = m.{{.FieldName}}{{end}}
	{{end}}
	return v
}
// Convert a view-model in a model entity
func (b *{{.TypeName}}) ToModel(v *{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}) *{{.Model.PackageName}}.{{.Model.TypeName}}{
	m := &{{.Model.PackageName}}.{{.Model.TypeName}}{}
	{{range .Model.Fields}}{{if .IsNullable}}m.{{.FieldName}} =  {{.FieldType}} { {{.NullableFieldType}}:v.{{.FieldName}} }{{else}}m.{{.FieldName}} = v.{{.FieldName}}{{end}}
	{{end}}
	return m
}
{{if .IsReadOnly}}{{else}}// Insert a new {{.Model.TypeName}} entity and returns the last insert Id.
func (b *{{.TypeName}}) Insert(v *{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}) (lastInsertId int64, err error) {
	return b.Dao.Insert(dao.Connection, b.ToModel(v))
}
// Update a {{.Model.TypeName}} entity and returns the number of affected rows.
func (b *{{.TypeName}}) Update(v *{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}) (rowsAffected int64, err error) {
	return b.Dao.Update(dao.Connection, b.ToModel(v))
}
// Delete a {{.Model.TypeName}} entity and returns the number of affected rows.
func (b *{{.TypeName}}) Delete(v *{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}) (rowsAffected int64, err error) {
	return b.Dao.Delete(dao.Connection, b.ToModel(v))
}
// Find the {{.Model.TypeName}} entity by primary keys, returns nil if not found.
func (b *{{.TypeName}}) FindByPrimaryKey({{range $i, $e := .Model.PKFields}}{{if $i}}, {{end}}{{.FieldName}} {{.FieldType}}{{end}}) (v *{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}, err error){
	m, err := b.Dao.FindByPrimaryKey(dao.Connection, {{range $i, $e := .Model.PKFields}}{{if $i}}, {{end}}{{.FieldName}}{{end}})
	if err != nil {
		return nil, err
	} else {
		return b.ToViewModel(m), nil
	}
}{{end}}
// List the {{.Model.TypeName}} entities.
func (b *{{.TypeName}}) List(take int32, skip int32) (list []*{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}, err error) {
	mlist, err := b.Dao.List(dao.Connection, take, skip)
	if err != nil {
		return nil, err
	} else {
		for _, m := range mlist {
			list = append(list, b.ToViewModel(m))
		}
		return list, nil
	}
}
// Count the {{.Model.TypeName}} entities.
func (b *{{.TypeName}}) Count() (count int64, err error){
	return b.Dao.Count(dao.Connection)
}
{{end}}