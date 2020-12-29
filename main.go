package main

import (
	"bufio"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var config Config

const (
	KEY_MODE     = "mode"
	KEY_USER     = "user"
	KEY_PASSWORD = "password"
	KEY_TOKEN    = "token"
	CONFIG_FILE  = "config.json"
)

func main() {
	err := loadConfig()
	if err != nil {
		panic(err)
	}

	configPrint(KEY_MODE, KEY_USER, KEY_PASSWORD, KEY_TOKEN)
	scanner()
}

func loadConfig() (err error) {
	viper.SetConfigFile(CONFIG_FILE)
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error %v ", err)
		return
	}

	config = *configParse()

	go func() {
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			err = viper.ReadInConfig()
			if err != nil {
				fmt.Printf("Error %v ", err)
				return
			}

			config = *configParse()
		})
	}()

	return
}

func scanner() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Type config name : ")

	for scanner.Scan() {
		in := scanner.Text()
		var val string

		switch in {
		case KEY_MODE:
			val = config.Mode
			break
		case KEY_USER:
			val = config.User
			break
		case KEY_PASSWORD:
			val = config.Password
			break
		case KEY_TOKEN:
			val = config.AuthToken
			break
		default:
			break
		}

		if len(val) == 0 {
			fmt.Println("Not Found")
		} else {
			configPrint(in)
		}

		fmt.Print("Type config name : ")
	}
}

type Config struct {
	Mode      string
	User      string
	Password  string
	AuthToken string
}

func configParse() *Config {
	mode := viper.GetString(KEY_MODE)
	user := viper.GetString(KEY_USER)
	password := viper.GetString(KEY_PASSWORD)
	authToken := viper.GetString(KEY_TOKEN)

	c := new(Config)
	c.Mode = mode
	c.User = user
	c.Password = password
	c.AuthToken = authToken
	return c
}

func configPrint(key ...string) {
	colorReset := "\033[0m"
	//colorCyan := "\033[36m"
	colorRed := "\033[31m"
	//colorGreen := "\033[32m"
	//colorYellow := "\033[33m"
	//colorBlue := "\033[34m"
	//colorPurple := "\033[35m"
	//colorWhite := "\033[37m"

	msg := fmt.Sprintf("Config { mode:%s , user:%s , password:%s , token:%s } \n", config.Mode, config.User, config.Password, config.AuthToken)
	for i := range key {
		msg = strings.Replace(msg, key[i], colorRed+key[i]+colorReset, 1)
	}

	fmt.Println(msg)
}
