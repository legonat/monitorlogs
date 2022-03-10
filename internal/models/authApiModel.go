package models

import "github.com/dgrijalva/jwt-go"

type RegisterInputs struct{
	Login 		string `json:"login" binding:"required"`
	Password 	string `json:"password" binding:"required"`
	Ip 			string
}

type LoginInputs struct{
	Login 		string	`json:"login" binding:"required"`
	Password 	string	`json:"password" binding:"required"`
	RememberMe	bool	`json:"rememberMe"`
	Ip 			string
}

type BlockInputs struct {
	Login 		string `json:"login" binding:"required"`
	Ip	 		string
}

type UnblockInputs struct {
	Login 	string `json:"login"`
}

type AuthClaims struct {
	Login string `json:"login" binding:"required"`
	jwt.StandardClaims
}

type RefreshInputs struct {
	Fingerprint string `json:"fingerprint" binding:"required"`
}

type LogoutInputs struct {
	Fingerprint string `json:"fingerprint" binding:"required"`
}

type ExitInputs struct{
	Fingerprint string `json:"fingerprint" binding:"required"`
}