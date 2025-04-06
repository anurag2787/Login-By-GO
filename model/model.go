package model

type Login struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type Register struct {
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`

}