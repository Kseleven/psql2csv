package pkg

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func ImportCsv2DB(conn *pgx.Conn, conf *Config) error {
	if !conf.ImportAction() {
		return fmt.Errorf("invalid action")
	}

	if err := checkTableHeader(conn, conf.ImportTable); err != nil {
		return err
	}
	defer conn.Close(context.Background())

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	for _, table := range conf.ImportTable.ImportTables {
		if err := table.importOneCsv2DB(tx); err != nil {
			tx.Rollback(context.Background())
			return fmt.Errorf("import table:%s failed:%s", table.TableName, err.Error())
		}
	}

	return tx.Commit(context.Background())
}

func checkTableHeader(conn *pgx.Conn, importConf *ImportTableConf) error {
	errMsgs := make([]string, 0, len(importConf.ImportTables))
	for _, table := range importConf.ImportTables {
		table.FileName = filepath.Join(importConf.ImportPath, table.FileName)
		if err := table.checkTableHeader(conn); err != nil {
			errMsgs = append(errMsgs, fmt.Sprintf("check table %s column failed:%s", table.TableName, err.Error()))
		}
	}

	if len(errMsgs) != 0 {
		return fmt.Errorf("check table header failed:\n%s", strings.Join(errMsgs, "\n"))
	}

	return nil
}

func (t *ImportTable) importOneCsv2DB(tx pgx.Tx) error {
	f, err := os.OpenFile(t.FileName, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file of %s failed:%s", t.TableName, err.Error())
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return fmt.Errorf("read csv from table %s failed:%s", t.TableName, err.Error())
	}
	if len(records) == 0 {
		return nil
	}

	placeHolders := make([]string, 0, len(records[0]))
	ignoreMap := make(map[int][]any, len(t.IgnoreColumns))
	for i := range t.Header {
		placeHolders = append(placeHolders, "$"+strconv.Itoa(i+1))
		for _, column := range t.IgnoreColumns {
			if column.Header == t.Header[i] {
				ignoreMap[i] = append(ignoreMap[i], column.Value)
			}
		}
	}

	var buf bytes.Buffer
	buf.WriteString("INSERT INTO ")
	buf.WriteString(t.TableName)
	buf.WriteString(" (")
	buf.WriteString(strings.Join(t.Header, ","))
	buf.WriteString(" )")
	buf.WriteString(" VALUES(")
	buf.WriteString(strings.Join(placeHolders, ","))
	buf.WriteString(")")

	br := &pgx.Batch{}
	for _, record := range records[1:] {
		vs := make([]any, 0, len(record))
		ignore := false
		for i, s := range record {
			if matchIgnoreColumn(ignoreMap, i, s) {
				ignore = true
				break
			}

			var slice []string
			if err := json.Unmarshal([]byte(s), &slice); err == nil {
				s = "{" + strings.Join(slice, ",") + "}"
			}
			if fd := t.TargetFds[i]; s == "" &&
				(fd.DataTypeOID == pgtype.InetArrayOID || fd.DataTypeOID == pgtype.TextArrayOID) {
				s = "{}"
			}

			vs = append(vs, s)
		}
		if ignore {
			continue
		}
		br.Queue(buf.String(), vs...)
	}

	result := tx.SendBatch(context.Background(), br)
	if err := result.Close(); err != nil {
		return fmt.Errorf("close failed:%s", err.Error())
	}
	return nil
}

func matchIgnoreColumn(ignoreMap map[int][]any, i int, s any) bool {
	if vs, ok := ignoreMap[i]; ok {
		for _, v := range vs {
			if reflect.DeepEqual(v, s) {
				return true
			}
		}
	}
	return false
}

type ImportTable struct {
	TableName     string         `json:"tableName"`
	FileName      string         `json:"fileName"`
	IgnoreColumns []IgnoreColumn `json:"ignoreColumns"`
	Header        []string
	TargetFds     []pgconn.FieldDescription
}

type IgnoreColumn struct {
	Header string `json:"header"`
	Value  any    `json:"value"`
}

func (t *ImportTable) String() string {
	return t.TableName + " : " + t.FileName
}

func (t *ImportTable) checkTableHeader(conn *pgx.Conn) error {
	if t.TableName == "" || t.FileName == "" {
		return fmt.Errorf("empty table:%s", t.TableName)
	}

	f, err := os.OpenFile(t.FileName, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	header, err := csv.NewReader(f).Read()
	if err != nil {
		return err
	}
	if len(header) == 0 {
		return fmt.Errorf("empty header")
	}

	rows, err := conn.Query(context.Background(), "select * from "+t.TableName+" limit 0")
	if err != nil {
		return err
	}
	defer rows.Close()

	targetHeader := make([]string, 0, len(rows.FieldDescriptions()))
	t.TargetFds = make([]pgconn.FieldDescription, len(rows.FieldDescriptions()))
	for _, fd := range rows.FieldDescriptions() {
		targetHeader = append(targetHeader, fd.Name)

		for i, s := range header {
			if s == fd.Name {
				t.TargetFds[i] = fd
				break
			}
		}
	}

	if err := columnDiff(header, targetHeader); err != nil {
		return fmt.Errorf("header:%v target:%v err:%s", header, targetHeader, err.Error())
	}

	t.Header = header
	return nil
}
