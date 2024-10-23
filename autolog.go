package autolog

import (
	"log"
	"time"
)

var (
	Location *time.Location
	TZTokyo  = "Asia/Tokyo"
)

func init() {
	var err error

	Location, err = time.LoadLocation(TZTokyo)
	if err != nil {
		log.Fatal(err)
	}
}
