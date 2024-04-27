package pharmacy_pl

import (
	"net/http"
	"time"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Nasledujúci kód je kópiou vygenerovaného a zakomentovaného kódu zo súboru api_ambulance_waiting_list.go

// CreatePrescription - Saves new entry into waiting list
func (api *implPrescriptionsAPI) CreatePrescription(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		var entry Prescription

		if err := c.ShouldBindJSON(&entry); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		if entry.PatientName == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Patient Name is required",
			}, http.StatusBadRequest
		}

		if entry.DoctorName == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Doctor Name is required",
			}, http.StatusBadRequest
		}

		if entry.Id == "" || entry.Id == "@new" {
			entry.Id = uuid.NewString()
		}

		conflictIndx := slices.IndexFunc(ambulance.PrescriptionList, func(prescription Prescription) bool {
			return entry.Id == prescription.Id || entry.PatientName == prescription.PatientName
		})

		if conflictIndx >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "Entry already exists",
			}, http.StatusConflict
		}

		ambulance.PrescriptionList = append(ambulance.PrescriptionList, entry)

		// entry was copied by value return reconciled value from the list
		entryIndx := slices.IndexFunc(ambulance.PrescriptionList, func(prescription Prescription) bool {
			return entry.Id == prescription.Id
		})
		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to save entry",
			}, http.StatusInternalServerError
		}
		return ambulance, ambulance.PrescriptionList[entryIndx], http.StatusOK
	})

}

// DeletePrescription - Deletes specific entry
func (api *implPrescriptionsAPI) DeletePrescription(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		entryId := ctx.Param("prescriptionId")

		if entryId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Prescription ID is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(ambulance.PrescriptionList, func(prescription Prescription) bool {
			return entryId == prescription.Id
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		ambulance.PrescriptionList = append(ambulance.PrescriptionList[:entryIndx], ambulance.PrescriptionList[entryIndx+1:]...)
		return ambulance, nil, http.StatusNoContent
	})
}

// GetAmbulancePrescriptions - Provides the ambulance waiting list
func (api *implPrescriptionsAPI) GetAmbulancePrescriptions(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		result := ambulance.PrescriptionList
		if result == nil {
			result = []Prescription{}
		}
		// return nil ambulance - no need to update it in db
		return nil, result, http.StatusOK
	})
}

// GetPrescriptionById - Provides details about waiting list entry
func (api *implPrescriptionsAPI) GetPrescriptionById(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		entryId := ctx.Param("entryId")

		if entryId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Entry ID is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(ambulance.PrescriptionList, func(prescription Prescription) bool {
			return entryId == prescription.Id
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		// return nil ambulance - no need to update it in db
		return nil, ambulance.PrescriptionList[entryIndx], http.StatusOK
	})
}

// UpdatePrescription - Updates specific entry
func (api *implPrescriptionsAPI) UpdatePrescription(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		var entry Prescription

		if err := c.ShouldBindJSON(&entry); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		prescriptionId := ctx.Param("prescriptionId")

		if prescriptionId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Entry ID is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(ambulance.PrescriptionList, func(prescription Prescription) bool {
			return prescriptionId == prescription.Id
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		if entry.PatientName != "" {
			ambulance.PrescriptionList[entryIndx].PatientName = entry.PatientName
		}

		if entry.DoctorName != "" {
			ambulance.PrescriptionList[entryIndx].DoctorName = entry.DoctorName
		}

		if entry.Instructions != "" {
			ambulance.PrescriptionList[entryIndx].Instructions = entry.Instructions
		}

		if entry.Notes != "" {
			ambulance.PrescriptionList[entryIndx].Notes = entry.Notes
		}

		if entry.Status != "" {
			ambulance.PrescriptionList[entryIndx].Status = entry.Status
		}

		if entry.Id != "" {
			ambulance.PrescriptionList[entryIndx].Id = entry.Id
		}

		if entry.IssuedDate.After(time.Time{}) {
			ambulance.PrescriptionList[entryIndx].IssuedDate = entry.IssuedDate
		}

		if entry.ValidUntil.After(time.Time{}) {
			ambulance.PrescriptionList[entryIndx].ValidUntil = entry.ValidUntil
		}

		if entry.Medicines != nil {
			ambulance.PrescriptionList[entryIndx].Medicines = entry.Medicines
		}

		return ambulance, ambulance.PrescriptionList[entryIndx], http.StatusOK
	})
}
