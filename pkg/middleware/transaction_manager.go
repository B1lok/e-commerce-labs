package middleware

import (
	"context"
	"gorm.io/gorm"
)

type TxManager interface {
	ExecuteTx(ctx context.Context, fn func(txCtx context.Context) error) error
}

type txManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) TxManager {
	return &txManager{db: db}
}

type txKey struct{}

func (m *txManager) ExecuteTx(ctx context.Context, fn func(txCtx context.Context) error) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}

func GetTx(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	return tx, ok
}
