package models

import (
	"interview-challenge-backend/database"
	"strconv"
)



type RentalUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func GetUserByID(userID int) (RentalUser, error) {
	// Prepare the query statement
	query := `
		SELECT id, first_name, last_name
		FROM users
		WHERE id = ` + strconv.Itoa(userID)
	
	// Execute the query
	row := database.DB.QueryRow(query)
	
	// Create a User struct to store the retrieved user data
	var user RentalUser
	
	// Scan the query result into the User struct
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName)
	if err != nil {
		return RentalUser{}, err
	}
	
	return user, nil
}