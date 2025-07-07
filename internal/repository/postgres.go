package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/toleubekov/check-iin-kaz/internal/model"
)

type PersonRepository struct {
	db *sqlx.DB
}

func NewPersonRepository(db *sqlx.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

func (r *PersonRepository) Create(person *model.Person) error {
	query := `INSERT INTO people (name, iin, phone) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, person.Name, person.IIN, person.Phone)
	if err != nil {

		if strings.Contains(err.Error(), "duplicate key") {
			return errors.New("a person with this IIN already exists")
		}
		return fmt.Errorf("failed to create person: %w", err)
	}
	return nil
}

func (r *PersonRepository) GetByIIN(iin string) (*model.Person, error) {
	var person model.Person
	query := `SELECT name, iin, phone FROM people WHERE iin = $1`
	err := r.db.Get(&person, query, iin)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, errors.New("person not found")
		}
		return nil, fmt.Errorf("failed to get person by IIN: %w", err)
	}
	return &person, nil
}

func (r *PersonRepository) FindByNamePart(namePart string) ([]model.Person, error) {
	var people []model.Person
	query := `SELECT name, iin, phone FROM people WHERE name ILIKE $1`
	err := r.db.Select(&people, query, "%"+namePart+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to find people by name part: %w", err)
	}
	return people, nil
}

func InitDB(connectionString string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
