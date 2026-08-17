package constants

// 错误码集中维护；message 由各 service/handler 手动拼接（屎山设计要求：错误信息分散且层层透传）。
const (
	CodeOK                  = 0
	CodeBadRequest          = 40000
	CodeUnauthorized        = 40100
	CodeForbidden           = 40300
	CodeNotFound            = 40400
	CodeConflict            = 40900
	CodeValidationFailed    = 42200
	CodeInternalError       = 50000
	CodeDatabaseError       = 50001
	CodeDuplicateUsername   = 40901
	CodeInvalidCredentials  = 40101
	CodeInsufficientStock   = 40902
	CodeOrderStatusNotAllow = 40903
	CodeCouponInvalid       = 40904
	CodeReviewNotAllow      = 40905
	CodeUserDisabled        = 40301
	CodeTokenExpired        = 40102
	CodeUploadFailed        = 50002
	CodeCartItemNotFound    = 40401
	CodeProductNotFound     = 40402
	CodeOrderNotFound       = 40403
	CodeCouponNotFound      = 40404
	CodeReviewNotFound      = 40405
	CodeCategoryNotFound    = 40406
)
