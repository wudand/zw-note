package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// namedReturningID runs a named INSERT ... RETURNING id and scans the new id.
func namedReturningID(ctx context.Context, q sqlx.QueryerContext, query string, arg interface{}, id *uint64) error {
	bound, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}
	bound = sqlx.Rebind(sqlx.DOLLAR, bound)
	return sqlx.GetContext(ctx, q, id, bound, args...)
}
