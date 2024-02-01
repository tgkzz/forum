package model

const (
	GOOGLECLIENTID     = "957724082433-c753taa2q7cahruqnifsedmou76pc2ua.apps.googleusercontent.com"
	GOOGLECLIENTSECRET = "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"
)

type GoogleRespBody struct {
	AccessToken string `json:"access_token"`
}
