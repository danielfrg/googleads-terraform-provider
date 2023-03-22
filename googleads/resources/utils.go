package resources

import (
	"fmt"
	"math/big"

	"google.golang.org/grpc/status"
)

func ParseClientError(err error) string {
	if e, ok := status.FromError(err); ok {
		return fmt.Sprintf("%s: %s\n%s", e.Code(), e.Message(), e.Details())
	} else {
		return fmt.Sprintf("not able to parse error returned %v", err)
	}
}

func ToBigFloat(val int64) *big.Float {
	return new(big.Float).SetInt(big.NewInt(val))
}
