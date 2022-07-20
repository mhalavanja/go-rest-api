package db

import (
	"context"
	"database/sql"
	"fmt"

	"main.go/db/sqlc"
)

type Store struct {
	*sqlc.Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: sqlc.New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := sqlc.New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

func (store *Store) LeaveGroup(ctx context.Context, arg sqlc.LeaveGroupParams) error {
	err := store.execTx(ctx, func(q *sqlc.Queries) error {
		if err := q.LeaveGroup(ctx, arg); err != nil {
			return err
		}

		numOfPeople, err := q.DecNumOfPeople(ctx, arg.GroupID)
		if err != nil {
			return err
		}

		if numOfPeople > 0 {
			return nil
		}

		if err := q.DeleteGroup(ctx, arg.GroupID); err != nil {
			return err
		}

		return nil
	})

	return err
}
