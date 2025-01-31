package dto


type SignupData struct {
	Name	 string `json:"name"`
	Email	 string `json:"email"`
	Password string `json:"password"`
}

type LoginData struct {
	Email	 string `json:"email"`
	Password string `json:"password"`
}


