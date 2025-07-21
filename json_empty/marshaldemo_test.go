package json_empty

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Marshal_Int_Str(t *testing.T) {
	v := Marshal_int_str()
	assert.Equal(t, v, `{"uid":"12345578"}`)
}
