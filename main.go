package main

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

func main() {
	router := gin.Default()
	points := make(map[string]int)

	// Path: /receipts/{id}/points
	// Method: GET
	// Response: A JSON object containing the number of points awarded.
	router.GET("/receipts/:id/points", func(ctx *gin.Context) {
		point, exists := points[ctx.Param("id")]

		if exists {
			ctx.JSON(http.StatusOK, gin.H{"points": point})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "item not found"})
		}
	})

	// Path: /receipts/process
	// Method: POST
	// Payload: Receipt JSON
	// Response: JSON containing an id for the receipt.
	router.POST("/receipts/process", func(ctx *gin.Context) {
		var receipt Receipt

		if err := ctx.BindJSON(&receipt); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "the receipt is invalid",
			})
			return
		}

		id := uuid.New().String()

		_, exists := points[id]

		for exists {
			id = uuid.New().String()
			_, exists = points[id]
		}

		current_points := calculatePoints(receipt)
		if current_points == -1 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "the receipt is invalid",
			})
			return
		}

		points[id] = current_points

		ctx.JSON(http.StatusOK, gin.H{
			id: current_points,
		})
	})

	router.Run(":8080")
}

// Calculating points for each receipt passed into POST endpoint according to provided api.yml.
func calculatePoints(receipt Receipt) int {
	currentPoints := 0

	for _, char := range receipt.Retailer {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			currentPoints += 1
		}
	}

	total, err := strconv.ParseFloat(receipt.Total, 64)

	if err != nil {
		return -1
	}

	if math.Mod(total, 1) == 0 {
		currentPoints += 50
	}

	if math.Mod(total*4, 1) == 0 {
		currentPoints += 25
	}

	everyTwoItems := len(receipt.Items) / 2
	currentPoints += everyTwoItems * 5

	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return -1
			}
			currentPoints += int(math.Ceil(price * .2))
		}
	}

	split_date := strings.Split(receipt.PurchaseDate, "-")
	if len(split_date) != 3 {
		return -1
	}
	day, err := strconv.Atoi(split_date[2])
	if err != nil {
		return -1
	}

	if day%2 == 1 {
		currentPoints += 6
	}

	split_time := strings.Split(receipt.PurchaseTime, ":")
	if len(split_time) < 2 {
		return -1
	}
	hour, err := strconv.Atoi(split_time[0])
	if err != nil {
		return -1
	}
	minute, err := strconv.Atoi(split_time[1])
	if err != nil {
		return -1
	}

	// Assuming rule "10 points if the time of purchase is after 2:00pm and before 4:00pm." means 2:00pm and 4:00pm are non
	// inclusive of those times.
	if (hour == 14 && minute != 0) || (hour == 15) {
		currentPoints += 10
	}

	return currentPoints
}
