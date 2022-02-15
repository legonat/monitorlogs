package models

import "github.com/dgrijalva/jwt-go"

type RegisterInputs struct{
	Login 		string `json:"login" binding:"required"`
	Password 	string `json:"password" binding:"required"`
}

type LoginInputs struct{
	Login 		string	`json:"login" binding:"required"`
	Password 	string	`json:"password" binding:"required"`
	RememberMe	bool	`json:"rememberMe"`
	//Fingerprint string `json:"fingerprint" binding:"required"`
}

//type CheckInputs struct {
//	Login 		string `json:"login" binding:"required"`
//	Password 	string `json:"password" binding:"required"`
//	Fingerprint string `json:"fingerprint" binding:"required"`
//}

type BlockInputs struct {
	Login 		string `json:"login" binding:"required"`
	//Ip 			string `json:"ip" binding:"required"`
	//AccessToken string `json:"accessToken" binding:"required"`
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