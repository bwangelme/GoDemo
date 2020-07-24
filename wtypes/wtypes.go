package wtypes

import "fmt"

type Bool uint8

var (
	Nil   Bool = 0
	True  Bool = 1
	False Bool = 2
)

func (b Bool) String() string {
	if b == True {
		return "true"
	} else if b == False {
		return "false"
	} else {
		return "null"
	}
}

func (b *Bool) ToBoolean() bool {
	if *b == True {
		return true
	}

	return false
}

func (b *Bool) MarshalJSON() ([]byte, error) {
	if *b == True {
		return []byte("true"), nil
	} else if *b == False {
		return []byte("false"), nil
	} else {
		// 注意：空值应该忽略的，不应该走到 Marshal 这一步
		return nil, fmt.Errorf("should add `,omitempty` on Bool value")
	}
}

func (b *Bool) UnmarshalJSON(data []byte) error {
	if string(data) == `true` {
		*b = True
	} else if string(data) == `false` {
		*b = False
	} else {
		return fmt.Errorf("invalid input %s", string(data))
	}

	return nil
}
