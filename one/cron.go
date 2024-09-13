package one

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func InitCron(k controller) {
	c := cron.New()
	// Cron every day at 8
	_, err := c.AddFunc("0 8 * * *", func() {
		k.GeneratesLoanInvoicesOnCron()
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	c.Start()
	fmt.Println("Cron started")

	select {}
}
