package one

import (
	"math"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func totalPaymentFromLoan(amount float64) float64 {
	return 1.1 * amount
}

func getInstalmentFromLoanAmount(amount float64) float64 {
	installment := totalPaymentFromLoan(amount) / 50

	return math.Ceil(installment)
}

const DBTX = "dbtx"

func getTx(c *gin.Context, db *gorm.DB) *gorm.DB {
	if valInt, ok := c.Get(DBTX); !ok {
		return db
	} else {
		if val, ok := valInt.(*gorm.DB); ok {
			return val
		}
	}
	return db
}

func endTx(c *gin.Context, err error) {
	if valInt, ok := c.Get(DBTX); ok {
		if val, ok := valInt.(*gorm.DB); ok {
			if err != nil {
				val.Rollback()
			} else {
				val.Commit()
			}
		}
	}
}
