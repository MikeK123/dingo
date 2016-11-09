package producers

import (
	"bytes"
	"log"
	"strings"

	"github.com/MikeK123/dingo/model"
)

func ProduceModelPackage(config *model.Configuration, schema *model.DatabaseSchema) (pkg *model.ModelPackage) {
	pkg = &model.ModelPackage{PackageName: "model", BasePackage: config.BasePackage}
	for _, table := range schema.Tables {
		mt := &model.ModelType{TypeName: getModelTypeName(table.TableName), PackageName: "model"}
		pkg.ModelTypes = append(pkg.ModelTypes, mt)
		for _, column := range table.Columns {
			field := &model.ModelField{FieldName: getModelFieldName(column.ColumnName), FieldType: getModelFieldType(config.DatabaseType, pkg, column), FieldMetadata: getFieldMetadata(pkg, column)}
			if column.IsPrimaryKey {
				field.IsPK = true
				mt.PKFields = append(mt.PKFields, field)
			} else {
				mt.OtherFields = append(mt.OtherFields, field)
			}
			if column.IsAutoIncrement {
				field.IsAutoInc = true
			}
			if column.IsNullable {
				if field.FieldType == "mysql.NullTime" {
					field.IsNullable = true
					field.NullableFieldType = field.FieldType[10:] // scorporate mysql.Null
				} else if field.FieldType == "pq.NullTime" {
					field.IsNullable = true
					field.NullableFieldType = field.FieldType[7:] // scorporate pq.Null
				} else if field.FieldType == "[]byte" {
					field.IsNullable = true
					field.NullableFieldType = "[]byte" // byte slice
				} else {
					field.IsNullable = true
					field.NullableFieldType = field.FieldType[8:] // scorporate sql.Null
				}
			}
			mt.Shortcut = getTitleLetters(mt.TypeName)
			mt.Fields = append(mt.Fields, field)
		}
	}
	for _, view := range schema.Views {
		mt := &model.ModelType{TypeName: getModelTypeName(view.ViewName), PackageName: "model"}
		pkg.ViewModelTypes = append(pkg.ViewModelTypes, mt)
		for _, column := range view.Columns {
			field := &model.ModelField{FieldName: getModelFieldName(column.ColumnName), FieldType: getModelFieldType(config.DatabaseType, pkg, column), FieldMetadata: getFieldMetadata(pkg, column)}
			if column.IsNullable {
				// if field.FieldType != "time.Time" { // exclude time fields
				// 	field.IsNullable = true
				// }
				// MK: follow same procedure as with standard tables
				if field.FieldType == "mysql.NullTime" {
					field.IsNullable = true
					field.NullableFieldType = field.FieldType[10:] // scorporate mysql.Null
				} else if field.FieldType == "pq.NullTime" {
					field.IsNullable = true
					field.NullableFieldType = field.FieldType[7:] // scorporate pq.Null
				} else if field.FieldType == "[]byte" {
					field.IsNullable = true
					field.NullableFieldType = "[]byte" // byte slice
				} else {
					field.IsNullable = true
					field.NullableFieldType = field.FieldType[8:] // scorporate sql.Null
				}
			}
			mt.Fields = append(mt.Fields, field)
		}
	}
	return pkg
}

func ProduceMixedModelPackage(config *model.Configuration) (pkg *model.ModelPackage) {
	pkg = &model.ModelPackage{PackageName: "model", BasePackage: config.BasePackage}
	for _, mdt := range config.MixedDaoTables {
		mt := &model.ModelType{PackageName: "model"}
		for i, tn := range mdt.Tables {
			field := &model.ModelField{
				FieldName: strings.Title(mdt.Shortcuts[i]),
				FieldType: getModelTypeName(tn),
			}
			//mt.TypeName += getModelTypeName(tn)
			mt.TypeName += getThreeLetters(getModelTypeName(tn))
			mt.Fields = append(mt.Fields, field)
			mt.Shortcut = mdt.Shortcuts[i]
		}
		pkg.ModelTypes = append(pkg.ModelTypes, mt)
	}
	return pkg
}

func getModelTypeName(tablename string) string {
	name := strings.Replace(tablename, "-", " ", -1)
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	name = strings.Replace(name, " ", "", -1)
	return name
}

func getModelFieldName(fieldname string) string {
	name := strings.Replace(fieldname, "-", " ", -1)
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	// MK: ifsuffix is id thenreplace with capital ID
	if strings.HasSuffix(name, " Id") {
		name = strings.TrimSuffix(name, "Id")
		name += "ID"
	}
	if name == "Id" {
		name = "ID"
	}
	name = strings.Replace(name, " ", "", -1)
	return name
}

func getModelFieldType(databaseType string, pkg *model.ModelPackage, column *model.Column) string {
	switch databaseType {
	case "mysql":
		return getMySQLModelFieldType(pkg, column)
	case "postgres":
		return getPostgresModelFieldType(pkg, column)
	default:
		return getMySQLModelFieldType(pkg, column)
	}

}

