package main

import (
	"bufio"
	"database/sql/driver"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/silbinarywolf/budgetboy/thirdparty/shopspring/decimal"
)

type ConfigProduct struct {
	PartialName string
}

type ConfigCategory struct {
	Name        string
	ProductList []*ConfigProduct
}

type Config struct {
	EarningPerWeek        decimal.Decimal
	RentPerWeek           decimal.Decimal
	NoCategory            *ConfigCategory
	PrintInConsole        bool
	DisallowUncategorized bool
	CategoryList          []*ConfigCategory
}

func ReadConfig() *Config {
	config := &Config{}
	config.EarningPerWeek, _ = decimal.NewFromString("")
	config.NoCategory = new(ConfigCategory)
	config.NoCategory.Name = "Uncategorized"

	file, err := os.Open("./config.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var category *ConfigCategory

	scanner := bufio.NewScanner(file)
	line_number := 0
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		line_number++
		if len(text) == 0 || text[0] == '#' {
			continue
		}

		info := strings.Split(text, ":")
		if len(info) == 1 {
			product := &ConfigProduct{
				PartialName: info[0],
			}
			category.ProductList = append(category.ProductList, product)
			continue
		}
		if len(info) == 2 {
			key := strings.TrimSpace(info[0])
			value := strings.TrimSpace(info[1])
			switch key {
			case "Category":
				if value == config.NoCategory.Name {
					// If explicitly using no category for other items, dont create twice
					continue
				}
				category = new(ConfigCategory)
				category.Name = value
				config.CategoryList = append(config.CategoryList, category)
			case "Rent Per Week":
				oldValue, _ := config.RentPerWeek.Value()
				config.RentPerWeek = ReadConfigMoney(key, value, oldValue, line_number)
			case "Disallow Uncategorized":
				if value == "true" {
					config.DisallowUncategorized = true
				}
			case "Print In Console":
				if value == "true" {
					config.PrintInConsole = true
				}
			case "No Category":
				config.NoCategory.Name = value
			case "Earning Per Week":
				oldValue, _ := config.EarningPerWeek.Value()
				config.EarningPerWeek = ReadConfigMoney(key, value, oldValue, line_number)
			default:
				log.Fatal(fmt.Sprintf("Line %d - Invalid key \"%s\".\n", line_number, key))
			}
			continue
		}
		log.Fatal(fmt.Sprintf("Line %d - Unable to process, too many ':' - %s.\n", line_number, text))
	}

	// Add uncategorized last
	config.CategoryList = append(config.CategoryList, config.NoCategory)

	return config
}

func ReadConfigMoney(key string, value string, oldValue driver.Value, line_number int) decimal.Decimal {
	if oldValue != "0" {
		log.Fatal(fmt.Sprintf("Line %d - Cannot declare \"%s\" more than once in one config file.", line_number, key))
	}
	if len(value) == 0 {
		log.Fatal(fmt.Sprintf("Line %d - Cannot have blank value for \"%s\".", line_number, key))
	}
	if value[0] == '$' {
		value = value[1:]
	}
	newValue, err := decimal.NewFromString(value)
	if err != nil {
		log.Fatal(fmt.Sprintf("Line %d - Unable to read \"%s\" from \"%s\".", line_number, value, key))
	}
	return newValue
}
