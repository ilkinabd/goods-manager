package dao

import (
	"context"
	"github.com/ilkinabd/goods-manager/app/internal/domain/product/filter"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type PostgreSQLClient interface {
	Begin(context.Context) (pgx.Tx, error)
	BeginFunc(ctx context.Context, f func(pgx.Tx) error) error
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) error
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

type ProductDAO interface {
	All(context.Context, []filter.Criteria, filter.Sortable) ([]*Product, error)
	One(context.Context, string) (*Product, error)
	Create(context.Context, map[string]interface{}) error
	Update(context.Context, string, map[string]interface{}) error
	Delete(context.Context, string) error
}
