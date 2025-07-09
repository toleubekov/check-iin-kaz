package iin

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// IINInfo содержит информацию, извлеченную из ИИН
type IINInfo struct {
	Valid       bool   `json:"valid"`
	Sex         string `json:"sex,omitempty"`           // "male" или "female"
	DateOfBirth string `json:"date_of_birth,omitempty"` // формат DD.MM.YYYY
	Century     int    `json:"century,omitempty"`       // век рождения (18, 19, 20)
	RegionCode  int    `json:"region_code,omitempty"`   // региональный код (1000-9999)
}

// Validate проверяет корректность ИИН и извлекает информацию
func Validate(iin string) (*IINInfo, error) {
	info := &IINInfo{}

	// Проверка длины
	if len(iin) != 12 {
		return info, errors.New("длина ИИН должна быть равна 12 символам")
	}

	// Проверка что все символы - цифры
	for _, c := range iin {
		if c < '0' || c > '9' {
			return info, errors.New("ИИН должен состоять только из цифр")
		}
	}

	// Проверка контрольной суммы
	if !validateChecksum(iin) {
		return info, errors.New("некорректная контрольная сумма ИИН")
	}

	// Извлечение даты рождения
	dateOfBirth, century, err := extractDateOfBirth(iin)
	if err != nil {
		return info, err
	}

	// Извлечение пола
	sex, err := extractSex(iin)
	if err != nil {
		return info, err
	}

	// Извлечение регионального кода
	regionCode, _ := strconv.Atoi(iin[7:11])

	info.Valid = true
	info.Sex = sex
	info.DateOfBirth = dateOfBirth
	info.Century = century
	info.RegionCode = regionCode

	return info, nil
}

// IsValid проверяет только корректность ИИН без извлечения информации
func IsValid(iin string) bool {
	info, err := Validate(iin)
	return err == nil && info.Valid
}

// ExtractSex извлекает пол из ИИН
func ExtractSex(iin string) (string, error) {
	if len(iin) != 12 {
		return "", errors.New("некорректная длина ИИН")
	}

	return extractSex(iin)
}

// ExtractDateOfBirth извлекает дату рождения из ИИН
func ExtractDateOfBirth(iin string) (string, error) {
	if len(iin) != 12 {
		return "", errors.New("некорректная длина ИИН")
	}

	dateOfBirth, _, err := extractDateOfBirth(iin)
	return dateOfBirth, err
}

// ValidateAndExtract - удобная функция для быстрого использования
// Возвращает: valid, sex, dateOfBirth, error (совместимо с текущим API)
func ValidateAndExtract(iin string) (bool, string, string, error) {
	info, err := Validate(iin)
	if err != nil {
		return false, "", "", err
	}

	return info.Valid, info.Sex, info.DateOfBirth, nil
}

// validateChecksum проверяет контрольную сумму ИИН
func validateChecksum(iin string) bool {
	weights1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

	sum := 0
	for i := 0; i < 11; i++ {
		digit, _ := strconv.Atoi(string(iin[i]))
		sum += digit * weights1[i]
	}

	controlDigit := sum % 11

	if controlDigit == 10 {
		weights2 := []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 1, 2}

		sum = 0
		for i := 0; i < 11; i++ {
			digit, _ := strconv.Atoi(string(iin[i]))
			sum += digit * weights2[i]
		}

		controlDigit = sum % 11
		if controlDigit == 10 {
			return false
		}
	}

	lastDigit, _ := strconv.Atoi(string(iin[11]))
	return controlDigit == lastDigit
}

// extractDateOfBirth извлекает дату рождения из ИИН
func extractDateOfBirth(iin string) (string, int, error) {
	yearTwoDigits, _ := strconv.Atoi(iin[:2])
	month, _ := strconv.Atoi(iin[2:4])
	day, _ := strconv.Atoi(iin[4:6])

	centurySexDigit, _ := strconv.Atoi(string(iin[6]))

	var baseYear int
	var century int
	switch centurySexDigit {
	case 1, 2:
		baseYear = 1800 // 1800-1899
		century = 19    // XIX век
	case 3, 4:
		baseYear = 1900 // 1900-1999
		century = 20    // XX век
	case 5, 6:
		baseYear = 2000 // 2000-2099
		century = 21    // XXI век
	default:
		return "", 0, errors.New("неверная цифра века/пола")
	}

	fullYear := baseYear + yearTwoDigits

	// Проверка месяца
	if month < 1 || month > 12 {
		return "", 0, errors.New("неверный месяц рождения")
	}

	// Проверка дня
	maxDays := 31
	if month == 4 || month == 6 || month == 9 || month == 11 {
		maxDays = 30
	} else if month == 2 {
		if (fullYear%4 == 0 && fullYear%100 != 0) || fullYear%400 == 0 {
			maxDays = 29
		} else {
			maxDays = 28
		}
	}

	if day < 1 || day > maxDays {
		return "", 0, errors.New("неверный день рождения")
	}

	// Проверка что дата не в будущем
	now := time.Now()
	birthDate := time.Date(fullYear, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if birthDate.After(now) {
		return "", 0, errors.New("дата рождения не может быть в будущем")
	}

	return fmt.Sprintf("%02d.%02d.%04d", day, month, fullYear), century, nil
}

// extractSex извлекает пол из ИИН
func extractSex(iin string) (string, error) {
	centurySexDigit, _ := strconv.Atoi(string(iin[6]))

	switch centurySexDigit {
	case 1, 3, 5:
		return "male", nil
	case 2, 4, 6:
		return "female", nil
	default:
		return "", errors.New("недопустимая цифра для определения пола")
	}
}
