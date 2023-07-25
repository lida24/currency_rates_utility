package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	XMLName  xml.Name `xml:"Valute"`
	NumCode  string   `xml:"NumCode"`
	CharCode string   `xml:"CharCode"`
	Nominal  int      `xml:"Nominal"`
	Name     string   `xml:"Name"`
	Value    string   `xml:"Value"`
}

func getCurrencyRate(code string, date string) (float64, error) {
	url := fmt.Sprintf("https://www.cbr.ru/scripts/XML_daily.asp?date_req=%s", date)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("ошибка при получении данных от ЦБ РФ: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("ошибка при получении данных от ЦБ РФ: код статуса %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("ошибка при чтении данных: %w", err)
	}

	var valCurs ValCurs
	err = xml.Unmarshal(body, &valCurs)
	if err != nil {
		return 0, fmt.Errorf("ошибка при парсинге XML-ответа: %w", err)
	}

	for _, valute := range valCurs.Valutes {
		if valute.CharCode == code {
			value, err := replaceComma(valute.Value)
			if err != nil {
				return 0, fmt.Errorf("ошибка при преобразовании значения: %w", err)
			}
			return value, nil
		}
	}

	return 0, fmt.Errorf("валюта с таким кодом не найдена: %s", code)
}

func replaceComma(value string) (float64, error) {
	replacedValue := value

	if replacedValue == "" {
		return 0, fmt.Errorf("некорректное значение")
	}

	replacedValue = replaceDot(replacedValue)

	if replacedValue[len(replacedValue)-3:] == ",00" {
		replacedValue = replacedValue[:len(replacedValue)-3]
	}

	val, err := strconv.ParseFloat(replacedValue, 64)
	if err != nil {
		return 0, fmt.Errorf("ошибка при преобразовании значения: %w", err)
	}

	return val, nil
}

func replaceDot(value string) string {
	return strings.ReplaceAll(value, ",", ".")
}

func main() {
	code := flag.String("code", "", "Код валюты в формате ISO 4217")
	date := flag.String("date", "", "Дата в формате YYYY-MM-DD")
	flag.Parse()

	if *code == "" {
		fmt.Println("Необходимо указать код валюты")
		os.Exit(1)
	}

	if *date == "" {
		*date = time.Now().Format("2006-01-02")
	}

	result, err := getCurrencyRate(*code, *date)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(result)
}
