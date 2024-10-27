package third_party

import (
	"embed"
)

//go:embed OpenAPI/withdrawal_service/*
var OpenAPIWithdrawal embed.FS
