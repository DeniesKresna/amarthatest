package one

import (
	"github.com/DeniesKresna/amarthatest/config"
	"github.com/gin-gonic/gin"
)

func InitQuestion(r *gin.Engine, conf *config.Config) {
	repos := repo{
		db: conf.DB,
	}
	handler := controller{
		repo: repos,
	}

	// create loan for users
	r.POST("/loan/create", handler.CreateLoan)

	// get loan by user id
	r.GET("/loan/user/:id", handler.GetLoansByUserID)

	// get all users with the profile
	r.GET("/users", handler.GetAllUserWithProfileList)

	// get repayment history of loan
	r.GET("/repayment/loan/code/:code", handler.GetRepaymentsHistoryByLoanCode)

	// generate payment invoice
	r.POST("/repayment/create", handler.GeneratePaymentInvoice)

	// pay user repayment invoices
	r.POST("/repayment/pay", handler.PayLoanInstallment)

	go InitCron(handler)
}
