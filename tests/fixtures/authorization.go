package fixtures

type AuthorizationCase struct {
	Name    string
	Token   string
	Message string
}

// AuthorizationCases задает общий набор отсутствующего и некорректного токенов.
func AuthorizationCases() []AuthorizationCase {
	return []AuthorizationCase{
		{Name: "missing_token", Message: "Authorization header required"},
		{Name: "invalid_token", Token: InvalidToken, Message: "Invalid or expired token"},
	}
}
