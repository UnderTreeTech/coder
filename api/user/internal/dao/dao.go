package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/UnderTreeTech/waterdrop/pkg/database/mongo"

	"github.com/UnderTreeTech/waterdrop/pkg/conf"
	"github.com/UnderTreeTech/waterdrop/pkg/database/redis"
	"github.com/UnderTreeTech/waterdrop/pkg/database/sql"
	"github.com/UnderTreeTech/waterdrop/pkg/log"

	"github.com/Masterminds/squirrel"
	"github.com/UnderTreeTech/drivers"
	"user/internal/dao/iface"
)

type Dao interface {
	Close() error
	Ping(ctx context.Context) error

	GetCollection(name string) *mongo.Collection
	Redis() *redis.Redis

	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	iface.TUser
}

// struct dao
type dao struct {
	db     *sql.DB
	driver string
	redis  *redis.Redis
	mongo  *mongo.DB
}

// New return a dao that implements interface Dao
func New() Dao {
	db, driver := NewDB()
	mongo := NewMongo()
	redis := NewRedis()
	return &dao{
		db:     db,
		driver: driver,
		redis:  redis,
		mongo:  mongo,
	}
}

// Close close backend base services
func (d *dao) Close() (err error) {
	d.db.Close()
	d.mongo.Close()
	d.redis.Close()
	return
}

// Ping ping backend base services, like db, redis, es etc.
func (d *dao) Ping(ctx context.Context) error {
	if err := d.db.Ping(ctx); err != nil {
		log.Error(ctx, "ping db fail", log.String("error", err.Error()))
		return err
	}

	if err := d.mongo.Ping(); err != nil {
		log.Error(ctx, "ping mongo fail", log.String("error", err.Error()))
		return err
	}

	if alive := d.redis.Ping(ctx); !alive {
		return errors.New("redis has gone")
	}

	return nil
}

func (d *dao) GetCollection(name string) *mongo.Collection {
	return d.mongo.GetCollection(name)
}

func (d *dao) Redis() *redis.Redis {
	return d.redis
}

// attach transaction flag to context
type txKey struct{}

func (d *dao) Begin(ctx context.Context) (context.Context, error) {
	tx, err := d.db.Begin(ctx)
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, txKey{}, tx)
	return ctx, err
}

