package relational

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type BatchInsertStmt struct {
	tableName   string
	columnNames []string
	args        map[string]interface{}
	rows        [][]interface{}
}

func NewBatchInsertStmt(table string, columnNames []string) (*BatchInsertStmt, error) {
	s := new(BatchInsertStmt)
	s.tableName = table
	for _, name := range columnNames {
		s.columnNames = append(s.columnNames, name)
	}
	return s, nil
}

func (s *BatchInsertStmt) AddBatch() {
	row := make([]interface{}, len(s.columnNames))
	for i, name := range s.columnNames {
		row[i] = s.args[name]
	}
	s.rows = append(s.rows, row)
	s.args = nil
}

func (s *BatchInsertStmt) ExecuteBatch(ctx context.Context, db *sql.DB) error {
	batchSQL := s.batchSQL()
	// TODO
	fmt.Println(batchSQL[:100])
	_, err := db.ExecContext(ctx, batchSQL)
	return err
}

func (s *BatchInsertStmt) CleanBatch() {
	s.rows = nil
}

func (s *BatchInsertStmt) HaveBatch() bool {
	return len(s.rows) != 0
}

func (s *BatchInsertStmt) Set(colName string, v interface{}) {
	if s.args == nil {
		s.args = make(map[string]interface{}, len(s.columnNames))
	}
	s.args[colName] = v
}

func (s *BatchInsertStmt) batchSQL() string {
	dml := fmt.Sprintf("INSERT INTO %s (%s) VALUES ", s.tableName, strings.Join(s.columnNames, ", "))

	var rowValues []string
	for _, row := range s.rows {
		sb := new(strings.Builder)
		sb.WriteString("(")
		for i, column := range row {
			switch v := column.(type) {
			case int:
				sb.WriteString(strconv.Itoa(v))
			case string:
				sb.WriteString(strconv.Quote(v))
			}
			if i != len(row)-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString(")")
		rowValues = append(rowValues, sb.String())
	}
	dml = dml + strings.Join(rowValues, ", ")

	return dml
}
