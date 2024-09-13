package models

import (
	"time"

	"gorm.io/gorm"
)

type Loan struct {
	ID          int64          `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
	UserID      int64          `json:"user_id" binding:"required"`
	Code        string         `json:"code"`
	Amount      int64          `json:"amount" binding:"required"`
	DisbursedAt *time.Time     `json:"disbursed_at"`
	IsPaidOff   *bool          `json:"is_paid_off"`
	Outstanding float64        `json:"outstanding"`
}

type Repayment struct {
	ID         int64          `json:"id" gorm:"primarykey"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
	LoanCode   string         `json:"loan_code" binding:"required"`
	PayAmount  float64        `json:"pay_amount"`
	PayCode    string         `json:"pay_code"`
	PayDueDate *time.Time     `json:"pay_due_date"`
	PaidAt     *time.Time     `json:"paid_at"`
}

type RepaymentInstallmentPayload struct {
	LoanCode  string  `json:"loan_code" binding:"required"`
	PayAmount float64 `json:"pay_amount" binding:"required"`
	PayCode   string  `json:"pay_code" binding:"required"`
}
