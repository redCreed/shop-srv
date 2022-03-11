package config

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Name         string       `mapstructure:"name" json:"name"`     //grpc服务名称
	Host         string       `mapstructure:"host" json:"host"`     //grpc服务主机
	Tags         []string     `mapstructure:"tags"  json:"tags"`    //grpc tag
	MysqlInfo    MysqlConfig  `mapstructure:"mysql" json:"mysql"`   //mysql 配置数据
	ConsulInfo   ConsulConfig `mapstructure:"consul" json:"consul"` //consul 配置
	JaegerConfig JaegerConfig `mapstructure:"jaeger" json:"jaeger"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}

type JaegerConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}
