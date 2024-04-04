package utils

const (
	Success = iota
	FailedBindInfo
	FailedBearer
	FailedParseJWT
	FailedFindUser
	FailedExpiredJWT
	FailedGenerateJWT
	FailedRefreshJWT
	FailedCreateChatRoom
	FailedFindChatRoom
	FailedGenerateSocket
	FailedReadMessage
	FailedLoadFriends
)
