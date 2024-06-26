/*
 * Pharmacy Prescription API
 *
 * Pharmacy Prescription management system
 *
 * API version: 1.0.0
 * Contact: xkabac@stuba.sk
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package pharmacy_pl

type Ambulance struct {

	// Unique identifier for the ambulance
	Id string `json:"id"`

	// Human readable display name of the ambulance
	Name string `json:"name"`

	PrescriptionList []Prescription `json:"prescriptionList,omitempty"`

	MedicineOrderList []MedicineOrder `json:"medicineOrderList,omitempty"`
}
