package models

import (
	"database/sql"

	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type Person struct {
    Id         int    `json:"id"`
    FirstName  string `json:"first_name" binding:"required"`
    LastName   string `json:"last_name" binding:"required"`
    Email      string `json:"email" binding:"required,email"`
    IpAddress  string `json:"ip_address" binding:"required,ip"`
}

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "./names.db")
	if err != nil {
		return err
	}

	DB = db
	return nil
}

func GetPersons(count int) ([]Person, error) {

	rows, err := DB.Query("SELECT id, first_name, last_name, email, ip_address from people LIMIT " + strconv.Itoa(count))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	people := make([]Person, 0)

	for rows.Next() {
		singlePerson := Person{}
    	err = rows.Scan(&singlePerson.Id, &singlePerson.FirstName, &singlePerson.LastName, &singlePerson.Email, &singlePerson.IpAddress)
		// Use // as reference to the variable so that it changes automaticly without configuration

		if err != nil {
			return nil, err
		}

		people = append(people, singlePerson)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return people, err
}
// Get all can only have one error regarding server error.

func GetPersonById(id string) (Person, error) {

	stmt, err := DB.Prepare("SELECT id, first_name, last_name, email, ip_address from people WHERE id = ?")
	// Use DB.Prepare instead of DB.Query to make sure there is no error in accepting the id string.
	// You should do this any time you’re accepting external input. 

	if err != nil {
		return Person{}, err
	}

	person := Person{}

	sqlErr := stmt.QueryRow(id).Scan(&person.Id, &person.FirstName, &person.LastName, &person.Email, &person.IpAddress)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Person{}, nil
		}
		return Person{}, sqlErr
	}
	return person, nil
}
// Separate between not available error and server error

func AddPerson(newPerson Person) (bool, error) {

	tx, err := DB.Begin()
	// This is like creating a new version of the database, prevent if any part of the transaction fails,
	// the entire transaction can be rolled back, preventing partial updates that could lead to data corruption.
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO people (first_name, last_name, email, ip_address) VALUES (?, ?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newPerson.FirstName, newPerson.LastName, newPerson.Email, newPerson.IpAddress)

	if err != nil {
		return false, err
	}

	tx.Commit() // Commnit to the main database if there is no error

	return true, nil
}

func UpdatePerson(ourPerson Person, id int) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE people SET first_name = ?, last_name = ?, email = ?, ip_address = ? WHERE Id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(ourPerson.FirstName, ourPerson.LastName, ourPerson.Email, ourPerson.IpAddress, id)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func DeletePerson(personId int) (bool, error) {

	tx, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("DELETE from people where id = ?")

	if err != nil {
		tx.Rollback() // Rollback the transaction if preparing the statement fails
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(personId)

	if err != nil {
		tx.Rollback() // Rollback the transaction if executing the statement fails
		return false, err
	}

	tx.Commit() // Commit the transaction
    // if err != nil {
    //     return false, err
    // }

	return true, nil
}