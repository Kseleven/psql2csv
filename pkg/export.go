package pkg

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5"
)

func Export2Csv(conn *pgx.Conn, conf *Config) error {
	if conf.ImportAction() {
		return fmt.Errorf("invalid action")
	}

	defer conn.Close(context.Background())
	for _, tableName := range conf.ExportTable.Tables {
		if err := exportOne2Csv(conn, tableName, conf.ExportTable.ExportPath); err != nil {
			return fmt.Errorf("export %s failed:%s", tableName, err.Error())
		}
	}

	return nil
}

func exportOne2Csv(conn *pgx.Conn, tableName, exportPath string) error {
	rows, err := conn.Query(context.Background(), "select * from "+tableName)
	if err != nil {
		return err
	}
	defer rows.Close()

	header := make([]string, 0, len(rows.FieldDescriptions()))
	fdTypes := make([]uint32, 0, len(rows.FieldDescriptions()))
	for _, fd := range rows.FieldDescriptions() {
		header = append(header, fd.Name)
		fdTypes = append(fdTypes, fd.DataTypeOID)
	}

	results := make([][]string, 0, len(rows.FieldDescriptions()))
	for rows.Next() {
		v, err := rows.Values()
		if err != nil {
			return err
		}

		result := make([]string, 0, len(rows.FieldDescriptions()))
		for i, a := range v {
			r := valueToString(a)
			if r == "" && fdTypes[i] == pgtype.TextArrayOID {
				r = "{}"
			}
			result = append(result, r)
		}
		results = append(results, result)
	}

	f, err := os.OpenFile(filepath.Join(exportPath, tableName+".csv"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	if err := w.Write(header); err != nil {
		return fmt.Errorf("write header failed:%s", err.Error())
	}
	for _, r := range results {
		if err := w.Write(r); err != nil {
			return fmt.Errorf("write contents failed:%s", err.Error())
		}
	}

	w.Flush()
	return w.Error()
}
