package cmd

import (
	"database/sql"
	"fmt"
	_ "gorm.io/driver/mysql"
	"os"
	"reflect"
	"strings"
)

var db *sql.DB

func initDB(h, p, u, s string) error {
	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", u, s, h, p, "information_schema")
	var err error
	db, err = sql.Open("mysql", dbURL)
	if err != nil {
		fmt.Println(dbURL)
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	return nil
}

type column struct {
	Field string
	Type  string
}

type tableColumn struct {
	TableName     string
	ColumnName    string
	DataType      string
	ColumnComment string
}

func tables(dbName string) (map[string][]*tableColumn, error) {
	rows, err := db.Query("SELECT TABLE_NAME, COLUMN_NAME, DATA_TYPE, COLUMN_COMMENT FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ?", dbName)
	if err != nil {
		return nil, err
	}
	result := make(map[string][]*tableColumn)
	for rows.Next() {
		tname, cname, dtype, ccomment := "","","",""
		err = rows.Scan(&tname, &cname, &dtype, &ccomment)
		if err != nil {
			return nil, err
		}
		column := &tableColumn{
			TableName:     tname,
			ColumnName:    cname,
			DataType:      dtype,
			ColumnComment: ccomment,
		}
		if _, ok := result[column.TableName]; ok {
			result[column.TableName] = append(result[column.TableName], column)
		} else {
			result[column.TableName] = []*tableColumn{column}
		}
	}
	return result, nil
}

func synTable(database, dst, pkg, prefix string) error {
	tbls, err := tables(database)
	if err != nil {
		return err
	}
	if len(tbls) == 0 {
		return fmt.Errorf("空数据库")
	}
	for tableName, table := range tbls {
		columns := make([]column, 0)
		for _, c := range table {
			columns = append(columns, column{
				Field: c.ColumnName,
				Type:  c.DataType,
			})
		}
		_, err = writeToTable(dst, pkg, strings.TrimLeft(tableName, prefix), columns)
		if err != nil {
			return err
		}
	}
	return nil
}
func toCamel(s string) string {
	buf := strings.Builder{}
	for i, b := range s {
		if b == '_' {
			continue
		}
		if i == 0 || s[i-1] == '_' {
			b -= 32
		}
		buf.WriteRune(b)
	}
	return buf.String()
}


func mysqlTypeToGoType(mysqlType string) string {
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint":
		return reflect.Int.String()
	case "bigint":
		return reflect.Int64.String()
	case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext", "json":
		return reflect.String.String()
	case "date", "datetime", "time", "timestamp":
		return "time.Time"
	case "decimal", "double":
		return reflect.Float64.String()
	case "float":
		return reflect.Float32.String()
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		return "[]byte"
	}
	return ""
}

func writeToTable(dst, pkg, table string, c []column) (string, error) {
	packageName := pkg
	entityName := toCamel(table)
	structFieldTpl := make([]string, 0)
	fields := make([]interface{}, 0)
	importTime := ""
	for _, v := range c {
		structFieldTpl = append(structFieldTpl, "%s")
		t := mysqlTypeToGoType(v.Type)
		fields = append(fields, fmt.Sprintf("    %s %s", toCamel(v.Field), t))
		if t == "time.Time" {
			importTime = "\nimport \"time\"\n"
		}
	}
	trueTemplate := fmt.Sprintf(template, packageName, importTime, entityName, strings.Join(structFieldTpl, "\n"), entityName, table)
	fileContent := fmt.Sprintf(trueTemplate, fields...)
	_, err := os.Stat(dst)
	if err != nil {
		return "", err
	}
	fileName := dst + "/" + table + ".go"
	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	_, err = f.WriteString(fileContent)
	if err != nil {
		return "", err
	}
	return fileContent, nil
}

const template = `package %s
%s
type %s struct {
%s
}

func (t %s) TableName() string {
	return "%s"
}
`
