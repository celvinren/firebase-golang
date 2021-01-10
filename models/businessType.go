package models

// BusinessType for dto
type BusinessType struct {
	ID              string            `json:"id"`
	BusinessType    string            `json:"businessType"`
	BusinessSubType []BusinessSubType `json:"businessSubType"`
}

// BusinessSubType for dto
type BusinessSubType struct {
	ID              string `json:"id"`
	BusinessSubType string `json:"businessSubType"`
}
