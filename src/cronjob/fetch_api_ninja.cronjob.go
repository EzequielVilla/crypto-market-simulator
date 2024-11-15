package cronjob

import (
	"crypto-market-simulator/internal/lib"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron"
	"net/http"
)

type MyCronjob struct {
	cronjob *cron.Cron
}

func (m *MyCronjob) FetchUpdateCryptoValues() {
	_ = "0 * * * * *" // everyMinute
	everyThreeHours := "0 0 */3 * * *"

	_ = m.cronjob.AddFunc(everyThreeHours, func() {
		fmt.Println("START_CRONJOB_UPDATE_CRYPTO_VALUES")
		url := "http://localhost:3000/api/crypto/values"
		req, err := http.NewRequest("PATCH", url, nil)
		if err != nil {
			fmt.Printf("ERROR_CRONJOB_UPDATE_CRYPTO_VALUES: %v\n", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("ERROR_CRONJOB_UPDATE_CRYPTO_VALUES: %v\n", err)
			return
		}
		var data lib.Response
		if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
			fmt.Printf("ERROR_CRONJOB_UPDATE_CRYPTO_VALUES: %v\n", err)
			return
		}
		fmt.Printf("END_CRONJOB_UPDATE_CRYPTO_VALUES: %v\n", data.Message)
	})
	m.cronjob.Start()
}

type IMyCronjob interface {
	FetchUpdateCryptoValues()
}

func NewMyCronjob() IMyCronjob {
	c := cron.New()
	return &MyCronjob{
		cronjob: c,
	}
}
