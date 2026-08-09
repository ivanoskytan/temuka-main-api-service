package database

import (
	"context"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type txKey struct{}

func InjectTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func ExtractTx(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}

	return defaultDB
}

type PostgresWrapper struct {
	DB *gorm.DB
}

func NewPostgreSQL(pgHost, pgUser, pgPassword, pgPort, dbName string) (*PostgresWrapper, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		pgHost, pgUser, pgPassword, dbName, pgPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, fmt.Errorf("Failed to connect to database: %w", err)
	}
	log.Println("Database connection established successfully")

	return &PostgresWrapper{DB: db}, nil
}

func (p *PostgresWrapper) Create(ctx context.Context, value interface{}) error {
	db := ExtractTx(ctx, p.DB)
	return db.WithContext(ctx).Create(value).Error
}

func (p *PostgresWrapper) Save(ctx context.Context, value interface{}) error {
	db := ExtractTx(ctx, p.DB)
	return db.WithContext(ctx).Save(value).Error
}

func (p *PostgresWrapper) Delete(ctx context.Context, value interface{}, conds ...interface{}) error {
	db := ExtractTx(ctx, p.DB)
	return db.WithContext(ctx).Delete(value, conds...).Error
}

func (p *PostgresWrapper) Model(ctx context.Context, value interface{}) *gorm.DB {
	db := ExtractTx(ctx, p.DB)
	return db.WithContext(ctx).Model(value)
}

func (p *PostgresWrapper) First(ctx context.Context, dest interface{}, conds ...interface{}) error {
	db := ExtractTx(ctx, p.DB)
	return db.WithContext(ctx).First(dest, conds...).Error
}

func (p *PostgresWrapper) Find(ctx context.Context, dest interface{}, conds ...interface{}) error {
	db := ExtractTx(ctx, p.DB)
	return db.WithContext(ctx).Find(dest, conds...).Error
}

func (p *PostgresWrapper) Where(ctx context.Context, query interface{}, args ...interface{}) *gorm.DB {
	db := ExtractTx(ctx, p.DB)
	return db.WithContext(ctx).Where(query, args...)
}

func (p *PostgresWrapper) Exec(ctx context.Context, query string, args ...interface{}) error {
	db := ExtractTx(ctx, p.DB)
	return db.WithContext(ctx).Exec(query, args...).Error
}

func (p *PostgresWrapper) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return p.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := InjectTx(ctx, tx)
		return fn(txCtx)
	})
}
