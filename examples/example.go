package main

import (
	"log"

	"github.com/wgyuuu/options"
)

type Config struct {
	App struct {
		Debug bool `options:"debug"`
	}

	Test struct {
		IntA int               `options:"int_a"`
		MapB map[string]string `options:"map_b"`
	}
}

func main() {
	var conf Config

	ops := options.NewOptions("./config.toml", handler())
	err := ops.Parsing(&conf)
	if err != nil {
		log.Println("err", err.Error())
	}

	log.Println(conf)
}

func handler() options.HandleGet {
	mmap := map[string]string{
		"debug": "true",
		"int_a": "2",
		"map_b": "{\"c\": \"test_C\"}",
	}

	return func(key, defVal string) string {
		value, ok := mmap[key]
		if !ok {
			return defVal
		}

		return value
	}
}
