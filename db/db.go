package db

type UserSecrets struct {
	//Secret string
	Salt string
}
type UserLoginInfo struct {
	UserName  string `json:"userName"`
	HPassword string
}
type UserInfo struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}
type UserInfoForOTP struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}
type UserLoginInfoForOTP struct {
	PhoneNumber       string `json:"phoneNumber"`
	VerificationToken string
}
type UserForOTPAuth struct {
	UserLoginInfoForOTP
	UserInfoForOTP
}
type User struct {
	UserLoginInfo
	UserInfo
	UserSecrets
}

type UserDB interface {
	Insert(u User) (*User, error)
	Find(userName string) (*User, error)
	FindByPhoneNumber(phoneNumber string) (*UserForOTPAuth, error)
	InsertForOTP(u UserForOTPAuth) (*UserForOTPAuth, error)
	UpdateVerificationToken(phoneNumber, VerificationToken string) (int64, error)
}

type DB interface {
	User() UserDB
}
