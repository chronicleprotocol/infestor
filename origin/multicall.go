package origin

import (
	"math/big"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/types"
)

var multicallMethod = abi.MustParseMethod(`
	function aggregate(
		(address target, bytes callData)[] memory calls
	) public returns (
		uint256 blockNumber, 
		bytes[] memory returnData
	)`,
)

type MultiCall struct {
	Target types.Address `abi:"target"`
	Data   []byte        `abi:"callData"`
}

func encodeMultiCallArgs(calls []MultiCall) ([]byte, error) {
	calldata, err := multicallMethod.EncodeArgs(calls)
	if err != nil {
		return nil, err
	}
	return calldata, nil
}

func encodeMultiCallResponse(blockNumber int64, data []any) ([]byte, error) {
	respEncoded, err := abi.EncodeValues(multicallMethod.Outputs(), big.NewInt(blockNumber).Uint64(), data)
	return respEncoded, err
}
