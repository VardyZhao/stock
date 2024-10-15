package db

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"stock/config"
	"stock/db/migration"
	"time"
)

type Manager struct {
	connections map[string]*gorm.DB
}

var Conn *Manager

type DatabaseConfig struct {
	Name     string `mapstructure:"name"`
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	DbName   string `mapstructure:"dbname"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Migrate  bool   `mapstructure:"migrate"`
	Charset  string `mapstructure:"charset"`
	MaxIdle  int    `mapstructure:"max_idle"`
	MaxOpen  int    `mapstructure:"max_open"`
	MaxLife  int    `mapstructure:"max_life"`
}

func Load() {
	var dbConfigs []DatabaseConfig
	if err := config.Conf.UnmarshalKey("database", &dbConfigs); err != nil {
		log.Fatalf("Error unmarshaling databases config: %v", err)
	}
	Conn = &Manager{
		connections: make(map[string]*gorm.DB),
	}

	for _, dbConfig := range dbConfigs {
		db, err := connect(dbConfig)
		if err != nil {
			log.Fatalf("Failed to connect to database %s: %v", dbConfig.Name, err)
		}
		Conn.connections[dbConfig.Name] = db
		if dbConfig.Migrate {
			migration.Run(dbConfig.Name, db)
		}
	}
}

func connect(cfg DatabaseConfig) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "mysql":
		dialector = ConnectMysql(cfg)
	default:
		return nil, fmt.Errorf("unsupported driver: %s", cfg.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 获取底层的 sql.DB 对象
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLife) * time.Second)

	return db, nil
}

// GetDB 获取指定名称的数据库连接
func (m *Manager) GetDB(name string) *gorm.DB {
	db, exists := m.connections[name]
	if !exists {
		log.Fatalf("database [%s] not found", name)
		return nil
	}
	return db
}

func (m *Manager) CloseAll() {
	for name, db := range m.connections {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Error getting sql.DB for [%s]: %v", name, err)
			continue
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database [%s]: %v", name, err)
		} else {
			log.Printf("Closed database [%s]", name)
		}
	}
}

func (m *Manager) New(cfg DatabaseConfig) {
	if _, exists := m.connections[cfg.Name]; exists {
		log.Fatalf("database [%s] already exists", cfg.Name)
		return
	}

	db, err := connect(cfg)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	m.connections[cfg.Name] = db
}

func (m *Manager) Close(name string) {
	db, exists := m.connections[name]
	if !exists {
		log.Fatalf("database [%s] does not exist", name)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Fatalf(err.Error())
		return
	}

	delete(m.connections, name)
}
