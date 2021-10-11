/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"my-instance/config"

	"github.com/spf13/viper"
)

var cfgFile string

var version = "v1.0.0"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "instance",
	Short: "golang test instance",
	Long: `*********************`,
    RunE: run,
    Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v string) {
	version = v
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ./config/source.toml)")

	//rootCmd.PersistentFlags().Int("log-level", 5, "setDebug=5, info=4, error=2, fatal=1, panic=0")
	//rootCmd.PersistentFlags().Bool("cache", false, "open cache")
	//rootCmd.PersistentFlags().String("api-http", "", "ip:port to http the api server(eg. 0.0.0.0:7012)")
	//rootCmd.PersistentFlags().String("mysql-database", "", "database name")
	//rootCmd.PersistentFlags().String("mysql-host", "", "e.g: 127.0.0.1")
	//rootCmd.PersistentFlags().String("mysql-port", "", "e.g: 3306")
	//rootCmd.PersistentFlags().String("mysql-user", "", "e.g: helium-admin")
	//rootCmd.PersistentFlags().String("mysql-password", "", "e.g: ******")
	//rootCmd.PersistentFlags().Bool("mysql-automigrate", false, "true/false")
	//
	//viper.BindPFlag("general.log_level", rootCmd.PersistentFlags().Lookup("log-level"))
	//viper.BindPFlag("general.cache", rootCmd.PersistentFlags().Lookup("cache"))
	//viper.BindPFlag("api.http", rootCmd.PersistentFlags().Lookup("api-http"))
	//viper.BindPFlag("mysql.database", rootCmd.PersistentFlags().Lookup("mysql-database"))
	//viper.BindPFlag("mysql.host", rootCmd.PersistentFlags().Lookup("mysql-host"))
	//viper.BindPFlag("mysql.port", rootCmd.PersistentFlags().Lookup("mysql-port"))
	//viper.BindPFlag("mysql.user", rootCmd.PersistentFlags().Lookup("mysql-user"))
	//viper.BindPFlag("mysql.password", rootCmd.PersistentFlags().Lookup("mysql-password"))
	//viper.BindPFlag("mysql.automigrate", rootCmd.PersistentFlags().Lookup("mysql-automigrate"))

	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	log.Info("扫描到的配置文件", cfgFile)
	if cfgFile != "" {
		b, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			log.WithError(err).WithField("config:", cfgFile).Fatal("error loading config file")
		}
		fmt.Printf("配置文件打印 \n %s \n",string(b))
		viper.SetConfigType("toml")
		// Use config file from the flag.
		if err = viper.ReadConfig(bytes.NewBuffer(b)); err != nil {
			log.WithError(err).WithField("config:", cfgFile).Fatal("error read config file")
		}
	} else {
		viper.SetConfigType("toml")
		viper.SetConfigName("source")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.config/source")
		viper.AddConfigPath("./config/")
		//viper.AddConfigPath("./config/")
		if err := viper.ReadInConfig(); err != nil {
			switch err.(type) {
			case viper.ConfigFileNotFoundError:
				log.Warning("no config file found")
			default:
				log.WithError(err).Fatal("read config file err, exit")
			}

		}
	}

	// If a config file is found, read it in.
	if err := viper.Unmarshal(&config.C); err != nil {
		log.WithError(err).Fatal("unmarshal config error")
	}
}
