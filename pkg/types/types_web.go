package types

type Web struct {
	Login loginRequestBody
}

type loginRequestBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
