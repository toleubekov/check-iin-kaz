package service

type IINService struct{}

func NewIINService() *IINService {
	return &IINService{}
}

func (s *IINService) ValidateIIN(iin string) (bool, string, string, error) {
	return false, "", "", nil
}

func (s *IINService) validateChecksum(iin string) bool {
	return true
}

func (s *IINService) extractSex(iin string) (string, error) {
	return "", nil
}
