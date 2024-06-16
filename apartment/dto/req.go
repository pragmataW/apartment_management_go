package dto

type PaymentSendReq struct {
	MerchantId     string `json:"merchant_id"`
	MerchantOID    string `json:"merchant_oid"`
	PaymentAmount  string `json:"payment_amount"`
	Currency       string `json:"currency"`
	Email          string `json:"email"`
	UserName       string `json:"user_name"`
	UserAddress    string `json:"user_address"`
	UserPhone      string `json:"user_phone"`
	OkURL          string `json:"merchant_ok_url"`
	FailURL        string `json:"merchant_fail_url"`
	UserBasket     string `json:"user_basket"`
	UserIP         string `json:"user_ip"`
	TimeOutLimit   string `json:"timeout_limit"`
	DebugOn        string `json:"debug_on"`
	TestMode       string `json:"test_mode"`
	NoInstallment  string `json:"no_installment"`
	MaxInstallment string `json:"max_installment"`
	PaytrToken     string `json:"paytr_token"`
}

type PaymentGetReq struct {
	Email       string          `json:"email" validate:"required,email"`
	UserName    string          `json:"user_name" validate:"required"`
	UserAddress string          `json:"user_address" validate:"required"`
	UserPhone   string          `json:"user_phone" validate:"required"`
	UserBasket  [][]interface{} `json:"user_basket" validate:"required"`
	DebugOn     string          `json:"debug_on" validate:"required"`
	TestMode    string          `json:"test_mode" validate:"required"`
}

type LoginAdminReq struct {
	Password string `json:"password" validate:"required"`
}

type LoginUserReq struct {
	FlatNo   int    `json:"flat_no" validate:"required"`
	Mail     string `json:"mail" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ChangeDuesPriceReq struct {
	Price float64 `json:"price" validate:"required"`
}

type ChangePayDayReq struct {
	PayDay int `json:"payday" validate:"required"`
}

type ApartmentRequest struct {
	FlatNo       int    `json:"flat_no" validate:"required"`
	OwnerName    string `json:"owner_name" validate:"required,min=2"`
	OwnerSurname string `json:"owner_surname" validate:"required,min=2"`
	Mail         string `json:"mail" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=8"`
	DuesCount    int    `json:"dues_count"`
}

type AnnouncementRequest struct {
	AnnouncementID int    `json:"announcement_id"`
	Title          string `json:"title" validate:"required"`
	Content        string `json:"content" validate:"required"`
}

type MailRequest struct {
	ToMail  string `json:"to_mail" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
}
