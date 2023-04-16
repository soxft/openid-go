package thirdpartutil

type GithubInterface interface {
	thirdPart
}

type Github struct {
	ThirdPartCtx
}

func NewThirdPartGithub(ctx ThirdPartCtx) *Github {
	return &Github{
		ThirdPartCtx: ctx,
	}
}

func GetUserInfo() (interface{}, error) {
	return nil, nil
}

func RedirectUrl() string {
	return ""
}
