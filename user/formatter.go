package user

type UserFormatter struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Email      string `json:"email"`
	Token      string `json:"token"`
	ImageUrl   string `json:"image_url"`
}
type EmailCheckFormatter struct {
	Validate bool `json:"validate"`
}

func FormatUser(user User, token string) UserFormatter {

	formatter := UserFormatter{
		ID:         user.ID,
		Name:       user.Name,
		Occupation: user.Occupation,
		Email:      user.Email,
		Token:      token,
		ImageUrl:   user.AvatarFileName,
	}
	return formatter
}

func CheckEmail(validate bool) EmailCheckFormatter {

	formatter := EmailCheckFormatter{
		Validate: validate,
	}

	return formatter
}
