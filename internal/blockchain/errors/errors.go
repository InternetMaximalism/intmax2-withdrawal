package errors

import "errors"

// ErrScrollChainIDInvalidStr error: the scroll chain ID must be equal: %s, %s.
const ErrScrollChainIDInvalidStr = "the scroll chain ID must be equal: %s, %s"

// ErrScrollNetworkChainLinkEvmJSONRPCFail error: failed to get the chain-link-evm-json-rpc of scroll network.
var ErrScrollNetworkChainLinkEvmJSONRPCFail = errors.New(
	"failed to get the chain-link-evm-json-rpc of scroll network",
)
