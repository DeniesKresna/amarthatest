package one

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/DeniesKresna/amarthatest/helpers"
	"github.com/DeniesKresna/amarthatest/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type controller struct {
	repo repo
}

func (k *controller) CreateLoan(c *gin.Context) {
	var loan models.Loan
	var user models.User

	err := c.ShouldBindJSON(&loan)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// loan data injection
	{
		loan.Code = helpers.CreateRandomString(15)

		loan.Outstanding = totalPaymentFromLoan(float64(loan.Amount))
		fmt.Println(loan)

		// i assume its disbursed immediatelly
		disburseAt := time.Now()
		loan.DisbursedAt = &disburseAt

		paidOff := false
		loan.IsPaidOff = &paidOff
	}

	// check if user exist
	user.ID = loan.UserID
	err = k.repo.GetUser(c, &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// i check if user still has active loan
	existedUserLoan, err := k.repo.GetActiveLoanByUserID(c, loan.UserID)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User still have active loan", "data": existedUserLoan})
		return
	} else {
		if !errors.Is(gorm.ErrRecordNotFound, err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// check if loan code has been exist since it was unique.
	// by the wait can be handled by db layer via unique column constraint
	err = k.repo.GetLoan(c, &loan)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Loan Code exist"})
		return
	} else {
		if !errors.Is(gorm.ErrRecordNotFound, err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	err = k.repo.CreateLoan(c, &loan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"data": loan})
}

func (k *controller) GetLoansByUserID(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loans, err := k.repo.GetLoanList(c, models.Loan{UserID: userID})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"data": loans})
}

func (k *controller) GetAllUserWithProfileList(c *gin.Context) {
	users, err := k.repo.GetAllUserWithProfileList(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"data": users})
}

func (k *controller) GetRepaymentsHistoryByLoanCode(c *gin.Context) {
	code := c.Param("code")
	users, err := k.repo.GetRepaymentList(c, models.Repayment{
		LoanCode: code,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"data": users})
}

