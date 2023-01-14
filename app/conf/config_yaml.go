package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	logger "nsq-demoset/app/_applib"
	"os"
)

type (
	mysql struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	}

	redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}

	nsq_ struct {
		Addr string `yaml:"addr"`
	}

	rsa_ struct {
		Private string `yaml:"private"`
		Public  string `yaml:"public"`
		Secret  string `yaml:"secret"`
	}

	telegram struct {
		TokenID string `yaml:"tokenID"`
		GroupID string `yaml:"groupID"`
	}
)

var (
	_c struct {
		Mysql    mysql    `yaml:"mysql"`
		Redis    redis    `yaml:"redis"`
		Nsq      nsq_     `yaml:"nsq"`
		Rsa      rsa_     `yaml:"rsa"`
		Telegram telegram `yaml:"telegram"`
	}
)

func InitYaml() {
	data, err := os.ReadFile("./conf/app-services.yaml")
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(data, &_c); err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}
}

func MysqlDNS() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		_c.Mysql.Username,
		_c.Mysql.Password,
		_c.Mysql.Host,
		_c.Mysql.Port,
		_c.Mysql.Database,
	)
}

func Redis() *redis {
	return &_c.Redis
}

func Nsq() *nsq_ {
	return &_c.Nsq
}

func RSA() *rsa_ {
	return &_c.Rsa
}

func Telegram() *telegram {
	return &_c.Telegram
}
