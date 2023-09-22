package datasource

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
	"sync"
)

type SQLCommon interface {
	sqlx.Ext
	Select(dest interface{}, query string, args ...interface{}) error
	// Get using this DB.
	// Any placeholder parameters are replaced with supplied args.
	// An error is returned if the result set is empty.
	Get(dest interface{}, query string, args ...interface{}) error
}
type DBCommon interface {
	SQLCommon
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

type TransactionManger struct {
	db    DBCommon
	mutex sync.Mutex
}

func NewTransactionManger(db DBCommon) *TransactionManger {
	return &TransactionManger{db: db}
}

// TransactionHolder 事务holder
//cnt 为嵌套计数，用于处理嵌套事务，这里没有考虑事务隔离级别和事务传播机制
type TransactionHolder struct {
	DB  DBCommon
	Tx  *sqlx.Tx
	Cnt int
}

func (t *TransactionManger) BeginTx(ctx *gin.Context) error {
	if ctx == nil {
		log.Debug("没有ctx,不开启事务")
		return nil
	}
	//从Context中获取transaction holder,如果没有则新增一个

	txHolder, exists := ctx.Get(transactionContextKey)
	if !exists {
		t.mutex.Lock()
		defer t.mutex.Unlock()
		txHolder, exists = ctx.Get(transactionContextKey)
		if !exists {
			txHolder = &TransactionHolder{
				DB: t.db,
				Tx: nil,
			}
			ctx.Set(transactionContextKey, txHolder)
		}
	}

	th, ok := txHolder.(*TransactionHolder)
	if !ok {
		return fmt.Errorf("error:获取事务管理错误:%v", txHolder)
	}

	if th.Tx == nil {
		var err error
		th.Tx, err = t.db.BeginTxx(ctx, nil)
		if err != nil {
			panic(fmt.Sprintf("error: 开启事务错误:%v\n", err))
		}
		th.Cnt++
	}
	return nil
}

func (*TransactionManger) CommitTx(ctx *gin.Context) error {
	if ctx == nil {
		log.Debug("没有ctx,不提交事务")
		return nil
	}
	txHolder, exists := ctx.Get(transactionContextKey)
	if !exists {
		log.Debug("error: 没有开启的事务.")
	}
	th, ok := txHolder.(*TransactionHolder)
	if !ok {
		return fmt.Errorf("error:获取事务管理错误:%v\n", txHolder)
	}
	th.Cnt--
	if th.Cnt != 0 {
		return nil
	}
	err := th.Tx.Commit()
	if err != nil {
		return fmt.Errorf("error: 提交事务错误:%v\n", err)
	}
	return nil
}

func (t *TransactionManger) GetConn(ctx *gin.Context) (SQLCommon, error) {
	if ctx == nil {
		return t.db, nil
	}

	txHolder, exists := ctx.Get(transactionContextKey)
	if !exists {
		//没有开启事务，则用普通连接
		return t.db, nil
	}
	th, ok := txHolder.(*TransactionHolder)
	if !ok {
		return nil, fmt.Errorf("error:获取事务管理错误:%v\n", txHolder)
	}
	if th.Tx != nil {
		return th.Tx, nil
	}
	return t.db, nil
}

//TransactionMiddleware 事务处理中间件
func TransactionMiddleware(ctx *gin.Context) {

	defer func() {
		if r := recover(); r != nil {
			txHolder, exists := ctx.Get(transactionContextKey)
			if exists {
				th, ok := txHolder.(*TransactionHolder)
				if ok && th.Tx != nil {
					err := th.Tx.Rollback()
					if err != nil {
						log.Debug("error:%v", err)
					}

				}

			}
			panic(r)
		}

	}()
	ctx.Next()
}
