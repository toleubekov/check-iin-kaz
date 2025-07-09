package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/toleubekov/check-iin-kaz/iin"
)

func main() {
	// Примеры ИИН для тестирования
	testIINs := []string{
		"031231500126", // Валидный ИИН
		"850515400786", // Валидный ИИН
		"123456789012", // Невалидный ИИН
		"03123150012",  // Короткий ИИН
	}

	fmt.Println("=== Примеры использования библиотеки check-iin-kaz ===")

	for i, testIIN := range testIINs {
		fmt.Printf("--- Пример %d: ИИН %s ---\n", i+1, testIIN)

		// Метод 1: Полная валидация с извлечением всей информации
		fmt.Println("1. Полная валидация:")
		info, err := iin.Validate(testIIN)
		if err != nil {
			fmt.Printf("   Ошибка: %v\n", err)
		} else {
			jsonData, _ := json.MarshalIndent(info, "   ", "  ")
			fmt.Printf("   Результат: %s\n", jsonData)
		}

		// Метод 2: Быстрая проверка валидности
		fmt.Println("2. Быстрая проверка:")
		fmt.Printf("   Валидный: %t\n", iin.IsValid(testIIN))

		// Метод 3: Совместимость с текущим API
		fmt.Println("3. Совместимый API:")
		valid, sex, dateOfBirth, err := iin.ValidateAndExtract(testIIN)
		if err != nil {
			fmt.Printf("   Ошибка: %v\n", err)
		} else {
			fmt.Printf("   Валидный: %t, Пол: %s, Дата рождения: %s\n", valid, sex, dateOfBirth)
		}

		// Метод 4: Извлечение отдельных компонентов
		if iin.IsValid(testIIN) {
			fmt.Println("4. Извлечение отдельных данных:")

			sex, err := iin.ExtractSex(testIIN)
			if err != nil {
				fmt.Printf("   Ошибка извлечения пола: %v\n", err)
			} else {
				fmt.Printf("   Пол: %s\n", sex)
			}

			dateOfBirth, err := iin.ExtractDateOfBirth(testIIN)
			if err != nil {
				fmt.Printf("   Ошибка извлечения даты: %v\n", err)
			} else {
				fmt.Printf("   Дата рождения: %s\n", dateOfBirth)
			}
		}

		fmt.Println()
	}

	// Пример использования в реальном приложении
	fmt.Println("=== Пример интеграции в приложение ===")
	exampleIntegration()
}

// exampleIntegration показывает как использовать библиотеку в реальном приложении
func exampleIntegration() {
	userIIN := "031231500123" // ИИН от пользователя

	// Проверяем и извлекаем информацию
	info, err := iin.Validate(userIIN)
	if err != nil {
		log.Printf("Ошибка валидации ИИН: %v", err)
		return
	}

	if !info.Valid {
		log.Println("ИИН не прошел валидацию")
		return
	}

	// Используем извлеченную информацию
	fmt.Printf("Пользователь: %s, %s\n", info.Sex, info.DateOfBirth)

	// Дополнительная логика на основе данных ИИН
	if info.Century == 20 {
		fmt.Println("Пользователь родился в XXI веке")
	}

	if info.RegionCode >= 1000 && info.RegionCode <= 1999 {
		fmt.Println("Региональный код указывает на центральные области")
	}

	// Можно использовать в HTTP handler'е
	handleUserRegistration(info)
}

// handleUserRegistration пример использования в HTTP handler
func handleUserRegistration(info *iin.IINInfo) {
	fmt.Println("Регистрация пользователя с валидным ИИН:")
	fmt.Printf("- Пол: %s\n", info.Sex)
	fmt.Printf("- Дата рождения: %s\n", info.DateOfBirth)
	fmt.Printf("- Век: %d\n", info.Century)
	fmt.Printf("- Региональный код: %d\n", info.RegionCode)
}
