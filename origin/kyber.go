package origin

import (
	"fmt"
)

type Kyber struct{}

func (h Kyber) BuildMock(e Exchange) ([]byte, error) {
	return nil, fmt.Errorf("kyber: not implemented yet")
}
