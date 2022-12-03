package conf

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

// mysql，数据库配置
type mysql struct {
	Host     string `toml:"host"` //toml跟json一样
	Port     string `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

// Log todo
type log struct {
	Level   string    `toml:"level" env:"LOG_LEVEL"`
	PathDir string    `toml:"path_dir" env:"LOG_PATH_DIR"`
	Format  LogFormat `toml:"format" env:"LOG_FORMAT"`
	To      LogTo     `toml:"to" env:"LOG_TO"`
}
