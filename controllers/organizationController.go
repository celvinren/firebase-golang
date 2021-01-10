package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	models "restaurant_golang/models"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

// CreateOrganization create organization
func CreateOrganization(organization models.Organization) (_ bool, err error) {
	isCreated := false
	var errGlobal error
	client, err := InitFireStore()

	if err != nil {
		errGlobal = errors.New("fail to create firebase client: " + err.Error())
	}

	_, _, errCreate := client.Collection("organization").Add(context.Background(), organization)
	if errCreate != nil {
		errGlobal = errors.New("An error has occurred when update organization:" + errCreate.Error())
	}
	isCreated = true
	return isCreated, errGlobal
}

// UpdateOrganization godoc
// @Summary Update Organization
// @Description Update Organization
// @Success 200 {string} json
// @Router /apis/organization/{id} [Put]
func UpdateOrganization(c *gin.Context) {
	id := c.Param("organizationId")
	var organizationDTO models.OrganizationDTO
	bindErr := c.BindJSON(&organizationDTO)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "json cannot be explain: " + bindErr.Error(),
			"data":   nil,
		})
		return
	}
	log.Print(id)
	log.Print(organizationDTO)

	organization := models.Organization{
		OrganizationName: organizationDTO.OrganizationName,
		Address:          organizationDTO.Address,
		Phone:            organizationDTO.Phone,
		Email:            organizationDTO.Email,
		UserID:           organizationDTO.UserID,
	}
	client, err := InitFireStore()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "fail to create firebase client: " + err.Error(),
			"data":   nil,
		})
		return
	}
	_, errCreate := client.Collection("organization").Doc(id).Set(context.Background(), organization)
	if errCreate != nil {
		log.Printf("An error has occurred when update organization: %s", errCreate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "An error has occurred when update organization:" + errCreate.Error(),
			"data":   nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "Successfully update organization",
		"data":   nil,
	})
}

func GetOrganzation(c *gin.Context) {
	var org models.OrganizationDTO

	userID := c.Param("userId")
	log.Print(userID)

	client, err := InitFireStore()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "fail to fetch organization details: " + err.Error(),
			"data":   nil,
		})
		return
	}

	orgIter := client.Collection("organization").Where("userID", "==", userID).Documents(context.Background())
	for {
		doc, err := orgIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf(err.Error())
		}
		log.Println(doc.Data())
		jsonString, err := json.Marshal(doc.Data())
		org = models.OrganizationFromJSON(string(jsonString))
		org.OrganizationID = doc.Ref.ID
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "Successfully fetch organization details",
		"data":   org,
	})
}
