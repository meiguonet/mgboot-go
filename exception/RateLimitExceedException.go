package BuiltinException

type RateLimitExceedException struct {
	total      int
	remaining  int
	retryAfter string
}

func NewRateLimitExceedException(data map[string]interface{}) RateLimitExceedException {
	var total int

	if n1, ok := data["total"].(int); ok && n1 > 0 {
		total = n1
	}

	var remaining int

	if n1, ok := data["remaining"].(int); ok && n1 > 0 {
		remaining = n1
	}

	var retryAfter string

	if s1, ok := data["retryAfter"].(string); ok && s1 != "" {
		retryAfter = s1
	}

	return RateLimitExceedException{
		total:      total,
		remaining:  remaining,
		retryAfter: retryAfter,
	}
}

func (ex RateLimitExceedException) Error() string {
	return "rate limit exceed"
}

func (ex RateLimitExceedException) Total() int {
	return ex.total
}

func (ex RateLimitExceedException) Remaining() int {
	return ex.remaining
}

func (ex RateLimitExceedException) RetryAfter() string {
	return ex.retryAfter
}
