package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var recPointMap map[string]recPoint

type receipt struct {
	ID           string `json:"id"`
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []item `json:"items"`
}

type item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type recPoint struct {
	Points int `json:"points"`
}

func getReceiptPoints(c *gin.Context) {
	id := c.Param("id")
	recPoints, ok := recPointMap[id]

	if ok {
		c.JSON(http.StatusOK, recPoints)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "No receipt found for that ID."})
}

func processReceipt(c *gin.Context) {
	var incomingReceipt receipt
	if err := c.ShouldBindJSON(&incomingReceipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "The receipt is invalid."})
		return
	}
	newId := uuid.New()
	incomingReceipt.ID = newId.String()
	err := calculatePointsForReceipt(incomingReceipt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "The receipt is invalid."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": incomingReceipt.ID})
}

func calculatePointsForReceipt(receipt receipt) error {
	pts := 0
	// One point for every alphanumeric character in the retailer name.
	pts += calculatePointsForRetailerName(receipt.Retailer)

	// 50 points if the total is a round dollar amount with no cents.
	totalBonus, err := calculateRoundTotalEvenBonus(receipt.Total)
	if err != nil {
		return err
	}
	pts += totalBonus

	//25 points if the total is a multiple of 0.25.
	total25Bonus, err := calculateRoundTotal25Bonus(receipt.Total)
	if err != nil {
		return err
	}
	pts += total25Bonus

	//5 points for every two items on the receipt.
	pts += calculateItemLengthBonus(receipt.Items)

	//If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2
	//and round up to the nearest integer. The result is the number of points earned.
	itemDescBonus, err := calculateItemDescBonus(receipt.Items)
	if err != nil {
		return err
	}
	pts += itemDescBonus

	//6 points if the day in the purchase date is odd.
	dateBonus, err := calculateOddDayBonus(receipt.PurchaseDate)
	if err != nil {
		return err
	}
	pts += dateBonus

	//10 points if the time of purchase is after 2:00pm and before 4:00pm.
	timeBonus, err := calculateTimeOfDayBonus(receipt.PurchaseTime)
	if err != nil {
		return err
	}
	pts += timeBonus
	recPointMap[receipt.ID] = recPoint{Points: pts}
	return nil
}

func calculatePointsForRetailerName(retailerName string) int {
	pts := 0

	for _, char := range retailerName {
		if regexp.MustCompile(`^[a-zA-Z0-9]$`).MatchString(string(char)) {
			pts++
		}
	}

	return pts
}

func calculateRoundTotalEvenBonus(total string) (int, error) {
	totalFloat, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return 0, err
	}
	if math.Mod(totalFloat, 1) == 0 {
		return 50, nil
	}
	return 0, nil
}

func calculateRoundTotal25Bonus(total string) (int, error) {
	totalFloat, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return 0, err
	}
	if math.Mod(totalFloat, .25) == 0 {
		return 25, nil
	}
	return 0, nil
}

func calculateItemLengthBonus(items []item) int {
	if items == nil {
		return 0
	}
	return int(math.Floor(float64(len(items)/2))) * 5
}

func calculateItemDescBonus(items []item) (int, error) {
	if items == nil {
		return 0, nil
	}

	pts := 0

	for _, item := range items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return 0, err
			}
			pts += int(math.Ceil(price * .2))
		}
	}

	return pts, nil
}

func calculateOddDayBonus(purchaseDate string) (int, error) {
	date, err := time.Parse("2006-01-02", purchaseDate)
	if err != nil {
		return 0, err
	}
	if date.Day()%2 != 0 {
		return 6, nil
	}
	return 0, nil
}

func calculateTimeOfDayBonus(purchaseTime string) (int, error) {
	startTime, _ := time.Parse("15:04", "14:00")
	endTime, _ := time.Parse("15:04", "16:00")

	parsedTime, err := time.Parse("15:04", purchaseTime)
	if err != nil {
		return 0, err
	}
	if parsedTime.Before(endTime) && parsedTime.After(startTime) {
		return 10, nil
	}
	return 0, nil
}

func main() {
	recPointMap = make(map[string]recPoint)
	router := gin.Default()
	router.GET("/receipts/:id/points", getReceiptPoints)
	router.POST("/receipts/process", processReceipt)
	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
