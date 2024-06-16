package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/pragmataW/apartment_management/dto"
	"github.com/pragmataW/apartment_management/pkg/jwt"
	"github.com/robfig/cron/v3"
)

func (s *service) LoginAdmin(password string) (string, error) {
	if s.ConfigManager.GetAdminPassword() == password {
		claim := jwt.JwtClaim{
			FlatNo: -1,
			Role:   "admin",
			Exp:    time.Now().Add(72 * time.Hour).Unix(),
		}

		jwtGenerator := jwt.NewJwtGenerator(claim, s.ConfigManager.GetJwtKey())

		token, err := jwtGenerator.GenerateJWT()
		if err != nil {
			return "", err
		}
		return token, nil
	}

	return "", dto.PasswordMatchError{
		Message: "password does not match",
	}
}

func (s *service) LoginUser(flatNo int, mail string, password string) (string, error) {
	passwordDb, flaNoDb, err := s.Repo.GetPasswordAndFlatNoByEmail(mail)
	if err != nil {
		return "", err
	}

	passwordDb, err = s.Encryptor.Decrypt(passwordDb)
	if err != nil {
		return "", err
	}

	if flaNoDb != flatNo || passwordDb != password {
		return "", dto.UserDoesNotExists{Message: "user does not exists"}
	}

	claim := jwt.JwtClaim{
		FlatNo: flatNo,
		Role:   "user",
		Exp:    time.Now().Add(72 * time.Hour).Unix(),
		Email:  mail,
	}

	jwtGenerator := jwt.NewJwtGenerator(claim, s.ConfigManager.GetJwtKey())

	token, err := jwtGenerator.GenerateJWT()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) CreateFlat(flatNo int) error {
	err := s.Repo.CreateFlat(flatNo)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateFlatOwner(apartment Apartment) error {
	apartmentModel := apartment.ToApartmentModel()

	var err error
	apartmentModel.Password, err = s.Encryptor.Encrypt(apartmentModel.Password)
	if err != nil {
		return err
	}

	if err := s.Repo.UpdateFlatOwner(apartmentModel); err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteFlat(flatNo int) error {
	err := s.Repo.DeleteFlat(flatNo)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetAllInfoAboutFlat(flatNo int) (Apartment, error) {
	modelApartments, err := s.Repo.GetAllInfoAboutFlat(flatNo)
	if err != nil {
		return Apartment{}, nil
	}

	modelApartments.Password, err = s.Encryptor.Decrypt(modelApartments.Password)
	if err != nil {
		return Apartment{}, err
	}

	var apartment Apartment
	apartment.ToApartmentServiceObject(modelApartments)

	return apartment, nil
}

func (s *service) GetAllInfoAboutAllFlat() ([]Apartment, error) {
	modelApartments, err := s.Repo.GetAllInfoAboutAllFlats()
	if err != nil {
		return []Apartment{}, err
	}

	var apartments []Apartment
	for _, modelApartment := range modelApartments {
		apartment := Apartment{}
		apartment.ToApartmentServiceObject(modelApartment)
		apartments = append(apartments, apartment)
	}
	return apartments, nil
}

func (s *service) AddDues(flatNo int) error {
	err := s.Repo.AddDues(flatNo)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteDues(flatNo int) error {
	err := s.Repo.DeleteDues(flatNo)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ChangeDuesPrice(price float64) {
	dto.Mutx.Lock()
	dto.DuesPrice = price
	dto.Mutx.Unlock()
}

func (s *service) ChangePayDay(payDay int) error {
	if payDay < 1 || payDay > 28 {
		return dto.PayDayRangeError{Message: "invalid range"}
	}
	dto.Mutx.Lock()
	dto.PayDay = payDay
	dto.Mutx.Unlock()

	return nil
}

func (s *service) AddAnnouncement(announcement Announcement) error {
	announcementModel := announcement.ToAnnouncementModel()
	err := s.Repo.AddAnnouncement(announcementModel)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetAllAnnouncements() ([]Announcement, error) {
	modelAnnouncements, err := s.Repo.GetAllAnnouncements()
	if err != nil {
		return []Announcement{}, err
	}

	var announcements []Announcement
	for _, announcement := range modelAnnouncements {
		a := Announcement{}
		a.ToAnnouncementServiceObject(announcement)
		announcements = append(announcements, a)
	}
	return announcements, nil
}

func (s *service) SendMail(subject string, body string, mail string) error {
	reqBody := map[string]interface{}{
		"from_name":  "apartment",
		"from_email": s.ConfigManager.GetFromMail(),
		"to_name":    "-",
		"to_email":   mail,
		"subject":    subject,
		"html":       body,
	}
	log.Println(reqBody)

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.ConfigManager.GetMailServer(), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.SendMailError{Message: "send mail error"}
	}

	return nil
}

func (s *service) GetPaymentToken(payment dto.PaymentSendReq) (string, error) {
	params := url.Values{
		"merchant_id":       {payment.MerchantId},
		"user_ip":           {payment.UserIP},
		"merchant_oid":      {payment.MerchantOID},
		"email":             {payment.Email},
		"payment_amount":    {payment.PaymentAmount},
		"paytr_token":       {payment.PaytrToken},
		"user_basket":       {payment.UserBasket},
		"debug_on":          {payment.DebugOn},
		"no_installment":    {payment.NoInstallment},
		"max_installment":   {payment.MaxInstallment},
		"user_name":         {payment.UserName},
		"user_address":      {payment.UserAddress},
		"user_phone":        {payment.UserPhone},
		"merchant_ok_url":   {payment.OkURL},
		"merchant_fail_url": {payment.FailURL},
		"timeout_limit":     {payment.TimeOutLimit},
		"currency":          {payment.Currency},
		"test_mode":         {payment.TestMode},
	}

	resp, err := http.PostForm("https://www.paytr.com/odeme/api/get-token", params)
	if err != nil {
		return "", fmt.Errorf("post request error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body error: %v", err)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		return "", fmt.Errorf("decode response body error: %v", err)
	}

	if status, ok := res["status"].(string); ok && status == "success" {
		if token, ok := res["token"].(string); ok {
			err := s.Repo.AddMerchant(payment.MerchantOID, payment.Email)
			if err != nil{
				return "", err
			}
			return token, nil
		} else {
			return "", fmt.Errorf("there is no token in response")
		}
	} else {
		return "", fmt.Errorf("status not ok: %s", string(body))
	}
}

func (s *service) PaymentCallback(merchantOID string) error {
	email, err := s.Repo.GetEmailFromMerchant(merchantOID)
	if err != nil {
		return err
	}

	err = s.Repo.DeleteDuesByEmail(email)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) IncreaseDuesAutomatically() error {
	cronTab := cron.New()

	_, err := cronTab.AddFunc(fmt.Sprintf("0 0 %d * *", dto.PayDay), func() {
		err := s.Repo.AddDuesForAll()
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		return err
	}
	cronTab.Start()
	select {}
}