func getMySQLModelFieldType(pkg *model.ModelPackage, column *model.Column) string {
	var ft string
	switch column.DataType {
	case "char", "varchar", "enum", "text", "longtext", "mediumtext", "tinytext":
		if column.IsNullable {
			ft = "sql.NullString"
			pkg.AppendImport("database/sql")
		} else {
			ft = "string"
		}
	case "blob", "mediumblob", "longblob", "varbinary", "binary":
		// MK: had to replace this one with string as Nullable type is tricky
		//ft = "[]byte"
		if column.IsNullable {
			ft = "sql.NullString"
			pkg.AppendImport("database/sql")
		} else {
			ft = "string"
		}
	case "date", "time", "datetime", "timestamp":
		if column.IsNullable {
			ft = "mysql.NullTime"
			pkg.AppendImport("github.com/go-sql-driver/mysql")
		} else {
			ft = "time.Time"
			pkg.AppendImport("time")
		}
	case "tinyint", "smallint":
		if column.IsNullable {
			ft = "sql.NullInt64"
			pkg.AppendImport("database/sql")
		} else {
			ft = "int64"
		}
	case "int", "mediumint", "bigint":
		if column.IsNullable {
			ft = "sql.NullInt64"
			pkg.AppendImport("database/sql")
		} else {
			ft = "int64"
		}
	case "float", "decimal", "double":
		if column.IsNullable {
			ft = "sql.NullFloat64"
			pkg.AppendImport("database/sql")
		} else {
			ft = "float64"
		}
	case "bit":
		column.IsNullable = false // nullable column of this type is not managed
		ft = "[]byte"             // sql/driver/Value does not supports bool
	}
	if ft == "" {
		log.Printf("WARNING Incompatible Go type for MySQL column %s %s -> using string\r\n", column.ColumnName, column.ColumnType)
		if column.IsNullable {
			ft = "sql.NullString"
			pkg.AppendImport("database/sql")
		} else {
			ft = "string"
		}
	}
	return ft
}

func getPostgresModelFieldType(pkg *model.ModelPackage, column *model.Column) string {
	var ft string
	switch column.ColumnType {
	case "char", "varchar", "text", "character":
		if column.IsNullable {
			ft = "sql.NullString"
			pkg.AppendImport("database/sql")
		} else {
			ft = "string"
		}
	case "bytea":
		column.IsNullable = false // nullable column of this type is not managed
		ft = "[]byte"
	case "date", "time", "timetz", "timestamptz", "timestamp", "interval":
		if column.IsNullable {
			ft = "pq.NullTime"
			pkg.AppendImport("github.com/lib/pq")
		} else {
			ft = "time.Time"
			pkg.AppendImport("time")
		}
	case "int2", "int4":
		if column.IsNullable {
			ft = "sql.NullInt64"
			pkg.AppendImport("database/sql")
		} else {
			ft = "int64"
		}
	case "int8":
		if column.IsNullable {
			ft = "sql.NullInt64"
			pkg.AppendImport("database/sql")
		} else {
			ft = "int64"
		}
	case "float4", "float8", "numeric":
		if column.IsNullable {
			ft = "sql.NullFloat64"
			pkg.AppendImport("database/sql")
		} else {
			ft = "float64"
		}
	case "bit", "bool":
		column.IsNullable = false // nullable column of this type is not managed
		ft = "string"             // seems that lib/pq does not supports bool ?
	}
	if ft == "" {
		log.Printf("WARNING Incompatible Go type for Postgres column %s %s -> using string\r\n", column.ColumnName, column.ColumnType)
		if column.IsNullable {
			ft = "sql.NullString"
			pkg.AppendImport("database/sql")
		} else {
			ft = "string"
		}
	}
	return ft
}

// Generate GORM like metadata
func getFieldMetadata(pkg *model.ModelPackage, column *model.Column) string {
	var buffer bytes.Buffer
	buffer.WriteString("sql:\"")
	buffer.WriteString("type:")
	buffer.WriteString(column.ColumnType)
	if !column.IsNullable {
		buffer.WriteString(";not null")
	}
	if column.IsUnique {
		buffer.WriteString(";unique")
	}
	if column.IsAutoIncrement {
		buffer.WriteString(";AUTO_INCREMENT")
	}
	buffer.WriteString("\"")
	return buffer.String()
}

func getTitleLetters(s string) string {
	var res string
	for _, l := range s {
		ls := string(l)
		if ls == strings.ToUpper(ls) {
			res += ls
		}
	}
	if res == "" && s != "" {
		res = string(s[0])
	}
	return res
}

func getThreeLetters(s string) string {
	var res string
	var i int
	for _, r := range s {
		sr := string(r)
		i++
		if sr == strings.ToUpper(sr) {
			i = 0
		}
		if i < 3 {
			res += sr
		}
	}
	return res
}
