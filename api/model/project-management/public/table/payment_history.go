//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var PaymentHistory = newPaymentHistoryTable("public", "payment_history", "")

type paymentHistoryTable struct {
	postgres.Table

	// Columns
	ID        postgres.ColumnInteger
	ProjectID postgres.ColumnString
	Amount    postgres.ColumnInteger
	CreatedAt postgres.ColumnTimestampz
	UpdatedAt postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type PaymentHistoryTable struct {
	paymentHistoryTable

	EXCLUDED paymentHistoryTable
}

// AS creates new PaymentHistoryTable with assigned alias
func (a PaymentHistoryTable) AS(alias string) *PaymentHistoryTable {
	return newPaymentHistoryTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new PaymentHistoryTable with assigned schema name
func (a PaymentHistoryTable) FromSchema(schemaName string) *PaymentHistoryTable {
	return newPaymentHistoryTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new PaymentHistoryTable with assigned table prefix
func (a PaymentHistoryTable) WithPrefix(prefix string) *PaymentHistoryTable {
	return newPaymentHistoryTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new PaymentHistoryTable with assigned table suffix
func (a PaymentHistoryTable) WithSuffix(suffix string) *PaymentHistoryTable {
	return newPaymentHistoryTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newPaymentHistoryTable(schemaName, tableName, alias string) *PaymentHistoryTable {
	return &PaymentHistoryTable{
		paymentHistoryTable: newPaymentHistoryTableImpl(schemaName, tableName, alias),
		EXCLUDED:            newPaymentHistoryTableImpl("", "excluded", ""),
	}
}

func newPaymentHistoryTableImpl(schemaName, tableName, alias string) paymentHistoryTable {
	var (
		IDColumn        = postgres.IntegerColumn("id")
		ProjectIDColumn = postgres.StringColumn("project_id")
		AmountColumn    = postgres.IntegerColumn("amount")
		CreatedAtColumn = postgres.TimestampzColumn("created_at")
		UpdatedAtColumn = postgres.TimestampzColumn("updated_at")
		allColumns      = postgres.ColumnList{IDColumn, ProjectIDColumn, AmountColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns  = postgres.ColumnList{ProjectIDColumn, AmountColumn, CreatedAtColumn, UpdatedAtColumn}
	)

	return paymentHistoryTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:        IDColumn,
		ProjectID: ProjectIDColumn,
		Amount:    AmountColumn,
		CreatedAt: CreatedAtColumn,
		UpdatedAt: UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}