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

import (
	//  "net/http"

	"github.com/gin-gonic/gin"
)

type PrescriptionsAPI interface {

	// internal registration of api routes
	addRoutes(routerGroup *gin.RouterGroup)

	// CreatePrescription - Creates a new prescription
	CreatePrescription(ctx *gin.Context)

	// DeletePrescription - Deletes a specific prescription
	DeletePrescription(ctx *gin.Context)

	// GetAmbulancePrescriptions - Provides the list of prescriptions for an ambulance
	GetAmbulancePrescriptions(ctx *gin.Context)

	// GetPrescriptionById - Retrieves a specific prescription by ID
	GetPrescriptionById(ctx *gin.Context)

	// UpdatePrescription - Updates a specific prescription
	UpdatePrescription(ctx *gin.Context)
}

// partial implementation of PrescriptionsAPI - all functions must be implemented in add on files
type implPrescriptionsAPI struct {
}

func newPrescriptionsAPI() PrescriptionsAPI {
	return &implPrescriptionsAPI{}
}

func (api *implPrescriptionsAPI) addRoutes(routerGroup *gin.RouterGroup) {
	prescriptionsGroup := routerGroup.Group("/ambulances/:ambulanceId/prescriptions")
	prescriptionsGroup.POST("", api.CreatePrescription)
	prescriptionsGroup.GET("", api.GetAmbulancePrescriptions)
	prescriptionsGroup.DELETE("/:prescriptionId", api.DeletePrescription)
	prescriptionsGroup.PUT("/:prescriptionId", api.UpdatePrescription)
	prescriptionsGroup.GET("/:prescriptionId", api.GetPrescriptionById)

}

// Copy following section to separate file, uncomment, and implement accordingly
