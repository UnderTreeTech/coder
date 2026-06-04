package dao

import (
	"strings"

	"github.com/UnderTreeTech/drivers"
	//_ "github.com/alexbrainman/odbc"
	_ "github.com/UnderTreeTech/drivers/dm"

	_ "github.com/UnderTreeTech/drivers/kingbase.com/gokb"

	_ "gitee.com/opengauss/openGauss-connector-go-pq"

	"github.com/UnderTreeTech/waterdrop/pkg/database/sql"
	"github.com/UnderTreeTech/waterdrop/pkg/log"
)

// NewKingbase 人大金仓DB
func NewKingbase(c *sql.Config) (db *sql.DB) {
	if c.QueryTimeout == 0 || c.ExecTimeout == 0 || c.TranTimeout == 0 {
		panic("kingbase must be set query/execute/transaction timeout")
	}
	c.WithParseDSN(drivers.ParseDSNAddr)

	db, err := sql.Open(c)
	if err != nil {
		log.Errorf("open kingbase error", log.String("error", err.Error()))
		panic(err)
	}
	return
}

// NewDm 达梦DB
func NewDm(c *sql.Config) (db *sql.DB) {
	if c.QueryTimeout == 0 || c.ExecTimeout == 0 || c.TranTimeout == 0 {
		panic("dm must be set query/execute/transaction timeout")
	}
	c.WithParseDSN(parseDMDSN)

	db, err := sql.Open(c)
	if err != nil {
		log.Errorf("open dm error", log.String("error", err.Error()))
		panic(err)
	}
	return
}

// parseDMDSN 解析dm dsn
func parseDMDSN(dsn string) (addr string) {
	splits := strings.Split(strings.Split(dsn, "@")[0], "/")
	addr = splits[0]
	return
}

// NewOpenGauss 高斯驱动
func NewOpenGauss(c *sql.Config) (db *sql.DB) {
	if c.QueryTimeout == 0 || c.ExecTimeout == 0 || c.TranTimeout == 0 {
		panic("opengauss must be set query/execute/transaction timeout")
	}
	db, err := sql.Open(c)
	if err != nil {
		log.Errorf("open odbc error", log.String("error", err.Error()))
		panic(err)
	}
	return
}

// NewODBC odbc通用驱动
func NewODBC(c *sql.Config) (db *sql.DB) {
	if c.QueryTimeout == 0 || c.ExecTimeout == 0 || c.TranTimeout == 0 {
		panic("odbc must be set query/execute/transaction timeout")
	}
	db, err := sql.Open(c)
	if err != nil {
		log.Errorf("open odbc error", log.String("error", err.Error()))
		panic(err)
	}
	return
}
