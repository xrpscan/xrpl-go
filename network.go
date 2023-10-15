package xrpl

import (
	"strconv"
)

const XRPL_NATIVE_ASSET = "XRP"
const XAHAU_NATIVE_ASSET = "XAH"

type Network int32

const (
	NetworkXrplMainnet   Network = 0
	NetworkXrplTestnet   Network = 1
	NetworkXrplDevnet    Network = 2
	NetworkXrplAmmDevnet Network = 25
	NetworkXahauMainnet  Network = 21337
	NetworkXahauTestnet  Network = 21338
)

func (n Network) Asset() string {
	switch n {
	// XRPL networks
	case NetworkXrplMainnet:
		return XRPL_NATIVE_ASSET
	case NetworkXrplTestnet:
		return XRPL_NATIVE_ASSET
	case NetworkXrplDevnet:
		return XRPL_NATIVE_ASSET
	case NetworkXrplAmmDevnet:
		return XRPL_NATIVE_ASSET

	// Xahau networks
	case NetworkXahauMainnet:
		return XAHAU_NATIVE_ASSET
	case NetworkXahauTestnet:
		return XAHAU_NATIVE_ASSET

	// Default is XRPL network
	default:
		return XRPL_NATIVE_ASSET
	}
}

func (n Network) Name() string {
	switch n {
	// XRPL networks
	case NetworkXrplMainnet:
		return "XrplMainnet"
	case NetworkXrplTestnet:
		return "XrplTestnet"
	case NetworkXrplDevnet:
		return "XrplDevnet"
	case NetworkXrplAmmDevnet:
		return "XrplAmmDevnet"

	// Xahau networks
	case NetworkXahauMainnet:
		return "XahauMainnet"
	case NetworkXahauTestnet:
		return "XahauTestnet"

	// Default
	default:
		return strconv.Itoa(int(n))
	}
}

func GetNetwork(networkId int) Network {
	switch networkId {
	// XRPL networks
	case int(NetworkXrplMainnet):
		return NetworkXrplMainnet
	case int(NetworkXrplTestnet):
		return NetworkXrplTestnet
	case int(NetworkXrplDevnet):
		return NetworkXrplDevnet
	case int(NetworkXrplAmmDevnet):
		return NetworkXrplAmmDevnet

	// Xahau networks
	case int(NetworkXahauMainnet):
		return NetworkXahauMainnet
	case int(NetworkXahauTestnet):
		return NetworkXahauTestnet

	// Default is XRPL network
	default:
		return NetworkXrplMainnet
	}
}
