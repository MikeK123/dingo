package {{.PackageName}}
{{if .HasImports}}
{{range .ImportPackages}}import "{{.}}"
{{end}}{{end}}

{{range .DaoMixedTypes}}
// ---------------------------- {{.TypeName}} ----------------------------

// {{.TypeName}} is a data access object for mixed {{range $i, $m := .Model}}{{if $i}}, {{end}}{{.TypeName}}{{end}} entities.
type {{.TypeName}} struct {
	{{range .Fields}}{{.FieldName}} {{.FieldType}} `{{.FieldMetadata}}`
	{{end}}
}

// List the mixed {{range $i, $m := .Model}}{{if $i}}, {{end}}{{.TypeName}}{{end}} entities.
func (dao *{{.TypeName}}) List(conn *sql.DB, take int, skip int, whereEx string) (list []*model.{{.Shortcut}}, err error) {
	q := "SELECT {{range $j, $e := .Entity}}{{if $j}}, {{end}}{{range $i, $c := .Columns}}{{if $i}}, {{end}}`{{$e.TableShortcut}}`.`{{.ColumnName}}`{{end}}{{end}} FROM {{range $j, $e := .Entity}}{{if $j}}, {{end}}`{{.TableName}}` `{{.TableShortcut}}`{{end}} WHERE {{.Where}}"
	if whereEx != "" {
		q += " AND " + whereEx
	}
	q += " LIMIT ? OFFSET ?"
	rows, err := conn.Query(q, take, skip)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		{{range $i, $m := .Model}}dto{{$i}} := &{{.PackageName}}.{{.TypeName}}{}
        {{end}}
		err := rows.Scan({{range $j, $m := .Model}}{{if $j}}, {{end}}{{range $i, $e := .Fields}}{{if $i}}, {{end}}&dto{{$j}}.{{.FieldName}}{{end}}{{end}})
		if err != nil {
			return nil, err
		}

		dto := &model.{{.Shortcut}}{		
			{{range $i, $m := .Model}}{{.Shortcut}}: *dto{{$i}},
			{{end}}
		}
		list = append(list, dto)
	}
	return list, nil
}

{{end}}

