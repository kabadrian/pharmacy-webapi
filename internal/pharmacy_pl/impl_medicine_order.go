package pharmacy_pl

import (
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// validateMedicineOrder - Validates the medicine order
func validateMedicineOrder(order MedicineOrder) []string {
	var errorMessages []string

	if order.OrderId == "" {
		errorMessages = append(errorMessages, "Order ID is required")
	}

	if order.OrderedBy == "" {
		errorMessages = append(errorMessages, "Ordered By is required")
	}

	if order.Medicines == nil || len(order.Medicines) == 0 {
		errorMessages = append(errorMessages, "At least one medicine is required")
	}

	return errorMessages
}

// CreateMedicineOrder - Creates a new medicine order
func (api *implMedicineOrdersAPI) CreateMedicineOrder(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		var entry MedicineOrder

		if err := c.ShouldBindJSON(&entry); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		if entry.OrderId == "" || entry.OrderId == "@new" {
			entry.OrderId = uuid.NewString()
		}

		if entry.OrderDate.IsZero() {
			entry.OrderDate = time.Now()
		}

		// Validate the MedicineOrder
		errorMessages := validateMedicineOrder(entry)
		if errorMessages != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Required information is missing",
				"errors":  errorMessages,
			}, http.StatusBadRequest
		}

		conflictIndx := slices.IndexFunc(ambulance.MedicineOrderList, func(medicineOrder MedicineOrder) bool {
			return entry.OrderId == medicineOrder.OrderId
		})

		if conflictIndx >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "MedicineOrder with api ID already exists",
			}, http.StatusConflict
		}

		ambulance.MedicineOrderList = append(ambulance.MedicineOrderList, entry)

		// entry was copied by value return reconciled value from the list
		entryIndx := slices.IndexFunc(ambulance.MedicineOrderList, func(medicineOrder MedicineOrder) bool {
			return entry.OrderId == medicineOrder.OrderId
		})
		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to save entry",
			}, http.StatusInternalServerError
		}
		return ambulance, ambulance.MedicineOrderList[entryIndx], http.StatusOK
	})
}

// DeleteMedicineOrder - Deletes a specific medicine order
func (api *implMedicineOrdersAPI) DeleteMedicineOrder(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		entryId := ctx.Param("orderId")

		if entryId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "MedicineOrder ID is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(ambulance.MedicineOrderList, func(medicineOrder MedicineOrder) bool {
			return entryId == medicineOrder.OrderId
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		ambulance.MedicineOrderList = append(ambulance.MedicineOrderList[:entryIndx], ambulance.MedicineOrderList[entryIndx+1:]...)
		return ambulance, nil, http.StatusNoContent
	})
}

// GetAllMedicineOrders - Retrieves all medicine orders for an ambulance
func (api *implMedicineOrdersAPI) GetAllMedicineOrders(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		result := ambulance.MedicineOrderList
		if result == nil {
			result = []MedicineOrder{}
		}
		// return nil ambulance - no need to update it in db
		return nil, result, http.StatusOK
	})
}

// GetMedicineOrderById - Retrieves a specific medicine order by ID
func (api *implMedicineOrdersAPI) GetMedicineOrderById(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		entryId := ctx.Param("orderId")

		if entryId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Entry ID is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(ambulance.MedicineOrderList, func(medicineOrder MedicineOrder) bool {
			return entryId == medicineOrder.OrderId
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		// return nil ambulance - no need to update it in db
		return nil, ambulance.MedicineOrderList[entryIndx], http.StatusOK
	})
}

// UpdateMedicineOrder - Updates a specific medicine order
func (api *implMedicineOrdersAPI) UpdateMedicineOrder(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		var entry MedicineOrder

		if err := c.ShouldBindJSON(&entry); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		orderId := ctx.Param("orderId")

		if orderId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Entry ID is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(ambulance.MedicineOrderList, func(medicineOrder MedicineOrder) bool {
			return orderId == medicineOrder.OrderId
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		if entry.Medicines != nil {
			ambulance.MedicineOrderList[entryIndx].Medicines = entry.Medicines
		}

		return ambulance, ambulance.MedicineOrderList[entryIndx], http.StatusOK
	})
}
