package pharmacy_pl

import (
	"net/http"
	"time"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func validatePrescription(entry Prescription) []string {
	var errorMessages []string

	if entry.Id == "" {
		errorMessages = append(errorMessages, "ID is required")
	}

	if entry.PatientName == "" {
		errorMessages = append(errorMessages, "Patient Name is required")
	}

	if entry.PatientId == "" {
		errorMessages = append(errorMessages, "Patient ID is required")
	}

	if entry.DoctorName == "" {
		errorMessages = append(errorMessages, "Doctor Name is required")
	}

	if entry.ValidUntil.IsZero() { // Assuming ValidUntil is a time.Time type
		errorMessages = append(errorMessages, "Valid Until is required and must be a valid date")
	}

	if entry.Medicines == nil || len(entry.Medicines) == 0 { // Assuming Medicines is a slice
		errorMessages = append(errorMessages, "At least one medicine is required")
	}

	// If there are error messages, return them with a Bad Request status
	if len(errorMessages) > 0 {
		return errorMessages
	}

	// If no errors, return nil for the gin.H and 0 for the status code
	return nil
}

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

		if entry.Id == "" || entry.Id == "@new" {
			entry.Id = uuid.NewString()
		}

		if entry.IssuedDate.IsZero() {
			entry.IssuedDate = time.Now()
		}

		defaultValidTime := 7 * 24 * time.Hour // One week

		if entry.ValidUntil.IsZero() {
			entry.ValidUntil = time.Now().Add(defaultValidTime)
		}

		// Validate the Prescription
		errorMessages := validatePrescription(entry)
		if errorMessages != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Required information is missing",
				"errors":  errorMessages,
			}, http.StatusBadRequest
		}

		conflictIndx := slices.IndexFunc(ambulance.PrescriptionList, func(prescription Prescription) bool {
			return entry.Id == prescription.Id
		})

		if conflictIndx >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "Prescription with this ID already exists",
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
