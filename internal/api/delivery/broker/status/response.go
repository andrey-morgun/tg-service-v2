package status

import "time"

type response struct {
	Name      string    `json:"name"`
	BuildDate time.Time `json:"build_date"`
}
