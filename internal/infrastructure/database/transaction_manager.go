package database

import (
	"context"

	"gorm.io/gorm"
	dbctx "github.com/dhanarrizky/Golang-template/pkg/database"
)

type GormTransactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) *GormTransactionManager {
	return &GormTransactionManager{db: db}
}

func (tm *GormTransactionManager) WithinTransaction(
	ctx context.Context,
	fn func(txCtx context.Context) error,
) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := dbctx.WithTx(ctx, tx)
		return fn(txCtx) // error → rollback otomatis
	})
}