package config

import "github.com/spf13/viper"

type Config struct{
	Db_port               string `mapstructure:"DB_PORT"`
	Db_host               string `mapstructure:"DB_HOST"`
	Db_username           string `mapstructure:"DB_USER"`
	Db_password           string `mapstructure:"DB_PASSWORD"`
	Db_name               string `mapstructure:"DB_NAME"`
	Port                  string `mapstructure:"PORT"`
	AWS_ACCESS_KEY_ID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWS_SECRET_ACCESS_KEY string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	Region                string `mapstructure:"REGION"`
}

var envs = []string{"DB_HOST","DB_PORT","DB_USER","DB_PASSWORD","DB_NAME","PORT","AWS_ACCESS_KEY_ID","AWS_SECRET_ACCESS_KEY","REGION"}

var config *Config

func LoadConfig() (config *Config,err error) {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")

	viper.ReadInConfig()

	for _,env:=range envs{
		if err :=viper.BindEnv(env);err != nil{
			return nil,err
		}
	}
	if err:=viper.Unmarshal(&config);err != nil{
		return nil,err
	}
	return config,nil
	
}

func GetConfig() *Config{
	return config
}