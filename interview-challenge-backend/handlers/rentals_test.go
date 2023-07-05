package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	// Import your package here
)

func TestGetRental(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new request with the specific ID parameter
	req := httptest.NewRequest(http.MethodGet, "/rentals/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Call the handler function
	err := GetRental(c)

	// Assert that the status code is OK (200)
	assert.Equal(t, http.StatusOK, rec.Code)
	// Assert that the error is nil
	assert.Nil(t, err)

	// Add additional assertions for the response body or other expected behaviors
}

func TestListRentals(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new request with the specific query parameters
	req := httptest.NewRequest(http.MethodGet, "/rentals?price_min=9000&price_max=75000&limit=3&offset=6&sort=price", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler function
	err := ListRentals(c)

	// Assert that the status code is OK (200)
	assert.Equal(t, http.StatusOK, rec.Code)
	// Assert that the error is nil
	assert.Nil(t, err)

	// Add additional assertions for the response body or other expected behaviors
}

func TestGetSortColumn(t *testing.T) {
	// Call the getSortColumn function with different sort strings
	sortColumn := getSortColumn("name")
	assert.Equal(t, "name", sortColumn)

	sortColumn = getSortColumn("vehicle_year")
	assert.Equal(t, "vehicle_year", sortColumn)

	sortColumn = getSortColumn("price_per_day")
	assert.Equal(t, "price_per_day", sortColumn)

	sortColumn = getSortColumn("created")
	assert.Equal(t, "created", sortColumn)

	sortColumn = getSortColumn("updated")
	assert.Equal(t, "updated", sortColumn)

	sortColumn = getSortColumn("unknown")
	assert.Equal(t, "created", sortColumn) // Default value

	sortColumn = getSortColumn("") // Empty string
	assert.Equal(t, "created", sortColumn) // Default value
}

func TestListRentalsWithInvalidParameters(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new request with invalid query parameters
	req := httptest.NewRequest(http.MethodGet, "/rentals?price_min=abc&price_max=xyz&limit=invalid&offset=invalid&sort=unknown", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler function
	err := ListRentals(c)

	// Assert that the status code is Internal Server Error (500)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	// Assert that the error is not nil
	assert.NotNil(t, err)

	// Add additional assertions for the response body or other expected behaviors
}