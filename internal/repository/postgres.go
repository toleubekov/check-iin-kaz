package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/toleubekov/kaspiCheckIIN/internal/model"
)

/// migrate create -ext sql -dir ./schema -seq init

type PersonRepository struct {
	db *sqlx.DB
}

func NewPersonRepository(db *sqlx.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

func (r *PersonRepository) Create(person *model.Person) error {
	return nil
}

func (r *PersonRepository) GetByIIN(iin string) (*model.Person, error) {
	var person model.Person
	return &person, nil
}

func (r *PersonRepository) FindByNamePart(namePart string) ([]model.Person, error) {
	var people []model.Person
	return people, nil
}

func InitDB(connectionString string) (*sqlx.DB, error) {
	db := &sqlx.DB{}
	return db, nil
}
