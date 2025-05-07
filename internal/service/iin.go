package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type IINService struct{}

func NewIINService() *IINService {
	return &IINService{}
}

func (s *IINService) ValidateIIN(iin string) (bool, string, string, error) {

	if len(iin) != 12 {
		return false, "", "", errors.New("длина ИИН должна быть равная 12")
	}

	for _, c := range iin {
		if c < '0' || c > '9' {
			return false, "", "", errors.New("ИИН должен состоять из цифр")
		}
	}

	if !s.validateChecksum(iin) {
		return false, "", "", errors.New("некорреткный иин")
	}

	dateOfBirth, err := s.extractDateOfBirth(iin)
	if err != nil {
		return false, "", "", err
	}

	sex, err := s.extractSex(iin)
	if err != nil {
		return false, "", "", err
	}

	return true, sex, dateOfBirth, nil
}

func (s *IINService) validateChecksum(iin string) bool {
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

func (s *IINService) extractDateOfBirth(iin string) (string, error) {

	yearTwoDigits, _ := strconv.Atoi(iin[:2])
	month, _ := strconv.Atoi(iin[2:4])
	day, _ := strconv.Atoi(iin[4:6])

	centurySexDigit, _ := strconv.Atoi(string(iin[6]))

	var century int
	switch centurySexDigit {
	case 1, 2:
		century = 18
	case 3, 4:
		century = 19
	case 5, 6:
		century = 20
	default:
		return "", errors.New("")
	}

	fullYear := century*100 + yearTwoDigits

	if month < 1 || month > 12 {
		return "", errors.New("неверный месяц рождения")
	}

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
		return "", errors.New("неверный день рождения")
	}

	now := time.Now()
	birthDate := time.Date(fullYear, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if birthDate.After(now) {
		return "", errors.New("день рождения в будущем")
	}

	return fmt.Sprintf("%02d.%02d.%04d", day, month, fullYear), nil
}

func (s *IINService) extractSex(iin string) (string, error) {
	centurySexDigit, _ := strconv.Atoi(string(iin[6]))

	switch centurySexDigit {
	case 1, 3, 5:
		return "male", nil
	case 2, 4, 6:
		return "female", nil
	default:
		return "", errors.New("недопустимая цифра для определния пола")
	}
}
