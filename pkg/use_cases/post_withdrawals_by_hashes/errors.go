package post_withdrawals_by_hashes

import "errors"

// ErrUCInputEmpty error: uc-input must not be empty.
var ErrUCInputEmpty = errors.New("uc-input must not be empty")

// ErrWithdrawalByHashesFail error: failed to get withdrawals by hashes.
var ErrWithdrawalByHashesFail = errors.New("failed to get withdrawals by hashes")
