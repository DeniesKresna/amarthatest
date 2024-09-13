package one

import (
	"database/sql"
	"time"

	"github.com/DeniesKresna/amarthatest/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func (r *repo) GetTx(c *gin.Context) *gorm.DB {
	return r.db.Begin()
}

func (r *repo) GetLoan(c *gin.Context, loan *models.Loan) (err error) {
	return r.db.Where(*loan).Order("created_at desc").First(loan).Error
}

func (r *repo) GetRepayment(c *gin.Context, repayment *models.Repayment) (err error) {
	return r.db.Where(*repayment).Order("created_at desc").First(repayment).Error
}

func (r *repo) GetLoanList(c *gin.Context, loan models.Loan) (loans []models.Loan, err error) {
	err = r.db.Where(loan).Order("created_at desc").Find(&loans).Error
	return
}

func (r *repo) GetLastDueRepayment(c *gin.Context, repayment *models.Repayment) (err error) {
	return r.db.Where(*repayment).Order("pay_due_date desc").First(repayment).Error
}

func (r *repo) GetRepaymentList(c *gin.Context, repayment models.Repayment) (repayments []models.Repayment, err error) {
	err = r.db.Where(repayment).Order("created_at desc").Find(&repayments).Error
	return
}

func (r *repo) GetPaidRepaymentTotal(c *gin.Context, repayment models.Repayment) (totalPayAmount float64, err error) {
	var totalPayment sql.NullFloat64
	err = r.db.Model(&models.Repayment{}).
		Select("SUM(pay_amount)").
		Where("pay_code = ?", repayment.PayCode).
		Where("paid_at is not null").
		Scan(&totalPayment).Error
	if err != nil {
		return
	}

	if totalPayment.Valid {
		totalPayAmount = totalPayment.Float64
	}
	return
}

func (r *repo) GetUnpaidRepayments(c *gin.Context, repayment models.Repayment) (repayments []models.Repayment, err error) {
	err = r.db.Where("loan_code = ?", repayment.LoanCode).
		Where("paid_at is null").
		Order("pay_due_date").Find(&repayments).Error
	return
}

func (r *repo) GetActiveLoanByUserID(c *gin.Context, userID int64) (loan models.Loan, err error) {
	err = r.db.Where("user_id = ?", userID).
		Where("(is_paid_off = ? or is_paid_off is null)", false).
		Where("outstanding > ?", 0).
		Order("created_at desc").
		First(&loan).Error

	return
}

func (r *repo) GetActiveLoanByLoanCode(c *gin.Context, code string) (loan models.Loan, err error) {
	err = r.db.Where("code = ?", code).
		Where("(is_paid_off = ? or is_paid_off is null)", false).
		Where("outstanding > ?", 0).
		Order("created_at desc").
		First(&loan).Error

	return
}

func (r *repo) GetAllActiveLoans(c *gin.Context) (loans []models.Loan, err error) {
	err = r.db.Where("(is_paid_off = ? or is_paid_off is null)", false).
		Where("outstanding > ?", 0).
		Order("created_at desc").
		Find(&loans).Error

	return
}

func (r *repo) GetAllRepaymentsByLoanCodeOnSpecificDay(c *gin.Context, loanCode string, day time.Time) (repayments []models.Repayment, err error) {
	err = r.db.Where("loan_code = ?", loanCode).Where("DATE(created_at) = ?", day.Format("2006-01-02")).Find(&repayments).Error

	return
}

func (r *repo) CreateLoan(c *gin.Context, loan *models.Loan) (err error) {
	db := getTx(c, r.db)
	return db.Create(loan).Error
}

func (r *repo) GetUser(c *gin.Context, user *models.User) (err error) {
	return r.db.Where(*user).Order("created_at desc").First(user).Error
}

func (r *repo) GetProfile(c *gin.Context, profile *models.Profile) (err error) {
	return r.db.Where(*profile).Order("created_at desc").First(profile).Error
}

func (r *repo) GetAllUserWithProfileList(c *gin.Context) (users []models.UserWithProfile, err error) {
	err = r.db.Raw(`select u.id, name, email, platform_account, npl_status from users u
		left join profiles p on p.user_id = u.id where u.deleted_at is null order by u.created_at desc
	`).Find(&users).Error
	return
}

func (r *repo) CreateRepayment(c *gin.Context, repayment *models.Repayment) (err error) {
	db := getTx(c, r.db)
	return db.Create(repayment).Error
}

func (r *repo) UpdateProfile(c *gin.Context, profile *models.Profile) (err error) {
	db := getTx(c, r.db)
	return db.Updates(profile).Error
}

func (r *repo) UpdateLoan(c *gin.Context, loan *models.Loan) (err error) {
	db := getTx(c, r.db)
	return db.Updates(loan).Error
}

func (r *repo) UpdateRepayment(c *gin.Context, repayment *models.Repayment) (err error) {
	db := getTx(c, r.db)
	return db.Updates(repayment).Error
}
