package models

import (
	"database/sql"
	"errors"
	"interview-challenge-backend/database"
	"strconv"
	"strings"
)

var db *sql.DB

type Rental struct {
	ID              int             `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Type            string          `json:"type"`
	Make            string          `json:"make"`
	Model           string          `json:"model"`
	Year            int             `json:"year"`
	Length          float64         `json:"length"`
	Sleeps          int             `json:"sleeps"`
	PrimaryImageURL string          `json:"primary_image_url"`
	Price           RentalPrice     `json:"price"`
	Location        RentalLocation  `json:"location"`
	User            RentalUser      `json:"user"`
} 

type RentalColumns struct {
	ID             int     `json:"id"`
	UserID         int     `json:"user_id"`
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	Description    string  `json:"description"`
	Sleeps         int     `json:"sleeps"`
	PricePerDay    int64   `json:"price_per_day"`
	HomeCity       string  `json:"home_city"`
	HomeState      string  `json:"home_state"`
	HomeZip        string  `json:"home_zip"`
	HomeCountry    string  `json:"home_country"`
	VehicleMake    string  `json:"vehicle_make"`
	VehicleModel   string  `json:"vehicle_model"`
	VehicleYear    int     `json:"vehicle_year"`
	VehicleLength  float64 `json:"vehicle_length"`
	Created        string  `json:"created"`
	Updated        string  `json:"updated"`
	Lat            float64 `json:"lat"`
	Lng            float64 `json:"lng"`
	PrimaryImageURL string `json:"primary_image_url"`
}

type RentalPrice struct {
	Day int64 `json:"day"`
}

type RentalLocation struct {
	City    string  `json:"city"`
	State   string  `json:"state"`
	Zip     string  `json:"zip"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}

type Response struct {
	Total   int      `json:"total"`
	Results []Rental `json:"results"`
}

func FetchRental(id string) (*Rental, error) {
	query := "SELECT * FROM rentals WHERE id = " + id
	var rentalCols RentalColumns
	row := database.DB.QueryRow(query)
	err := row.Scan(
		&rentalCols.ID, &rentalCols.UserID, &rentalCols.Name, &rentalCols.Type, &rentalCols.Description,
		&rentalCols.Sleeps, &rentalCols.PricePerDay, &rentalCols.HomeCity, &rentalCols.HomeState, &rentalCols.HomeZip,
		&rentalCols.HomeCountry, &rentalCols.VehicleMake, &rentalCols.VehicleModel, &rentalCols.VehicleYear,
		&rentalCols.VehicleLength, &rentalCols.Created, &rentalCols.Updated, &rentalCols.Lat, &rentalCols.Lng,
		&rentalCols.PrimaryImageURL,
	)
	if err != nil {
		return nil, errors.New("Internal Server Error")
	}

	userId, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	rentalUser, err := GetUserByID(userId)
	if err != nil {
		return nil, errors.New("Internal Server Error")
	}

	rentalPrice := RentalPrice{
		Day: rentalCols.PricePerDay,
	}

	rentalLocation := RentalLocation{
		City:    rentalCols.HomeCity, 
		State:   rentalCols.HomeState,
		Zip:    rentalCols.HomeZip, 
		Country: rentalCols.HomeCountry,
		Lat:    rentalCols.Lat,
		Lng:     rentalCols.Lng,
	}

	rental := &Rental{
		ID:          rentalCols.ID,
		Name:        rentalCols.Name,
		Description: rentalCols.Description,
		Type:         rentalCols.Type,        
		Make:         rentalCols.VehicleMake,         
		Model:         rentalCols.VehicleModel,
		Year:           rentalCols.VehicleYear,         
		Length:       rentalCols.VehicleLength,      
		Sleeps:          rentalCols.Sleeps,        
		PrimaryImageURL: rentalCols.PrimaryImageURL,          
		Price:          rentalPrice,
		Location:       rentalLocation,  
		User:            rentalUser,      
	}

	return rental, nil
}

func ParseIntArray(input string) []string {
	if input == "" {
		return nil
	}
	return strings.Split(input, ",")
}

func ParseLatLng(input string) []float64 {
	if input == "" {
		return nil
	}

	coordinates := strings.Split(input, ",")
	lat, err := strconv.ParseFloat(coordinates[0], 64)
	if err != nil {
		return nil
	}
	lng, err := strconv.ParseFloat(coordinates[1], 64)
	if err != nil {
		return nil
	}

	return []float64{lat, lng}
}