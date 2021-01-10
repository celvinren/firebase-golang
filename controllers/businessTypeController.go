package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	models "restaurant_golang/models"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

// GetBusinessType doc
func GetBusinessType(c *gin.Context) {
	//TODO:
	client, err := InitFireStore()
	if err != nil {
		log.Printf(err.Error())
		return
	}
	var businessTypeList []models.BusinessType
	iter := client.Collection("businessType").Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf(err.Error())
		}

		businessTypeList = append(businessTypeList, models.BusinessType{ID: doc.Ref.ID, BusinessType: doc.Data()["businessType"].(string)})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "Successfully get business type",
		"data":   businessTypeList,
	})
}

// GetBusinessSubType doc
func GetBusinessSubType(c *gin.Context) {
	businessParentTypeID := c.Param("id")
	log.Println(businessParentTypeID)

	client, err := InitFireStore()
	if err != nil {
		log.Printf(err.Error())
		return
	}

	iter := client.Collection("businessSubType").Where("businessTypeID", "==", businessParentTypeID).Documents(context.Background())
	fmt.Printf("%v", iter)
	var businessSubTypeList []models.BusinessSubType
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf(err.Error())
		}
		fmt.Println(doc.Data())
		businessSubTypeList = append(businessSubTypeList, models.BusinessSubType{ID: doc.Ref.ID, BusinessSubType: doc.Data()["businessSubType"].(string)})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "Successfully get business type",
		"data":   businessSubTypeList,
	})
}