func (k *controller) GeneratePaymentInvoice(c *gin.Context) {
	var (
		repayment models.Repayment
		loan      models.Loan
		profile   models.Profile
		err       error
	)

	// use tx
	c.Set(DBTX, k.repo.GetTx(c))
	defer func() {
		endTx(c, err)
	}()

	err = c.ShouldBindJSON(&repayment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check active loan by code
	loan, err = k.repo.GetActiveLoanByLoanCode(c, repayment.LoanCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check user profile exist
	profile.UserID = loan.UserID
	err = k.repo.GetProfile(c, &profile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check wheter user is delinquent
	{
		// get creation date of last paid repayment
		var isDelinquent bool
		var unpaidRepayments []models.Repayment
		unpaidRepayments, err = k.repo.GetUnpaidRepayments(c, models.Repayment{
			LoanCode: loan.Code,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(unpaidRepayments) > 0 {
			// if there is unpaid invoices
			for _, unpaidRepayment := range unpaidRepayments {
				// check now if two weeks after last unpaid active loan
				twoWeeksLater := unpaidRepayment.CreatedAt.AddDate(0, 0, 14)
				if time.Now().After(twoWeeksLater) {
					// inject for deliquet status
					profile.NplStatus = models.NPL_DELINQUENT_PAYMENT
					isDelinquent = true
					break
				}
			}
		} else {
			twoWeeksLater := loan.DisbursedAt.AddDate(0, 0, 14)
			if time.Now().After(twoWeeksLater) {
				// inject for deliquet status
				profile.NplStatus = models.NPL_DELINQUENT_PAYMENT
				isDelinquent = true
			}
		}
		if !isDelinquent {
			profile.NplStatus = models.NPL_SMOOTH_PAYMENT
			isDelinquent = true
		}
	}

	// update loan and profile
	{
		err = k.repo.UpdateLoan(c, &loan)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = k.repo.UpdateProfile(c, &profile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// set outstanding and paidoff status of loan
	var outstanding float64
	{
		totalPaidPayments, err := k.repo.GetPaidRepaymentTotal(c, repayment)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// inject for updating outstanding later
		outstanding = float64(loan.Amount) - totalPaidPayments
		loan.Outstanding = outstanding

		// if outstanding is 0 or lower, means that this loan has been paid off
		isLoanPaidOff := false
		if outstanding <= 0 {
			isLoanPaidOff = true
		}
		loan.IsPaidOff = &isLoanPaidOff

		if isLoanPaidOff {

			c.JSON(http.StatusBadRequest, gin.H{"error": "this loan has been paid off"})
			return
		}
	}

	// repayment injection
	{
		// re-init
		repayment = models.Repayment{
			LoanCode: repayment.LoanCode,
		}

		// inject payment code
		trial := 0
		for {
			trial++
			payCode := fmt.Sprintf("%s-%s", loan.Code, helpers.CreateRandomString(7))
			existRepayment := models.Repayment{
				PayCode: payCode,
			}

			err = k.repo.GetRepayment(c, &existRepayment)
			if err != nil {
				if !errors.Is(gorm.ErrRecordNotFound, err) {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				repayment.PayCode = payCode
				break
			}
			if trial == 3 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "paycode generate existed code 3 times"})
				return
			}
		}

		installment := getInstalmentFromLoanAmount(float64(loan.Amount))
		// if outstanding same or lower than installment, so this is last installment of client
		if outstanding < installment {
			installment = outstanding
		}
		repayment.PayAmount = installment

		// i will check for latest repayment on this loan code
		var nextWeekLater time.Time
		lastRepayment := models.Repayment{
			LoanCode: loan.Code,
		}
		err = k.repo.GetLastDueRepayment(c, &lastRepayment)
		if err == nil {
			// if previous repayment exist, i will add 7 days from due date as new repayment due date
			dueDate := *lastRepayment.PayDueDate
			nextWeekLater = dueDate.AddDate(0, 0, 7)
		} else {
			if !errors.Is(gorm.ErrRecordNotFound, err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			// if previous repayment exist, i will add 7 days from disbursed time as new repayment due date
			dueDate := *loan.DisbursedAt
			nextWeekLater = dueDate.AddDate(0, 0, 7)
		}
		repayment.PayDueDate = &nextWeekLater

		// create repayment invoice
		err = k.repo.CreateRepayment(c, &repayment)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"data": repayment})
}

func (k *controller) PayLoanInstallment(c *gin.Context) {
	var (
		payload models.RepaymentInstallmentPayload
		loan    models.Loan
		profile models.Profile
		err     error
	)

	// use tx
	c.Set(DBTX, k.repo.GetTx(c))
	defer func() {
		endTx(c, err)
	}()

	err = c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check active loan by code
	loan, err = k.repo.GetActiveLoanByLoanCode(c, payload.LoanCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check user profile exist
	profile.UserID = loan.UserID
	err = k.repo.GetProfile(c, &profile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get repayment
	repayment := models.Repayment{
		LoanCode:  payload.LoanCode,
		PayAmount: payload.PayAmount,
		PayCode:   payload.PayCode,
	}
	err = k.repo.GetRepayment(c, &repayment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// pay the repayment data
	now := time.Now()
	repayment.PaidAt = &now

	err = k.repo.UpdateRepayment(c, &repayment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check wheter user is delinquent
	{
		// get creation date of last paid repayment
		var isDelinquent bool
		var unpaidRepayments []models.Repayment
		unpaidRepayments, err = k.repo.GetUnpaidRepayments(c, models.Repayment{
			LoanCode: loan.Code,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(unpaidRepayments) > 0 {
			// if there is unpaid invoices
			for _, unpaidRepayment := range unpaidRepayments {
				// check now if two weeks after last unpaid active loan
				twoWeeksLater := unpaidRepayment.CreatedAt.AddDate(0, 0, 14)
				if time.Now().After(twoWeeksLater) {
					// inject for deliquet status
					profile.NplStatus = models.NPL_DELINQUENT_PAYMENT
					isDelinquent = true
					break
				}
			}
		} else {
			twoWeeksLater := loan.DisbursedAt.AddDate(0, 0, 14)
			if time.Now().After(twoWeeksLater) {
				// inject for deliquet status
				profile.NplStatus = models.NPL_DELINQUENT_PAYMENT
				isDelinquent = true
			}
		}
		if !isDelinquent {
			profile.NplStatus = models.NPL_SMOOTH_PAYMENT
			isDelinquent = true
		}
	}

	// update loan and profile
	{
		err = k.repo.UpdateLoan(c, &loan)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = k.repo.UpdateProfile(c, &profile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"data": repayment})
}

func (k *controller) GeneratesLoanInvoicesOnCron() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// first i get all unpaidoff loans
	var (
		loans []models.Loan
		err   error
	)
	loans, err = k.repo.GetAllActiveLoans(c)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	today := time.Now()

	for _, loan := range loans {
		// check if the loan disbursed at and today is in 7 days multiple
		daysBetween := int(math.Abs(today.Sub(*loan.DisbursedAt).Hours() / 24))
		isMultipleBy7 := false
		if daysBetween%7 == 0 {
			isMultipleBy7 = true
		}

		// if yes, check if loan has some invoices today.
		var repayments []models.Repayment
		if isMultipleBy7 {
			repayments, err = k.repo.GetAllRepaymentsByLoanCodeOnSpecificDay(c, loan.Code, today)
			if err != nil {
				continue
			}

			// only create new invoice when no invoices generated today
			if len(repayments) <= 0 {
				jsonPayload := fmt.Sprintf(`{"loan_code": "%s"}`, loan.Code)
				req, errs := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(jsonPayload)))
				if errs != nil {
					continue
				}

				// prepare hiting generate payment invoice
				req.Header.Set("Content-Type", "application/json")
				innerW := httptest.NewRecorder()
				innerC, _ := gin.CreateTestContext(innerW)
				innerC.Request = req
				go k.GeneratePaymentInvoice(innerC)
			}
		}
	}
	fmt.Println("daily check invoice generation done")
}
