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
			field := &model.ModelField{FieldName: getModelFieldName(column.ColumnName), FieldType: getModelFieldType(pkg, column), FieldMetadata: getFieldMetadata(pkg, column)}
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
				if field.FieldType != "mysql.NullTime" { // exclude time fields
					field.IsNullable = true
					field.NullableFieldType = field.FieldType[8:] // scorporate sql.Null
				} else {
					field.IsNullable = true
					field.NullableFieldType = field.FieldType[10:] // scorporate mysql.Null
				}
			}
			mt.Fields = append(mt.Fields, field)
		}
	}
	for _, view := range schema.Views {
		mt := &model.ModelType{TypeName: getModelTypeName(view.ViewName), PackageName: "model"}
		pkg.ViewModelTypes = append(pkg.ViewModelTypes, mt)
		for _, column := range view.Columns {
			field := &model.ModelField{FieldName: getModelFieldName(column.ColumnName), FieldType: getModelFieldType(pkg, column), FieldMetadata: getFieldMetadata(pkg, column)}
			if column.IsNullable {
				if field.FieldType != "time.Time" { // exclude time fields
					field.IsNullable = true
				}
			}
			mt.Fields = append(mt.Fields, field)
		}
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
	name = strings.Replace(name, " ", "", -1)
	return name
}

func getModelFieldType(pkg *model.ModelPackage, column *model.Column) string {
	var ft string = ""
	switch column.DataType {
	case "char", "varchar", "enum", "text", "longtext", "mediumtext", "tinytext":
		if column.IsNullable {
			ft = "sql.NullString"
			pkg.AppendImport("database/sql")
		} else {
			ft = "string"
		}
	case "blob", "mediumblob", "longblob", "varbinary", "binary":
		ft = "[]byte"
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
			// MK: seams NullInt32 doesn't exist, replaced with NullInt64
			ft = "sql.NullInt64"
			pkg.AppendImport("database/sql")
		} else {
			// MK: Same here
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
		ft = "[]byte" // sql/driver/Value does not supports bool
	}
	if ft == "" {
		log.Printf("WARNING Incompatible Go type for column %s %s -> using string\r\n", column.ColumnName, column.ColumnType)
		ft = "string"
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