func (d *dao) Commit(ctx context.Context) error {
	tx, err := d.GetTxFromCtx(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (d *dao) Rollback(ctx context.Context) error {
	tx, err := d.GetTxFromCtx(ctx)
	if err != nil {
		return err
	}

	return tx.Rollback()
}

func (d *dao) GetTxFromCtx(ctx context.Context) (*sql.Tx, error) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	if !ok {
		return nil, errors.New("assert tx err")
	}

	return tx, nil
}

// NewDB new db instance according by db driver
func NewDB() (*sql.DB, string) {
	config := &sql.Config{}
	if err := conf.Unmarshal("db", config); err != nil {
		panic(fmt.Sprintf("unmarshal db config fail,err msg %s", err.Error()))
	}
	log.Debugf("db config", log.Any("config", config))

	var db *sql.DB
	switch config.DriverName {
	case drivers.DBDriverType_mysql.String():
		db = sql.NewMySQL(config)
	case drivers.DBDriverType_postgres.String():
		db = sql.NewPostgres(config)
	case drivers.DBDriverType_kingbase.String():
		db = NewKingbase(config)
	case drivers.DBDriverType_dm.String():
		db = NewDm(config)
	case drivers.DBDriverType_gbase.String():
		config.DriverName = drivers.DBDriverType_opengauss.String()
		db = NewOpenGauss(config)
	case drivers.DBDriverType_vastbase.String(),
		drivers.DBDriverType_highgo.String():
		config.DriverName = drivers.DBDriverType_postgres.String()
		db = sql.NewPostgres(config)
	case drivers.DBDriverType_oceanbase.String():
		config.DriverName = drivers.DBDriverType_mysql.String()
		db = sql.NewMySQL(config)
	default:
		panic(fmt.Sprintf("unsupport db driver type:%s", config.DriverName))
	}
	return db, config.DriverName
}

// NewRedis returns redis instance
func NewRedis() *redis.Redis {
	config := &redis.Config{}
	if err := conf.Unmarshal("redis", config); err != nil {
		panic(fmt.Sprintf("unmarshal redis config fail,err msg %s", err.Error()))
	}
	log.Infof("redis config", log.Any("config", config))

	redis, err := redis.New(config)
	if err != nil {
		panic(fmt.Sprintf("new redis client fail,err msg %s", err.Error()))
	}
	return redis
}

// NewMongo return mongo instance
func NewMongo() *mongo.DB {
	cfg := &mongo.Config{}
	if err := conf.Unmarshal("mongo", cfg); err != nil {
		panic(fmt.Sprintf("unmarshal mongo config fail, err msg %s", err.Error()))
	}

	db := mongo.Open(cfg)
	return db
}

// PlaceHolder returns placeholder format by driver
func (d *dao) PlaceHolder() squirrel.PlaceholderFormat {
	switch d.driver {
	case drivers.DBDriverType_mysql.String(),
		drivers.DBDriverType_dm.String():
		return squirrel.Question
	case drivers.DBDriverType_postgres.String(),
		drivers.DBDriverType_kingbase.String(),
		drivers.DBDriverType_opengauss.String():
		return squirrel.Dollar
	case drivers.DBDriverType_mssql.String():
		return squirrel.AtP
	case drivers.DBDriverType_oracle.String():
		return squirrel.Colon
	default:
		return squirrel.Question
	}
}

func (d *dao) Analytic(build squirrel.SelectBuilder, condition map[string]interface{}) (squirrel.SelectBuilder, error) {
	// add order by
	if orderBy, ok := condition[drivers.OpAction__order_by.String()]; ok {
		if orderBy, ok := orderBy.(string); ok {
			build = build.OrderBy(orderBy)
			delete(condition, drivers.OpAction__order_by.String())
		} else {
			return build, errors.New("_orderBy type is string")
		}
	}

	// add group by
	if groupBy, ok := condition[drivers.OpAction__group_by.String()]; ok {
		if groupBy, ok := groupBy.(string); ok {
			build = build.GroupBy(groupBy)
			delete(condition, drivers.OpAction__group_by.String())
		} else {
			return build, errors.New("_groupBy type is string")
		}

		//add having condition
		if having, ok := condition[drivers.OpAction__having.String()]; ok {
			if having, ok := having.(string); ok {
				build = build.Having(having)
				delete(condition, drivers.OpAction__having.String())
			} else {
				return build, errors.New("_having type is string")
			}

		}
	}

	// add offset
	if offset, ok := condition[drivers.OpAction__offset.String()]; ok {
		if offset, ok := offset.(uint64); ok {
			build = build.Offset(offset)
			delete(condition, drivers.OpAction__offset.String())
		} else {
			return build, errors.New("_offset type is uint64")
		}

	}

	// add limit
	if limit, ok := condition[drivers.OpAction__limit.String()]; ok {
		if limit, ok := limit.(uint64); ok {
			build = build.Limit(limit)
			delete(condition, drivers.OpAction__limit.String())
		} else {
			return build, errors.New("_limit type is uint64")
		}

	}

	// add like
	if like, ok := condition[drivers.OpAction__like.String()]; ok {
		if likeCondition, ok := like.(map[string]interface{}); ok {
			build = build.Where(squirrel.Like(likeCondition))
			delete(condition, drivers.OpAction__like.String())
		} else {
			return build, errors.New("_like type need map[string]interface{}")
		}
	}

	// add not like
	if like, ok := condition[drivers.OpAction__not_like.String()]; ok {
		if likeCondition, ok := like.(map[string]interface{}); ok {
			build = build.Where(squirrel.NotLike(likeCondition))
			delete(condition, drivers.OpAction__not_like.String())
		} else {
			return build, errors.New("_notLike type need map[string]interface{}")
		}
	}

	// add not equal
	if noteq, ok := condition[drivers.OpAction__not_eq.String()]; ok {
		if noteqCondition, ok := noteq.(map[string]interface{}); ok {
			build = build.Where(squirrel.NotEq(noteqCondition))
			delete(condition, drivers.OpAction__not_eq.String())
		} else {
			return build, errors.New("_notEq type need map[string]interface{}")
		}
	}

	// add gt
	if gt, ok := condition[drivers.OpAction__gt.String()]; ok {
		if gtCond, ok := gt.(map[string]interface{}); ok {
			build = build.Where(squirrel.Gt(gtCond))
			delete(condition, drivers.OpAction__gt.String())
		}
	}

	// add GtOrEq
	if gtOrEq, ok := condition[drivers.OpAction__gte.String()]; ok {
		if gtOrEqCond, ok := gtOrEq.(map[string]interface{}); ok {
			build = build.Where(squirrel.GtOrEq(gtOrEqCond))
			delete(condition, drivers.OpAction__gte.String())
		}
	}

	// add lt
	if lt, ok := condition[drivers.OpAction__lt.String()]; ok {
		if ltCond, ok := lt.(map[string]interface{}); ok {
			build = build.Where(squirrel.Lt(ltCond))
			delete(condition, drivers.OpAction__lt.String())
		}
	}

	// add LtOrEq
	if ltOrEq, ok := condition[drivers.OpAction__lte.String()]; ok {
		if ltOrEqCond, ok := ltOrEq.(map[string]interface{}); ok {
			build = build.Where(squirrel.LtOrEq(ltOrEqCond))
			delete(condition, drivers.OpAction__lte.String())
		} else {
			return build, errors.New("_ltOrEq type need map[string]interface{}")
		}
	}

	return build.Where(condition), nil
}

func (d *dao) AnalyticUpdate(build squirrel.UpdateBuilder, condition map[string]interface{}) (squirrel.UpdateBuilder, error) {
	// add order by
	if orderBy, ok := condition[drivers.OpAction__order_by.String()]; ok {
		if orderBy, ok := orderBy.(string); ok {
			build = build.OrderBy(orderBy)
			delete(condition, drivers.OpAction__order_by.String())
		} else {
			return build, errors.New("_orderBy type is string")
		}
	}

	// add offset
	if offset, ok := condition[drivers.OpAction__offset.String()]; ok {
		if offset, ok := offset.(uint64); ok {
			build = build.Offset(offset)
			delete(condition, drivers.OpAction__offset.String())
		} else {
			return build, errors.New("_offset type is uint64")
		}

	}

	// add limit
	if limit, ok := condition[drivers.OpAction__limit.String()]; ok {
		if limit, ok := limit.(uint64); ok {
			build = build.Limit(limit)
			delete(condition, drivers.OpAction__limit.String())
		} else {
			return build, errors.New("_limit type is uint64")
		}

	}

	// add like
	if like, ok := condition[drivers.OpAction__like.String()]; ok {
		if likeCondition, ok := like.(map[string]interface{}); ok {
			build = build.Where(squirrel.Like(likeCondition))
			delete(condition, drivers.OpAction__like.String())
		} else {
			return build, errors.New("_like type need map[string]interface{}")
		}
	}

	// add not like
	if like, ok := condition[drivers.OpAction__not_like.String()]; ok {
		if likeCondition, ok := like.(map[string]interface{}); ok {
			build = build.Where(squirrel.NotLike(likeCondition))
			delete(condition, drivers.OpAction__not_like.String())
		} else {
			return build, errors.New("_notLike type need map[string]interface{}")
		}
	}

	// add not equal
	if noteq, ok := condition[drivers.OpAction__not_eq.String()]; ok {
		if noteqCondition, ok := noteq.(map[string]interface{}); ok {
			build = build.Where(squirrel.NotEq(noteqCondition))
			delete(condition, drivers.OpAction__not_eq.String())
		} else {
			return build, errors.New("_notEq type need map[string]interface{}")
		}
	}

	// add gt
	if gt, ok := condition[drivers.OpAction__gt.String()]; ok {
		if gtCond, ok := gt.(map[string]interface{}); ok {
			build = build.Where(squirrel.Gt(gtCond))
			delete(condition, drivers.OpAction__gt.String())
		}
	}

	// add GtOrEq
	if gtOrEq, ok := condition[drivers.OpAction__gte.String()]; ok {
		if gtOrEqCond, ok := gtOrEq.(map[string]interface{}); ok {
			build = build.Where(squirrel.GtOrEq(gtOrEqCond))
			delete(condition, drivers.OpAction__gte.String())
		}
	}

	// add lt
	if lt, ok := condition[drivers.OpAction__lt.String()]; ok {
		if ltCond, ok := lt.(map[string]interface{}); ok {
			build = build.Where(squirrel.Lt(ltCond))
			delete(condition, drivers.OpAction__lt.String())
		}
	}

	// add LtOrEq
	if ltOrEq, ok := condition[drivers.OpAction__lte.String()]; ok {
		if ltOrEqCond, ok := ltOrEq.(map[string]interface{}); ok {
			build = build.Where(squirrel.LtOrEq(ltOrEqCond))
			delete(condition, drivers.OpAction__lte.String())
		} else {
			return build, errors.New("_ltOrEq type need map[string]interface{}")
		}
	}

	return build.Where(condition), nil
}

func (d *dao) AnalyticDelete(build squirrel.DeleteBuilder, condition map[string]interface{}) (squirrel.DeleteBuilder, error) {
	// add order by
	if orderBy, ok := condition[drivers.OpAction__order_by.String()]; ok {
		if orderBy, ok := orderBy.(string); ok {
			build = build.OrderBy(orderBy)
			delete(condition, drivers.OpAction__order_by.String())
		} else {
			return build, errors.New("_orderBy type is string")
		}
	}

	// add offset
	if offset, ok := condition[drivers.OpAction__offset.String()]; ok {
		if offset, ok := offset.(uint64); ok {
			build = build.Offset(offset)
			delete(condition, drivers.OpAction__offset.String())
		} else {
			return build, errors.New("_offset type is uint64")
		}

	}

	// add limit
	if limit, ok := condition[drivers.OpAction__limit.String()]; ok {
		if limit, ok := limit.(uint64); ok {
			build = build.Limit(limit)
			delete(condition, drivers.OpAction__limit.String())
		} else {
			return build, errors.New("_limit type is uint64")
		}

	}

	// add like
	if like, ok := condition[drivers.OpAction__like.String()]; ok {
		if likeCondition, ok := like.(map[string]interface{}); ok {
			build = build.Where(squirrel.Like(likeCondition))
			delete(condition, drivers.OpAction__like.String())
		} else {
			return build, errors.New("_like type need map[string]interface{}")
		}
	}

	// add not like
	if like, ok := condition[drivers.OpAction__not_like.String()]; ok {
		if likeCondition, ok := like.(map[string]interface{}); ok {
			build = build.Where(squirrel.NotLike(likeCondition))
			delete(condition, drivers.OpAction__not_like.String())
		} else {
			return build, errors.New("_notLike type need map[string]interface{}")
		}
	}

	// add not equal
	if noteq, ok := condition[drivers.OpAction__not_eq.String()]; ok {
		if noteqCondition, ok := noteq.(map[string]interface{}); ok {
			build = build.Where(squirrel.NotEq(noteqCondition))
			delete(condition, drivers.OpAction__not_eq.String())
		} else {
			return build, errors.New("_notEq type need map[string]interface{}")
		}
	}

	// add gt
	if gt, ok := condition[drivers.OpAction__gt.String()]; ok {
		if gtCond, ok := gt.(map[string]interface{}); ok {
			build = build.Where(squirrel.Gt(gtCond))
			delete(condition, drivers.OpAction__gt.String())
		}
	}

	// add GtOrEq
	if gtOrEq, ok := condition[drivers.OpAction__gte.String()]; ok {
		if gtOrEqCond, ok := gtOrEq.(map[string]interface{}); ok {
			build = build.Where(squirrel.GtOrEq(gtOrEqCond))
			delete(condition, drivers.OpAction__gte.String())
		}
	}

	// add lt
	if lt, ok := condition[drivers.OpAction__lt.String()]; ok {
		if ltCond, ok := lt.(map[string]interface{}); ok {
			build = build.Where(squirrel.Lt(ltCond))
			delete(condition, drivers.OpAction__lt.String())
		}
	}

	// add LtOrEq
	if ltOrEq, ok := condition[drivers.OpAction__lte.String()]; ok {
		if ltOrEqCond, ok := ltOrEq.(map[string]interface{}); ok {
			build = build.Where(squirrel.LtOrEq(ltOrEqCond))
			delete(condition, drivers.OpAction__lte.String())
		} else {
			return build, errors.New("_ltOrEq type need map[string]interface{}")
		}
	}

	return build.Where(condition), nil
}
