package models

import (
	"encoding/json"
	"log"
)

// StoreDTO for DTO
type StoreDTO struct {
	StoreName          string `json:"storeName"`
	Address            string `json:"address"`
	BusinessType       string `json:"businessType"`
	BusinessSubType    string `json:"businessSubType"`
	BackgroundColorHex string `json:"backgroundColorHex"`
	OrganizationID     string `json:"organizationID"`
	Logo               string `json:"logo"`
	Phone              string `json:"phone"`
	Email              string `json:"email"`
}

// Store godoc
type Store struct {
	StoreName          string `json:"storeName,omitempty"`
	Address            string `json:"address,omitempty"`
	BusinessType       string `json:"businessType,omitempty"`
	BusinessSubType    string `json:"businessSubType,omitempty"`
	BackgroundColorHex string `json:"backgroundColorHex,omitempty"`
	OrganizationID     string `json:"organizationID,omitempty"`
	Logo               string `json:"logo,omitempty"`
	Phone              string `json:"phone,omitempty"`
	Email              string `json:"email,omitempty"`
}

// // StoretoJSON to json
// func StoretoJSON(store Store) interface{} {
// 	return gin.H{
// 		"storeName":      store.StoreName,
// 		"address":        store.Address,
// 		"phone":          store.Phone,
// 		"email":          store.Email,
// 		"organizationID": store.OrganizationID,
// 	}
// }

// OrganizationToJSON to json
func StoreFromJSON(storeStr string) Store {
	var store Store
	err := json.Unmarshal([]byte(storeStr), &store)
	if err != nil {
		log.Println(err.Error())
	}
	return store
}
