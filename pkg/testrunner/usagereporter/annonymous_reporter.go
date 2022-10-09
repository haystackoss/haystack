package usagereporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/nabaz-io/nabaz/pkg/testrunner/models"
)

func SendTelemetry(telemetry models.Telemetry) error {
	j, err := json.Marshal(telemetry)
    if err != nil {
        fmt.Println("failed to marshal telemetry")
        return err
    }

	res, err := http.Post("https://api.nabaz.io/stats", "application/json", bytes.NewBuffer(j))
	if err != nil || res.StatusCode != 200 {
		fmt.Println("failed to send telemetry")
		return err
	}

	return nil
}