package github

type AccessTokenStruct struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`

	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	// ErrorUri         string `json:"error_uri"`
}

type UserStruct struct {
	Login string `json:"login"`
	ID    int    `json:"id"`

	Message string `json:"message"`
}
