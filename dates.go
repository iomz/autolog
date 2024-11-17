package autolog

import (
	"fmt"
	"os"
	"time"
)

func LogDates(filename string) {
	today := time.Now()

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%s\n", today.Format("2006-01-02"))); err != nil {
		panic(err)
	}
}
