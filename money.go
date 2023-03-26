package swiss_qr_code

import "fmt"

type MoneyValue struct {
	Base  uint32
	Cents uint32
}

func (m MoneyValue) String() string {
	return fmt.Sprintf("%d.%02d", m.Base, m.Cents)
}
