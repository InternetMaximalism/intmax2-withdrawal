package mnemonic_wallet

import "intmax2-withdrawal/internal/mnemonic_wallet/models"

type MnemonicWallet interface {
	WalletGenerator(mnemonicDerivationPath, password string) (w *models.Wallet, err error)
	WalletFromMnemonic(
		mnemonic, password, mnemonicDerivationPath string,
	) (w *models.Wallet, err error)
	WalletFromPrivateKeyHex(
		privateKeyHex string,
	) (w *models.Wallet, err error)
}
