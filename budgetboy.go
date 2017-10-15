package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/silbinarywolf/budgetboy/thirdparty/shopspring/decimal"
)

type Product struct {
	Date              time.Time
	Name              string
	Price             decimal.Decimal
	BalanceAfterPrice decimal.Decimal
	Category          *ConfigCategory
}

type OutputCategory struct {
	Name        string
	ProductList []*Product
	Config      *ConfigCategory
}

type OutputDay struct {
	Name         string
	CategoryList []*OutputCategory
}

type OutputWeek struct {
	DayList [7]OutputDay
}

const netbankDateFormat = "02/01/2006" // 02 = %D, 01 = %M, 2006 = %Y

func main() {
	// Read Config
	config := ReadConfig()

	// Read/Process CSV
	product_list, err := ReadProducts(config, "./CSVData.csv")
	if err != nil {
		log.Fatal(err)
	}

	outputWeekMap := make(map[string]*OutputWeek)

	// Pull transaction data into history
	for _, product := range product_list {
		date := product.Date
		lastDayOfWeek := TimeEndOfWeek(date, false)

		// ie. 2017_October_1
		weekIndex := fmt.Sprintf("%d_%s_%s", lastDayOfWeek.Year(), lastDayOfWeek.Month().String(), DayOrdinal(lastDayOfWeek.Day()))
		outputWeek, ok := outputWeekMap[weekIndex]
		if !ok {
			// Create new week
			outputWeek = new(OutputWeek)
			for i := 0; i < len(outputWeek.DayList); i++ {
				day := &outputWeek.DayList[i]
				weekday := time.Weekday(i)
				day.Name = fmt.Sprintf("%s", weekday.String())
				for _, config_category := range config.CategoryList {
					category := new(OutputCategory)
					category.Name = config_category.Name
					category.Config = config_category
					day.CategoryList = append(day.CategoryList, category)
				}
			}
			outputWeekMap[weekIndex] = outputWeek
		}

		// Add product to that day of the week
		dayIndex := int(date.Weekday())
		day := &outputWeek.DayList[dayIndex]
		for _, category := range day.CategoryList {
			if category.Config != product.Category {
				continue
			}
			category.ProductList = append(category.ProductList, product)
		}
	}

	// Create output folder
	os.Mkdir("./output", 0644)

	// Output CSV
	//var fileBuffer bytes.Buffer
	for weekIndex, week := range outputWeekMap {
		var buffer bytes.Buffer
		buffer.WriteString("\"\",")
		for i, config_category := range config.CategoryList {
			if config_category.Name == "_" {
				continue
			}
			if i != 0 {
				buffer.WriteRune(',')
			}
			buffer.WriteString(fmt.Sprintf("\"%s\"", config_category.Name))
		}
		buffer.WriteRune('\n')

		total := config.EarningPerWeek
		total = total.Sub(config.RentPerWeek)

		total = total.Add(WriteDay(&buffer, week.DayList[time.Monday]))
		total = total.Add(WriteDay(&buffer, week.DayList[time.Tuesday]))
		total = total.Add(WriteDay(&buffer, week.DayList[time.Wednesday]))
		total = total.Add(WriteDay(&buffer, week.DayList[time.Thursday]))
		total = total.Add(WriteDay(&buffer, week.DayList[time.Friday]))
		total = total.Add(WriteDay(&buffer, week.DayList[time.Saturday]))
		total = total.Add(WriteDay(&buffer, week.DayList[time.Sunday]))

		buffer.WriteRune('\n')
		buffer.WriteString("\"Earning This Week\",")
		buffer.WriteString(fmt.Sprintf("\"%s\"", config.EarningPerWeek.String()))
		buffer.WriteRune('\n')
		buffer.WriteString("\"Rent This Week\",")
		buffer.WriteString(fmt.Sprintf("\"-%s\"", config.RentPerWeek.String()))
		buffer.WriteRune('\n')
		buffer.WriteString("\"Total Saved\",")
		buffer.WriteString(fmt.Sprintf("\"%s\"", total.String()))

		if config.PrintInConsole {
			fmt.Printf("\n\nFilename: %s\n", weekIndex)
			fmt.Printf("----------------------------------\n")
			fmt.Printf("%s\n", buffer.String())
			fmt.Printf("----------------------------------\n")
		}

		err := ioutil.WriteFile("./output/"+weekIndex+".csv", buffer.Bytes(), 0644)
		if err != nil {
			fmt.Printf("- File write error: %v\n", err)
		}
	}
	fmt.Printf("Done. Check your \"output\" folder.")
}

