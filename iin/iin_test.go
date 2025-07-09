package iin_test

import (
	"fmt"
	"log"

	"github.com/toleubekov/check-iin-kaz/iin"
)

// ExampleValidate демонстрирует основное использование функции Validate
func ExampleValidate() {
	info, err := iin.Validate("031231500126")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Валидный: %t\n", info.Valid)
	fmt.Printf("Пол: %s\n", info.Sex)
	fmt.Printf("Дата рождения: %s\n", info.DateOfBirth)
	fmt.Printf("Век: %d\n", info.Century)
	// Output:
	// Валидный: true
	// Пол: male
	// Дата рождения: 31.12.2003
	// Век: 21
}

// ExampleIsValid демонстрирует быструю проверку валидности ИИН
func ExampleIsValid() {
	// Валидный ИИН
	fmt.Println(iin.IsValid("031231500126"))

	// Невалидный ИИН
	fmt.Println(iin.IsValid("123456789012"))

	// Output:
	// true
	// false
}

// ExampleExtractSex показывает извлечение пола из ИИН
func ExampleExtractSex() {
	sex, err := iin.ExtractSex("031231500126")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Пол: %s\n", sex)
	// Output:
	// Пол: male
}

// ExampleExtractDateOfBirth показывает извлечение даты рождения
func ExampleExtractDateOfBirth() {
	dateOfBirth, err := iin.ExtractDateOfBirth("031231500126")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Дата рождения: %s\n", dateOfBirth)
	// Output:
	// Дата рождения: 31.12.2003
}

// ExampleValidateAndExtract демонстрирует совместимый API
func ExampleValidateAndExtract() {
	valid, sex, dateOfBirth, err := iin.ValidateAndExtract("031231500126")
	if err != nil {
		log.Fatal(err)
	}

	if valid {
		fmt.Printf("ИИН валидный\n")
		fmt.Printf("Пол: %s\n", sex)
		fmt.Printf("Дата рождения: %s\n", dateOfBirth)
	}
	// Output:
	// ИИН валидный
	// Пол: male
	// Дата рождения: 31.12.2003
}

// ExampleValidate_errorHandling показывает обработку ошибок
func ExampleValidate_errorHandling() {
	// Невалидный ИИН - слишком короткий
	_, err := iin.Validate("12345")
	if err != nil {
		fmt.Printf("Ошибка: %s\n", err.Error())
	}

	// Невалидный ИИН - неверные символы
	_, err = iin.Validate("12345678901a")
	if err != nil {
		fmt.Printf("Ошибка: %s\n", err.Error())
	}

	// Output:
	// Ошибка: длина ИИН должна быть равна 12 символам
	// Ошибка: ИИН должен состоять только из цифр
}

// ExampleIINInfo показывает работу со структурой IINInfo
func ExampleIINInfo() {
	info, err := iin.Validate("031231500126")
	if err != nil {
		log.Fatal(err)
	}

	// Проверяем различные поля
	if info.Valid {
		fmt.Println("ИИН прошел валидацию")

		if info.Sex == "male" {
			fmt.Println("Пол: мужской")
		}

		if info.Century == 21 {
			fmt.Println("Родился в XXI веке")
		}

		fmt.Printf("Региональный код: %d\n", info.RegionCode)
	}
	// Output:
	// ИИН прошел валидацию
	// Пол: мужской
	// Родился в XXI веке
	// Региональный код: 12
}
