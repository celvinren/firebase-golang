package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	models "restaurant_golang/models"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

// GetStoreList doc
func GetStoreList(c *gin.Context) {
	//TODO:
	organizationId := c.Param("organizationId")
	log.Println(organizationId)

	client, err := InitFireStore()
	if err != nil {
		log.Println(err)
	}

	iter := client.Collection("store").Where("OrganizationID", "==", organizationId).Documents(context.Background())
	var storeList []models.Store
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": 400,
				"msg":    "get store list error: " + err.Error(),
				"data":   nil,
			})
		}
		log.Println(doc.Data())
		log.Println(doc.Data()["StoreName"].(string))
		jsonStr, err := json.Marshal(doc.Data())
		if err != nil {
			log.Println("convert map to json string failed")
		}
		storeList = append(storeList, models.StoreFromJSON(string(jsonStr)))
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "get store list successfully",
		"data":   storeList,
	})
}

// GetStore doc
func GetStore(c *gin.Context) {
	//TODO:

}

// CreateStore doc
func CreateStore(c *gin.Context) {
	//bind request json to object
	var storeDTO models.StoreDTO
	bindErr := c.BindJSON(&storeDTO)
	if bindErr != nil {
		log.Println("Bind json data error:", bindErr)
	}

	//extract logo base64 string and remove in the storeDTO so create store without base64 string
	logoStr := storeDTO.Logo
	storeDTO.Logo = ""

	//create Store
	client, err := InitFireStore()
	if err != nil {
		log.Println("fail to create firebase client: " + err.Error())
	}

	s, _, errCreate := client.Collection("store").Add(context.Background(), storeDTO)
	if errCreate != nil {
		log.Println("An error has occurred when update organization:" + errCreate.Error())
	}
	log.Println(s.ID)

	//save base64 string to cloud storage and return url string
	path := storeDTO.OrganizationID + "/" + s.ID + "/logo/logo.jpg"
	urlStr := SaveBase64ToCloudStore(logoStr, path, c)

	//update store logo url
	storeDTO.Logo = urlStr
	_, err = client.Collection("store").Doc(s.ID).Set(context.Background(), storeDTO)
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}

	log.Println(urlStr)
}

// UpdateStore doc
func UpdateStore(c *gin.Context) {
	//TODO:
}
