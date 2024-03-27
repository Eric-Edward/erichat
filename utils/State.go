package utils

const (
	Success = iota
	FailedBindInfo
	FailedBearer
	FailedParseJWT
	FailedNotFoundUser
	FailedExpiredJWT
	FailedGenerateJWT
	FailedRefreshJWT
)