func WriteDay(buffer *bytes.Buffer, day OutputDay) decimal.Decimal {
	totalDayPrice := decimal.NewFromFloat(0)

	buffer.WriteString(fmt.Sprintf("\"%s\",", day.Name))
	for i, category := range day.CategoryList {
		if category.Name == "_" {
			continue
		}
		if i != 0 {
			buffer.WriteRune(',')
		}
		// Calculate total
		totalPrice := decimal.NewFromFloat(0)
		for _, product := range category.ProductList {
			totalPrice = totalPrice.Add(product.Price)
		}
		totalDayPrice = totalDayPrice.Add(totalPrice)
		totalPriceString := totalPrice.String()
		buffer.WriteRune('"')
		if totalPriceString != "0" {
			buffer.WriteString(totalPriceString)
		}
		buffer.WriteRune('"')
	}
	buffer.WriteRune('\n')
	return totalDayPrice
}

func ReadProducts(config *Config, filepath string) ([]*Product, error) {
	// Read CSV
	file_bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(strings.NewReader(string(file_bytes)))
	csv_records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	// Process CSV
	products := make([]*Product, 0, len(csv_records))
	line_number := 0
	for _, record := range csv_records {
		line_number++

		product := &Product{}
		product.Name = record[2]

		var date string = record[0]
		deffered_date := strings.Split(product.Name, "Value Date:")
		if len(deffered_date) > 1 {
			// Use "Value Date: " information instead if available.
			date = strings.TrimSpace(deffered_date[1])
		}
		// Product.Date
		{
			t, err := time.Parse(netbankDateFormat, date)
			if err != nil {
				log.Fatal(fmt.Sprintf("Line %d - Unable to parse date \"%s\"", line_number, date))
			}
			product.Date = t
		}
		// Product.Price
		{
			moneyString := record[1]
			if moneyString[0] == '$' {
				moneyString = moneyString[1:]
			}
			p, err := decimal.NewFromString(moneyString)
			if err != nil {
				log.Fatal(fmt.Sprintf("Line %d - Unable to parse the price %s as currency.", line_number, moneyString))
			}
			product.Price = p
		}
		// Product.BalanceAfterPrice
		{
			moneyString := record[1]
			if moneyString[0] == '$' {
				moneyString = moneyString[1:]
			}
			p, err := decimal.NewFromString(moneyString)
			if err != nil {
				log.Fatal(fmt.Sprintf("Line %d - Unable to parse the bank balance %s as currency.", line_number, moneyString))
			}
			product.BalanceAfterPrice = p
		}

		// Attach category to Product
	CategoryMatchLoop:
		for _, category := range config.CategoryList {
			for _, config_product := range category.ProductList {
				if !strings.Contains(product.Name, config_product.PartialName) {
					continue
				}
				product.Category = category
				break CategoryMatchLoop
			}
		}
		if product.Category == nil {
			product.Category = config.NoCategory
		}

		products = append(products, product)
	}

	// Print all uncategorized products
	if config.DisallowUncategorized {
		var buffer bytes.Buffer
		for _, product := range products {
			if product.Category != config.NoCategory {
				continue
			}
			buffer.WriteString(fmt.Sprintf("- %s", product.Name))
			buffer.WriteRune('\n')
		}
		if buffer.Len() > 0 {
			fmt.Printf("Uncategorized transactions found:\n")
			fmt.Printf(buffer.String())
			fmt.Printf("Stopping due to \"Disallow Uncategorized\" being set to true.\n")
			os.Exit(0)
		}
	}
	return products, nil
}
