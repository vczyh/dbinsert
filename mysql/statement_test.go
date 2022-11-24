package mysql

import "testing"

func TestNewClientPreparedStmt(t *testing.T) {
	stmt, err := NewBatchInsertStmt("tbl", []string{"id", "val"})
	if err != nil {
		t.Fatal(err)
	}

	stmt.Set("id", 1)
	stmt.Set("val", "a")
	stmt.AddBatch()

	stmt.Set("id", 2)
	stmt.Set("val", "b")
	stmt.AddBatch()

	batchSQL := stmt.batchSQL()
	t.Logf(batchSQL)
}
