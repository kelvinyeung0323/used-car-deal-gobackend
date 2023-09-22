package datasource

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"sync"
)

const transactionContextKey = "transaction-context-key"

var mysqlMutex = sync.Mutex{}

type MysqlConn struct {
	*sqlx.DB
	dsn string
}

func NewMysqlConn(dsn string) (*MysqlConn, error) {
	conn := &MysqlConn{dsn: dsn}
	err := conn.initDB()
	return conn, err
}

func (c *MysqlConn) initDB() error {
	// dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=UTC",
	// "root", "123456", "localhost", 3306, "go-tunnel")
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?allowNativePasswords=True&charset=utf8&parseTime=True&loc=UTC",
	// "root", "123456", "localhost", 33060, "go-tunnel")
	//dsn := config.MySQLConfig.Username + ":" + config.MySQLConfig.Password + "@" + config.MySQLConfig.Url
	// sqlx连接池
	var err error
	c.DB, err = sqlx.Connect("mysql", c.dsn)
	if err != nil {
		return err
	}
	c.DB.SetMaxOpenConns(500)
	c.DB.SetMaxIdleConns(100)
	return nil
}
