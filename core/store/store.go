package store

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/peers-touch/peers-touch-go/core/option"
	"gorm.io/gorm"
)

// region Errors

// ErrDBNotFound is returned when the requested database is not found in the rdsMap.
var ErrDBNotFound = errors.New("database not found")
var ErrStoreNotDefined = errors.New("store not defined")
var ErrStoreInitFailed = func(name string) error { return fmt.Errorf("store[%s] init failed", name) }
var ErrStoreAlreadyInjected = errors.New("store already injected")

// endregion

var (
	store    Store
	initLock sync.Mutex
)

type Store interface {
	Init(ctx context.Context, opts ...option.Option) error
	RDS(ctx context.Context, opts ...RDSDMLOption) (*gorm.DB, error)
	Name() string
}

func GetRDS(ctx context.Context, opts ...RDSDMLOption) (*gorm.DB, error) {
	st, err := GetStore(ctx)
	if err != nil {
		return nil, err
	}

	return st.RDS(ctx, opts...)
}

func GetStore(ctx context.Context, opts ...GetOption) (Store, error) {
	return store, nil
}

// InjectStore is the entry point for saving Store implement for the application.
// It must be used after completing the initialization of Store.
// We use the store injected to access RDS or other resources.
func InjectStore(ctx context.Context, s Store) (err error) {
	initLock.Lock()
	defer initLock.Unlock()

	if store != nil {
		return ErrStoreAlreadyInjected
	}

	store = s

	return
}
