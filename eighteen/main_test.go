package eighteen

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_IsOver18(t *testing.T) {
	for _, tt := range []struct {
		Birthday string
		Now_     string
		Wanted   bool
	}{
		{
			"2006-03-01",
			"2024-02-29",
			false,
		},
	} {
		birthday, _ := time.Parse("2006-01-02", tt.Birthday)
		now, _ := time.Parse("2006-01-02", tt.Now_)
		assert.Equal(t, tt.Wanted, IsOver18(birthday, now))
	}
}
