package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"os"
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

func getCurrencyRate(code string, date string) (string, error) {
	url := fmt.Sprintf("https://www.cbr.ru/scripts/XML_daily.asp?date_req=%s", date)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("ошибка при получении данных от ЦБ РФ: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка при получении данных от ЦБ РФ: код статуса %d", resp.StatusCode)
	}

	var valCurs ValCurs
	err = xml.NewDecoder(resp.Body).Decode(&valCurs)
	if err != nil {
		return "", fmt.Errorf("ошибка при парсинге XML-ответа: %w", err)
	}

	for _, valute := range valCurs.Valutes {
		if valute.CharCode == code {
			return fmt.Sprintf("%s (%s): %s", code, valute.Name, valute.Value), nil
		}
	}

	return "", fmt.Errorf("валюта с таким кодом не найдена: %s", code)
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
		*date = time.Now().Format("2006-01-02") // Если дата не указана, используем текущую дату
	}

	result, err := getCurrencyRate(*code, *date)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(result)
}
