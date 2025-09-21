package auth

type LoginResponse struct {
	RefreshToken string `json:"refresh"`
	AccessToken  string `json:"access"`
	User         struct {
		ID          string `json:"uuid"`
		Username    string `json:"username"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		IsActive    bool   `json:"is_active"`
		IsStaff     bool   `json:"is_staff"`
		IsSuperuser bool   `json:"is_superuser"`
	} `json:"user"`
}

type UserInfo struct {
	ID            int    `json:"id"`
	UUID          string `json:"uuid"`
	Username      string `json:"username"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	IsActive      bool   `json:"is_active"`
	IsStaff       bool   `json:"is_staff"`
	IsSuperuser   bool   `json:"is_superuser"`
	S3AccessToken string `json:"s3_access_token"`
}
