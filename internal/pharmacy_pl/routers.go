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
    "github.com/gin-gonic/gin"
)

func AddRoutes(engine *gin.Engine) {
  group := engine.Group("/api")
  
  {
    api := newAmbulancesAPI()
    api.addRoutes(group)
  }
  
  {
    api := newMedicineOrdersAPI()
    api.addRoutes(group)
  }
  
  {
    api := newPrescriptionsAPI()
    api.addRoutes(group)
  }
  
}