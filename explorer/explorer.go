package explorer

import (
	"database/sql"
	"log"

	"github.com/MikeK123/dingo/model"
	_ "github.com/go-sql-driver/mysql"
)

type DatabaseExplorer interface {
	ExploreSchema(config *model.Configuration) (schema *model.DatabaseSchema)
}

type MySqlExplorer struct {
}

func NewMySqlExplorer() *MySqlExplorer {
	e := &MySqlExplorer{}
	return e
}

func (e *MySqlExplorer) ExploreSchema(config *model.Configuration) (schema *model.DatabaseSchema) {
	conn, err := sql.Open("mysql", config.Username+":"+config.Password+"@tcp("+config.Hostname+":"+config.Port+")/information_schema?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	schema = &model.DatabaseSchema{}
	schema.SchemaName = config.DatabaseName
	e.readTables(config, conn, schema)
	e.readViews(config, conn, schema)
	return schema
}

func (e *MySqlExplorer) readTables(config *model.Configuration, conn *sql.DB, schema *model.DatabaseSchema) {
	q := "SELECT TABLE_NAME FROM information_schema.TABLES Where TABLE_SCHEMA=? AND TABLE_TYPE='BASE TABLE' ORDER BY TABLE_NAME"
	rows, err := conn.Query(q, schema.SchemaName)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		table := &model.Table{}
		err := rows.Scan(&table.TableName)
		if err != nil {
			log.Fatal(err)
		}
		if config.IsIncluded(table.TableName) && !config.IsExcluded(table.TableName) {
			schema.Tables = append(schema.Tables, table)
			log.Printf("Examining table %s\r\n", table.TableName)
			e.readColums(conn, schema, table.TableName, &table.Columns)
			for _, col := range table.Columns {
				if col.IsPrimaryKey {
					table.PrimaryKeys = append(table.PrimaryKeys, col)
				} else {
					table.OtherColumns = append(table.OtherColumns, col)
				}
			}
		} else {
			log.Printf("Table %s is excluded\r\n", table.TableName)
		}
	}
}

func (e *MySqlExplorer) readViews(config *model.Configuration, conn *sql.DB, schema *model.DatabaseSchema) {
	q := "SELECT TABLE_NAME FROM information_schema.VIEWS Where TABLE_SCHEMA=? ORDER BY TABLE_NAME"
	rows, err := conn.Query(q, schema.SchemaName)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		view := &model.View{}
		err := rows.Scan(&view.ViewName)
		if err != nil {
			log.Fatal(err)
		}
		if config.IsIncluded(view.ViewName) && !config.IsExcluded(view.ViewName) {
			schema.Views = append(schema.Views, view)
			log.Printf("Examining view %s\r\n", view.ViewName)
			e.readColums(conn, schema, view.ViewName, &view.Columns)
		} else {
			log.Printf("View %s is excluded\r\n", view.ViewName)
		}
	}
}

func (e *MySqlExplorer) readColums(conn *sql.DB, schema *model.DatabaseSchema, tableName string, colums *[]*model.Column) {
	q := "SELECT TABLE_NAME, COLUMN_NAME, IS_NULLABLE, DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE, COLUMN_TYPE, COLUMN_KEY, EXTRA"
	q += " FROM information_schema.COLUMNS "
	q += " WHERE TABLE_SCHEMA=? AND TABLE_NAME=? ORDER BY ORDINAL_POSITION"
	rows, err := conn.Query(q, schema.SchemaName, tableName)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		column := &model.Column{}
		nullable := "NO"
		columnKey, extra := "", ""
		err := rows.Scan(&column.TableName, &column.ColumnName, &nullable, &column.DataType, &column.CharacterMaximumLength, &column.NumericPrecision, &column.NumericScale, &column.ColumnType, &columnKey, &extra)
		if err != nil {
			log.Fatal(err)
		}
		//log.Printf("Examining column %s\r\n", column.ColumnName)
		if "NO" == nullable {
			column.IsNullable = false
		} else {
			column.IsNullable = true
		}
		if "PRI" == columnKey {
			column.IsPrimaryKey = true
		}
		if "UNI" == columnKey {
			column.IsUnique = true
		}
		if "auto_increment" == extra {
			column.IsAutoIncrement = true
		}
		*colums = append(*colums, column)
	}
}
