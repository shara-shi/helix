package database

import (
	log "github.com/Sirupsen/logrus"
	"github.com/shara/helix/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func Open(dns string) (*DB, error) {
	// 配置Mysql数据库
	// 使用默认gorm的MaxOpenConns
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dns, // DSN data source name
	}), &gorm.Config{})

	if err != nil {
		log.Error("database open error", "dns", dns, "error", err.Error())
		return nil, err
	}

	instance := &DB{
		db: db,
	}
	return instance, err
}

// OpenDB的时候，用gorm默认max open conns ，这里增加方法可以随时设置数据库MaxOpenConns
func (instance *DB) SetMaxOpenConns(maxConns int) {
	_i, err := instance.db.DB()
	if err != nil {
		log.Error("set max open conns error", "error", err.Error())
	}
	_i.SetMaxOpenConns(maxConns)
}

func (instance *DB) AutoMigration(value interface{}) {
	instance.db.AutoMigrate(value)
}

// 以下业务层使用-------------------------------------------

func (instance *DB) GetRowByID(value interface{}, id uint64) {
	instance.db.First(value, id)
}

func (instance *DB) GetRowByStringID(value interface{}, id string) {
	instance.db.First(value, "id = ?", id)
}

func (instance *DB) GetRowsByConditions(value interface{}, conditions map[string]string) int64 {
	tx := instance.db.Where(conditions).Find(value)
	return tx.RowsAffected
}

func (instance *DB) CreateRow(value interface{}) int64 {
	tx := instance.db.Create(value)
	return tx.RowsAffected
}

func (instance *DB) DeleteRow(value interface{}) int64 {
	tx := instance.db.Delete(value)
	return tx.RowsAffected
}

func (instance *DB) UpdateRowById(value interface{}, id uint64) int64 {
	tx := instance.db.Where("id=?", id).Updates(value)
	return tx.RowsAffected
}

func (instance *DB) UpdateRowByStringId(value interface{}, id string) int64 {
	tx := instance.db.Where("id=?", id).Updates(value)
	return tx.RowsAffected
}

func (instance *DB) Save(value interface{}) int64 {
	tx := instance.db.Save(value)
	return tx.RowsAffected
}

func (instance *DB) Close() {
	sqldb, err := instance.db.DB()
	if err != nil {
		log.Error("Can't Close DB %s", config.DATABASE_DNS)
	}
	sqldb.Close()
}
