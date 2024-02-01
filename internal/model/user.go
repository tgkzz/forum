package model

type User struct {
	Id         int
	Username   string
	Email      string
	Password   string
	AuthMethod string
}

type GithubUser struct {
	Login  string `json:"login"`
	NodeId string `json:"node_id"`
}

type GoogleUsers struct {
	Sub      string `json:"sub"`
	Username string `json:"name"`
}
