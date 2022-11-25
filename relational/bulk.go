package relational

//type BulkInsert struct {
//	//db              *sql.DB
//	//table           *relational.Table
//	//batchInsertStmt *relational.BatchInsertStmt
//}
//
////func NewBulkInsert(db *sql.DB, table *relational.Table) (*BulkInsert, error) {
////	bi := &BulkInsert{
////		db:    db,
////		table: table,
////	}
////
////	batchInsertStmt, err := relational.NewBatchInsertStmt(bi.table.Name, bi.table.NotAutoIncrementColumnNames())
////	if err != nil {
////		return nil, err
////	}
////	bi.batchInsertStmt = batchInsertStmt
////	return bi, nil
////}
//
//func NewBulkInsert() (*BulkInsert, error) {
//	return &BulkInsert{}, nil
//}
//
//func (bi *BulkInsert) Insert(ctx context.Context, db *sql.DB, table *Table, rows []map[string]interface{}) error {
//	batchInsertStmt, err := NewBatchInsertStmt(table.Name, table.NotAutoIncrementColumnNames())
//	if err != nil {
//		return err
//	}
//
//	for _, row := range rows {
//		for _, field := range table.Fields {
//			if !field.AutoIncrement() {
//				batchInsertStmt.Set(field.Name, row[field.TypeName])
//			}
//		}
//	}
//
//	if batchInsertStmt.HaveBatch() {
//		if err := batchInsertStmt.ExecuteBatch(ctx, db); err != nil {
//			return err
//		}
//		batchInsertStmt.CleanBatch()
//	}
//
//	return nil
//}
