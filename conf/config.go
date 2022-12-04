package conf

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/infraboard/mcube/logger/zap"
)

// conf pkg 的全局变量
// 全局配置对象
var global *Config

func C() *Config {
	if global == nil {
		panic("config required")
	}
	return global
}

func SetGlobalConfig(conf *Config) {
	global = conf
}

func NewDefaultConfig() *Config {
	return &Config{
		App:   NewDefaultApp(),
		MySQL: NewDefaultMysql(),
		Log:   NewDefaultLog(),
	}
}

type Config struct {
	App   *app
	MySQL *mysql
	Log   *log
}

// 配置通过对象来进行映射的
// 我们定义的是，配置对象的数据结构

// 应用程序本身的一些配置
type app struct {
	// restful-api
	Name string
	// 127.0.0.1, 0.0.0.0
	Host string `toml:"host"`
	// 8080, 8050
	Port string `toml:"port"`
	// 比较敏感的数据，入库的是加密后的内容，加密的秘钥就是该配置
	Key string `toml:"key"`
}

func NewDefaultApp() *app {
	return &app{
		Name: "restful-api",
		Host: "127.0.0.1",
		Port: "8050",
		Key:  "default app key",
	}
}

// mysql，数据库配置
type mysql struct {
	Host        string `toml:"host"` //toml跟json一样
	Port        string `toml:"port"`
	UserName    string `toml:"username"`
	Password    string `toml:"password"`
	Database    string `toml:"database"`
	MaxOpenConn int    `toml:"max_open_conn"`
	MaxIdleConn int    `toml:"max_idle_conn"`
	//单位是秒
	MaxLifeTime int `toml:"max_life_time"`
	//单位是秒
	MaxIdleTime int `toml:"max_idle_time"`
	//互斥锁
	lock sync.Mutex
}

func NewDefaultMysql() *mysql {
	return &mysql{
		Host:        "192.168.15.203",
		Port:        "3306",
		UserName:    "jack",
		Password:    "123456",
		MaxOpenConn: 100,
		MaxIdleConn: 20,
		//10分钟
		MaxLifeTime: 10 * 60 * 60,
		MaxIdleTime: 5 * 60 * 60,
	}
}

// 利用MySQL的配置，构建全局MySQL单例链接
var (
	db *sql.DB
)

// getDBConn use to get db connection pool
func (m *mysql) getDBConn() (*sql.DB, error) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&multiStatements=true", m.UserName, m.Password, m.Host, m.Port, m.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to mysql<%s> error, %s", dsn, err.Error())
	}
	db.SetMaxOpenConns(m.MaxOpenConn)
	db.SetMaxIdleConns(m.MaxIdleConn)
	db.SetConnMaxLifetime(time.Second * time.Duration(m.MaxLifeTime))
	// 设置db ping超时时间为5秒
	db.SetConnMaxIdleTime(time.Second * time.Duration(m.MaxIdleTime))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 回收上下文超时时间器
	defer cancel()
	// db ping，检测是否存活
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping mysql<%s> error, %s", dsn, err.Error())
	}
	return db, nil
}

// GetDB todo
func (m *mysql) GetDB() (*sql.DB, error) {
	// 加载全局数据量单例
	m.lock.Lock()
	defer m.lock.Unlock()
	if db == nil {
		conn, err := m.getDBConn()
		if err != nil {
			return nil, err
		}
		db = conn
	}
	return db, nil
}

// Log todo
type log struct {
	Level   string    `toml:"level" env:"LOG_LEVEL"`
	PathDir string    `toml:"path_dir" env:"LOG_PATH_DIR"`
	Format  LogFormat `toml:"format" env:"LOG_FORMAT"`
	To      LogTo     `toml:"to" env:"LOG_TO"`
}

func NewDefaultLog() *log {
	return &log{
		Level:  zap.DebugLevel.String(),
		To:     ToStdout,
		Format: TextFormat,
	}
}
