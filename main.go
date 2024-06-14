package main

import (
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Global variables -----------------------------------------------------------
// ID to points mapping
var receiptID = map[string]int{}

// Validator
var validate = validator.New()

// ---------------------------------------------------------------------------

// Custom types --------------------------------------------------------------
type Receipt struct {
	Retailer     string `json:"retailer" validate:"required,retailerValidator"`
	PurchaseDate string `json:"purchaseDate" validate:"required,datetime=2006-01-02"`
	PurchaseTime string `json:"purchaseTime" validate:"required,datetime=15:04"`
	Items        []Item `json:"items" validate:"required,min=1,dive"`
	Total        string `json:"total" validate:"required,priceValidator"`
}

type Item struct {
	ShortDescription string `json:"shortDescription" validate:"required,shortDescriptionValidator"`
	Price            string `json:"price" validate:"required,priceValidator"`
}

// ---------------------------------------------------------------------------

// Custom validators ---------------------------------------------------------
func RetailerValidator(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^[\w\s\-&]+$`)
	return re.MatchString(fl.Field().String())
}

func PriceValidator(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^\d+\.\d{2}$`)
	return re.MatchString(fl.Field().String())
}

func ShortDescriptionValidator(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^[\w\s\-]+$`)
	return re.MatchString(fl.Field().String())
}

// ---------------------------------------------------------------------------

// Rules for calculating points ----------------------------------------------
func pointsOfAlphaNumericCount(s string) int {
	count := 0
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			count++
		}
	}
	return count
}

func pointsIfRoundDollarAmount(total string) int {
	t, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return 0
	}
	if t == float64(int(t)) {
		return 50
	}
	return 0
}

func pointsIfMultipleOf25(total string) int {
	t := strings.Split(total, ".")
	if len(t) != 2 {
		return 0
	}
	t1, err := strconv.ParseInt(t[1], 10, 64)
	if err != nil {
		return 0
	}
	if t1%25 == 0 {
		return 25
	}
	return 0
}

func pointsForEveryTwoItems(items []Item) int {
	return len(items) / 2 * 5
}

func pointsByShortDescription(items []Item) int {
	points := 0
	for _, item := range items {
		sd := strings.Trim(item.ShortDescription, " ")
		if len(sd)%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				continue
			}
			points += int(math.Ceil(price * 0.2))
		}
	}
	return points
}

func pointsIfOddDay(date string) int {
	d, err := strconv.Atoi(date[len(date)-2:])
	if err != nil {
		return 0
	}
	if d%2 == 1 {
		return 6
	}
	return 0
}

func pointsIfBetween2And4(time string) int {
	if time > "14:00:00" && time < "16:00:00" {
		return 10
	}
	return 0
}

// ---------------------------------------------------------------------------

func getID(receipt Receipt) (string, int) {
	points := 0
	points += pointsOfAlphaNumericCount(receipt.Retailer)
	points += pointsIfRoundDollarAmount(receipt.Total)
	points += pointsIfMultipleOf25(receipt.Total)
	points += pointsForEveryTwoItems(receipt.Items)
	points += pointsByShortDescription(receipt.Items)
	points += pointsIfOddDay(receipt.PurchaseDate)
	points += pointsIfBetween2And4(receipt.PurchaseTime)

	id := uuid.New().String()
	return id, points
}

// Handlers ------------------------------------------------------------------
// POST /receipts/process
func processReceipt(c *gin.Context) {
	var receipt Receipt
	err := c.BindJSON(&receipt)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "The receipt is invalid"})
		return
	}
	err = validate.Struct(receipt)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "The receipt is invalid"})
		fmt.Println(err)
		return
	}
	id, points := getID(receipt)
	receiptID[id] = points
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// GET /receipts/:id/points
func getPoints(c *gin.Context) {
	id := c.Param("id")
	points, ok := receiptID[id]
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No receipt found for that id"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"points": points})
}

// GET /receipts/points
// func getAllPoints(c *gin.Context) {
// 	c.IndentedJSON(http.StatusOK, receiptID)
// }

// ---------------------------------------------------------------------------

func main() {
	router := gin.Default()
	router.POST("/receipts/process", processReceipt)
	router.GET("/receipts/:id/points", getPoints)
	// router.GET("/receipts/points", getAllPoints)

	validate.RegisterValidation("retailerValidator", RetailerValidator)
	validate.RegisterValidation("priceValidator", PriceValidator)
	validate.RegisterValidation("shortDescriptionValidator", ShortDescriptionValidator)

	router.Run("localhost:8080")
}
