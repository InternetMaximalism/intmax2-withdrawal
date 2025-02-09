package utils

import (
	"fmt"
	"intmax2-withdrawal/internal/logger"
	"intmax2-withdrawal/internal/mnemonic_wallet"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func NewClient(url string) (*ethclient.Client, error) {
	ethClient, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("error connecting to rpc service: %+v", err)
		return nil, err
	}
	defer ethClient.Close()
	return ethClient, nil
}

func CreateTransactor(ethereumPrivateKeyHex, networkChainID string) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(ethereumPrivateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}
	const (
		int10Key = 10
		int64Key = 64
	)
	chainID, err := strconv.ParseInt(networkChainID, int10Key, int64Key)
	if err != nil {
		return nil, fmt.Errorf("invalid chain ID: %w", err)
	}
	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(chainID))
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}
	if transactOpts == nil {
		return nil, fmt.Errorf("transactOpts is nil")
	}
	return transactOpts, nil
}

func IsValidEthereumPrivateKey(key string) error {
	key = strings.TrimPrefix(key, "0x")
	_, err := crypto.HexToECDSA(key)
	if err != nil {
		return fmt.Errorf("invalid ethereum private key: %w", err)
	}
	return nil
}

func PrivateKeyToAddress(pkHex string) (*common.Address, error) {
	if pkHex == "" {
		return nil, fmt.Errorf("private key cannot be empty")
	}
	wallet, err := mnemonic_wallet.New().WalletFromPrivateKeyHex(pkHex)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet from private key: %w", err)
	}
	if wallet.WalletAddress == nil {
		return nil, fmt.Errorf("wallet address is nil")
	}
	return wallet.WalletAddress, nil
}

func LogTransactionDebugInfo(_log logger.Logger, privateKeyHex, contractAddres string, args ...interface{}) error {
	address, err := PrivateKeyToAddress(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to get address from private key: %w", err)
	}
	_log.Debugf("Contract address: %s", contractAddres)
	_log.Debugf("Transaction sender address: %s", address.Hex())
	for i, arg := range args {
		_log.Debugf("Transaction Argument %d: %v (Type: %T)", i, arg, arg)
	}
	return nil
}
