package cmd

import (
	clilog "coscli/logger"
	"coscli/util"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"log"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var initSkip bool
var logPath string
var config util.Config
var param util.Param
var cmdCnt int //控制某些函数在一个命令中被调用的次数

var rootCmd = &cobra.Command{
	Use:   "coscli",
	Short: "Welcome to use coscli",
	Long:  "Welcome to use coscli!",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
	Version: util.Version,
}

func Execute() error {
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config-path", "c", "", "config file path(default is $HOME/.cos.yaml)")
	rootCmd.PersistentFlags().StringVarP(&param.SecretID, "secret-id", "i", "", "config secretId")
	rootCmd.PersistentFlags().StringVarP(&param.SecretKey, "secret-key", "k", "", "config secretKey")
	rootCmd.PersistentFlags().StringVarP(&param.SessionToken, "token", "", "", "config sessionToken")
	rootCmd.PersistentFlags().StringVarP(&param.Endpoint, "endpoint", "e", "", "config endpoint")
	rootCmd.PersistentFlags().BoolVarP(&param.Customized, "customized", "", false, "config customized")
	rootCmd.PersistentFlags().StringVarP(&param.Protocol, "protocol", "p", "", "config protocol")
	rootCmd.PersistentFlags().BoolVarP(&initSkip, "init-skip", "", false, "skip config init")
	rootCmd.PersistentFlags().StringVarP(&logPath, "log-path", "", "", "coscli log dir")
}

func initConfig() {
	// 初始化日志路径
	clilog.InitLoggerWithDir(logPath)

	home, err := homedir.Dir()
	cobra.CheckErr(err)
	viper.SetConfigType("yaml")
	firstArg := ""
	if len(os.Args) > 1 {
		firstArg = os.Args[1]
	}

	if cfgFile != "" {
		if cfgFile[0] == '~' {
			cfgFile = home + cfgFile[1:]
		}
		if !strings.HasSuffix(cfgFile, ".yaml") {
			fmt.Println("config file need end with .yaml ")
			os.Exit(1)
		}
		viper.SetConfigFile(cfgFile)
	} else {
		_, err = os.Stat(home + "/.cos.yaml")
		if os.IsNotExist(err) {
			if firstArg != "config" {
				if !initSkip {
					log.Println("Welcome to coscli!\nWhen you use coscli for the first time, you need to input some necessary information to generate the default configuration file of coscli.")
					initConfigFile(false)
					cmdCnt++
				} else {
					// 若无配置文件，则需有输入ak，sk及endpoint
					if param.SecretID == "" {
						logger.Fatalln("missing parameter SecretID")
						os.Exit(1)
					}
					if param.SecretKey == "" {
						logger.Fatalln("missing parameter SecretKey")
						os.Exit(1)
					}
					if param.Endpoint == "" {
						logger.Fatalln("missing parameter Endpoint")
						os.Exit(1)
					}
					return
				}
			} else {
				if !initSkip {
					log.Println("Welcome to coscli!\nWhen you use coscli for the first time, you need to input some necessary information to generate the default configuration file of coscli.")
					initConfigFile(false)
					cmdCnt++
				} else {
					return
				}
			}

		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".cos")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		if err := viper.UnmarshalKey("cos", &config); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if config.Base.Protocol == "" {
			config.Base.Protocol = "https"
		}
		// 若未关闭秘钥加密，则先解密秘钥
		if config.Base.DisableEncryption != "true" {
			secretKey, err := util.DecryptSecret(config.Base.SecretKey)
			if err == nil {
				config.Base.SecretKey = secretKey
			}
			secretId, err := util.DecryptSecret(config.Base.SecretID)
			if err == nil {
				config.Base.SecretID = secretId
			}
			sessionToken, err := util.DecryptSecret(config.Base.SessionToken)
			if err == nil {
				config.Base.SessionToken = sessionToken
			}
		}

	} else {
		fmt.Println(err)
		os.Exit(1)
	}
}
