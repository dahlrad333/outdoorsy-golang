package handlers

import (
	"fmt"
	"interview-challenge-backend/database"
	"interview-challenge-backend/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func GetRental(c echo.Context) error {
	id := c.Param("id")

	rental, err := models.FetchRental(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Rental not found"})
	}

	return c.JSON(http.StatusOK, rental)
}

func ListRentals(c echo.Context) error {
	// Parse and validate query parameters
	priceMin, _ := strconv.Atoi(c.QueryParam("price_min"))
	priceMax, _ := strconv.Atoi(c.QueryParam("price_max"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	ids := models.ParseIntArray(c.QueryParam("ids"))
	near := models.ParseLatLng(c.QueryParam("near"))
	sort := c.QueryParam("sort")

	// Construct the SQL query
	query := "SELECT * FROM rentals WHERE true"
	args := []interface{}{}

	if len(ids) > 0 {
		query += " AND id = ANY($" + strconv.Itoa(len(args)+1) + ")"
		args = append(args, pq.Array(ids))
	}

	if priceMin > 0 {
		query += " AND price_per_day >= $" +strconv.Itoa(len(args)+1)
		args = append(args, priceMin)
	}

	if priceMax > 0 {
		query += " AND price_per_day <= $" + strconv.Itoa(len(args)+1)
		args = append(args, priceMax)
	}

	if len(near) > 0 {
		query += " AND ST_Distance(ST_MakePoint(lat, lng), ST_MakePoint($" +strconv.Itoa(len(args)+1)+ ", $" + strconv.Itoa(len(args)+2)+ ")) <= 100"
		args = append(args, near[0], near[1])
	}

	if sort != "" {
		sortable := getSortColumn(sort)
		query += fmt.Sprintf(" ORDER BY %s", sortable)
	}

	if limit > 0 {
		query += " LIMIT $" + strconv.Itoa(len(args)+1)
		args = append(args, limit)
	}

	if offset > 0 {
		query += " OFFSET $" +strconv.Itoa(len(args)+1)
 		args = append(args, offset)
	}

	// Execute the SQL query
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	defer rows.Close()


	// Prepare the response struct
	var response models.Response
	var rentals []models.Rental

	// Iterate over the rows and scan the values into Rental struct
	for rows.Next() {
		var rentalCols models.RentalColumns
		err := rows.Scan(
			&rentalCols.ID, &rentalCols.UserID, &rentalCols.Name, &rentalCols.Type, &rentalCols.Description,
			&rentalCols.Sleeps, &rentalCols.PricePerDay, &rentalCols.HomeCity, &rentalCols.HomeState, &rentalCols.HomeZip,
			&rentalCols.HomeCountry, &rentalCols.VehicleMake, &rentalCols.VehicleModel, &rentalCols.VehicleYear,
			&rentalCols.VehicleLength, &rentalCols.Created, &rentalCols.Updated, &rentalCols.Lat, &rentalCols.Lng,
			&rentalCols.PrimaryImageURL,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}

		rentalUser, err := models.GetUserByID(rentalCols.UserID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}

		rentalPrice := models.RentalPrice{
			Day: rentalCols.PricePerDay,
		}

		rentalLocation := models.RentalLocation{
			City:    rentalCols.HomeCity, 
			State:   rentalCols.HomeState,
			Zip:    rentalCols.HomeZip, 
			Country: rentalCols.HomeCountry,
			Lat:    rentalCols.Lat,
			Lng:     rentalCols.Lng,
		}

		rental := models.Rental{
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

		rentals = append(rentals, rental)
	}

	// Check for any errors during rows iteration
	if err = rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	// Set the response values
	response.Total = len(rentals)
	response.Results = rentals

	// Return the response
	return c.JSON(http.StatusOK, response)
}

func getSortColumn(sort string) string {
    // Define your column names or relevant keywords for sorting
    columnNames := []string{"name", "vehicle_year", "price_per_day", "created", "updated"}

    // Iterate over the words and check if any of them match the column names
    for _, columnName := range columnNames {
        if strings.Contains(sort, columnName) {
            return columnName
        }
    }

    // If no matching column is found, return a default column name
    return "created"
}