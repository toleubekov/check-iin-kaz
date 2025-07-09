package service

import (
	"github.com/toleubekov/check-iin-kaz/iin"
)

type IINService struct{}

func NewIINService() *IINService {
	return &IINService{}
}

// ValidateIIN оборачивает функцию из пакета iin для совместимости с существующим API
func (s *IINService) ValidateIIN(iinStr string) (bool, string, string, error) {
	return iin.ValidateAndExtract(iinStr)
}

// GetFullInfo возвращает полную информацию об ИИН
func (s *IINService) GetFullInfo(iinStr string) (*iin.IINInfo, error) {
	return iin.Validate(iinStr)
}
