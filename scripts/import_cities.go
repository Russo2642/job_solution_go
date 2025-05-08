package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"job_solition/internal/config"
	"job_solition/internal/db"
)

type City struct {
	Name    string
	Region  string
	Country string
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Использование: %s <путь_к_csv_файлу>", os.Args[0])
	}

	csvPath := os.Args[1]

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	postgres, err := db.NewPostgreSQL(cfg.PostgreSQL)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer postgres.Close()

	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatalf("Ошибка открытия файла: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.LazyQuotes = true

	header, err := reader.Read()
	if err != nil {
		log.Fatalf("Ошибка чтения заголовка: %v", err)
	}

	if len(header) < 3 {
		log.Fatalf("Неверный формат CSV. Требуются столбцы: name, region, country")
	}

	nameIdx, regionIdx, countryIdx := -1, -1, -1
	for i, col := range header {
		colName := strings.ToLower(strings.TrimSpace(col))
		switch colName {
		case "name", "город":
			nameIdx = i
		case "region", "регион", "область":
			regionIdx = i
		case "country", "страна":
			countryIdx = i
		}
	}

	if nameIdx == -1 || regionIdx == -1 || countryIdx == -1 {
		log.Fatalf("Не найдены обязательные столбцы: name, region, country")
	}

	var cities []City
	lineNum := 1

	for {
		lineNum++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Ошибка чтения строки %d: %v", lineNum, err)
			continue
		}

		if len(record) <= max(nameIdx, regionIdx, countryIdx) {
			log.Printf("Пропуск строки %d: недостаточно полей", lineNum)
			continue
		}

		name := strings.TrimSpace(record[nameIdx])
		region := strings.TrimSpace(record[regionIdx])
		country := strings.TrimSpace(record[countryIdx])

		if name == "" || country == "" {
			log.Printf("Пропуск строки %d: пустое название города или страны", lineNum)
			continue
		}

		if region == "" {
			region = name
		}

		cities = append(cities, City{
			Name:    name,
			Region:  region,
			Country: country,
		})
	}

	fmt.Printf("Найдено %d городов для импорта\n", len(cities))

	tx, err := postgres.Begin()
	if err != nil {
		log.Fatalf("Ошибка начала транзакции: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO cities (name, region, country)
		VALUES ($1, $2, $3)
		ON CONFLICT (name, region, country) DO NOTHING
	`)
	if err != nil {
		log.Fatalf("Ошибка подготовки запроса: %v", err)
	}
	defer stmt.Close()

	for _, city := range cities {
		_, err := stmt.Exec(city.Name, city.Region, city.Country)
		if err != nil {
			log.Printf("Ошибка вставки города %s: %v", city.Name, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Ошибка при коммите транзакции: %v", err)
	}

	fmt.Println("Импорт городов успешно завершен!")
}

func max(nums ...int) int {
	result := nums[0]
	for _, num := range nums[1:] {
		if num > result {
			result = num
		}
	}
	return result
}
