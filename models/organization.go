package models

import (
	"encoding/json"
	"log"
)

// OrganizationDTO for dto
type OrganizationDTO struct {
	OrganizationName string `json:"organizationName"`
	Address          string `json:"address"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	UserID           string `json:"userID"`
	OrganizationID   string `json:"organizationID"`
}

// Organization organization object
type Organization struct {
	OrganizationName string `json:"organizationName,omitempty"`
	Address          string `json:"address,omitempty"`
	Phone            string `json:"phone,omitempty"`
	Email            string `json:"email,omitempty"`
	UserID           string `json:"userID,omitempty"`
}

// // OrganizationToJSON to json
// func OrganizationToJSON(organization Organization) interface{} {
// 	return gin.H{
// 		"organizationName": organization.OrganizationName,
// 		"address":          organization.Address,
// 		"phone":            organization.Phone,
// 		"email":            organization.Email,
// 		"userID":           organization.UserID,
// 	}
// }

// OrganizationToJSON to json
func OrganizationFromJSON(organizationStr string) OrganizationDTO {
	var org OrganizationDTO
	err := json.Unmarshal([]byte(organizationStr), &org)
	if err != nil {
		log.Println(err.Error())
	}
	return org
}
