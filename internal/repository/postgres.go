package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/toleubekov/kaspiCheckIIN/internal/model"
)

/// docker run --name=postgres -e POSTGRES_PASSWORD='qwerty' -p 5432:5432 -d --rm postgres
/// migrate create -ext sql -dir ./schema -seq init

// migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' up
// docker ps for getting container name what we created right now
// docker exec -it CONTAINER_ID_RIGHT_THERE /bin/bash
/// psql -U postgres
// \d

// migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' down

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
