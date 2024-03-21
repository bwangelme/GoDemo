package eighteen

import (
	"fmt"
	"time"
)

func IsOver18(birthday time.Time, now_ time.Time) bool {
	ago18 := now_.AddDate(-18, 0, 0)
	fmt.Println(birthday, ago18)
	return birthday == ago18 || birthday.Before(ago18)
}
