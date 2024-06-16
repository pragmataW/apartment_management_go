package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pragmataW/apartment_management/dto"
	randomkeygen "github.com/pragmataW/apartment_management/pkg/random_keygen"
	"github.com/pragmataW/apartment_management/services"
)

func (ctrl *controller) LoginAdmin(c *fiber.Ctx) error {
	var body dto.LoginAdminReq
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request",
		})
	}

	token, err := ctrl.Service.LoginAdmin(body.Password)
	if err != nil {
		if err, ok := err.(dto.PasswordMatchError); ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "Authentication",
		Value:    token,
		HTTPOnly: true,
		Expires:  time.Now().Add(72 * time.Hour),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "status ok",
	})
}

func (ctrl *controller) LoginUser(c *fiber.Ctx) error {
	var body dto.LoginUserReq
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request",
		})
	}

	token, err := ctrl.Service.LoginUser(body.FlatNo, body.Mail, body.Password)
	if err != nil {
		if _, ok := err.(dto.PasswordMatchError); ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid credentials",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})

	}

	c.Cookie(&fiber.Cookie{
		Name:     "Authentication",
		Value:    token,
		HTTPOnly: true,
		Expires:  time.Now().Add(72 * time.Hour),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "status ok",
	})
}

func (ctrl *controller) Logout(c *fiber.Ctx) error {

	c.Cookie(&fiber.Cookie{
		Name:     "Authentication",
		Value:    "",
		HTTPOnly: true,
		Expires:  time.Now().Add(-1 * time.Hour),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "logout successful",
	})
}

func (ctrl *controller) CreateFlat(c *fiber.Ctx) error {
	flatNo, err := strconv.Atoi(c.Params("flatNo"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "missing query parameter: flatNo - " + strconv.Itoa(flatNo),
		})
	}

	err = ctrl.Service.CreateFlat(flatNo)
	if err != nil {
		if err, ok := err.(dto.FlatAlreadyExists); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "status ok",
	})
}

func (ctrl *controller) UpdateFlatOwner(c *fiber.Ctx) error {
	var body dto.ApartmentRequest
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request " + err.Error(),
		})
	}

	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	apartment := services.Apartment{
		FlatNo:       body.FlatNo,
		OwnerName:    body.OwnerName,
		OwnerSurname: body.OwnerSurname,
		Mail:         body.Mail,
		Password:     body.Password,
		DuesCount:    body.DuesCount,
	}

	err = ctrl.Service.UpdateFlatOwner(apartment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "status ok",
	})
}

func (ctrl *controller) DeleteFlat(c *fiber.Ctx) error {
	flatNoParam := c.Params("flatNo")
	if flatNoParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "missing query parameter: flatNo",
		})
	}

	flatNo, err := strconv.Atoi(flatNoParam)
	if err != nil || flatNo <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid query parameter: flatNo",
		})
	}

	err = ctrl.Service.DeleteFlat(flatNo)
	if err != nil {
		if err, ok := err.(dto.ThereIsNoFlat); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "status ok",
	})
}

func (ctrl *controller) GetAllInfoAboutFlat(c *fiber.Ctx) error {
	flatNo, err := strconv.Atoi(c.Params("flatNo"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "missing query parameter: flatNo - " + strconv.Itoa(flatNo),
		})
	}

	apartment, err := ctrl.Service.GetAllInfoAboutFlat(flatNo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := dto.ApartmentResponse{
		FlatNo:       apartment.FlatNo,
		OwnerName:    apartment.OwnerName,
		OwnerSurname: apartment.OwnerSurname,
		Mail:         apartment.Mail,
		Password:     apartment.Password,
		DuesCount:    apartment.DuesCount,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) GetAllInfoAboutAllFlat(c *fiber.Ctx) error {
	apartments, err := ctrl.Service.GetAllInfoAboutAllFlat()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var resp []dto.ApartmentResponse
	for _, apartment := range apartments {
		respApartment := dto.ApartmentResponse{
			FlatNo:       apartment.FlatNo,
			OwnerName:    apartment.OwnerName,
			OwnerSurname: apartment.OwnerSurname,
			Mail:         apartment.Mail,
			Password:     apartment.Password,
			DuesCount:    apartment.DuesCount,
		}
		resp = append(resp, respApartment)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) AddDues(c *fiber.Ctx) error {
	flatNo, err := strconv.Atoi(c.Params("flatNo"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "missing query parameter: flatNo - " + strconv.Itoa(flatNo),
		})
	}

	if err := ctrl.Service.AddDues(flatNo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "status ok",
	})
}

func (ctrl *controller) DeleteDues(c *fiber.Ctx) error {
	flatNo, err := strconv.Atoi(c.Params("flatNo"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "missing query parameter: flatNo - " + strconv.Itoa(flatNo),
		})
	}

	if err := ctrl.Service.DeleteDues(flatNo); err != nil {
		if err, ok := err.(dto.ThereIsNoDues); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "status ok",
	})
}

func (ctrl *controller) ChangeDuesPrice(c *fiber.Ctx) error {
	var body dto.ChangeDuesPriceReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request " + err.Error(),
		})
	}

	ctrl.Service.ChangeDuesPrice(body.Price)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "status ok",
	})
}

func (ctrl *controller) ChangePayDay(c *fiber.Ctx) error {
	var body dto.ChangePayDayReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request " + err.Error(),
		})
	}

	if err := ctrl.Service.ChangePayDay(body.PayDay); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "status ok",
	})
}

func (ctrl *controller) GetDuesPrice(c *fiber.Ctx) error {
	resp := dto.DuesPriceResponse{
		DuesPrice: dto.DuesPrice,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) GetPayDay(c *fiber.Ctx) error {
	resp := dto.PayDayResponse{
		PayDay: dto.PayDay,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) AddAnnouncement(c *fiber.Ctx) error {
	var body dto.AnnouncementRequest
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request " + err.Error(),
		})
	}

	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	serviceAnnouncement := services.Announcement{
		AnnouncementID: body.AnnouncementID,
		Title:          body.Title,
		Content:        body.Content,
	}

	err = ctrl.Service.AddAnnouncement(serviceAnnouncement)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "status ok",
	})
}

func (ctrl *controller) GetAllAnnouncements(c *fiber.Ctx) error {
	announcements, err := ctrl.Service.GetAllAnnouncements()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var respAnnouncements []dto.AnnouncementResponse
	for _, announcement := range announcements {
		newAnnouncement := dto.AnnouncementResponse{
			AnnouncementID: announcement.AnnouncementID,
			Title:          announcement.Title,
			Content:        announcement.Content,
		}
		respAnnouncements = append(respAnnouncements, newAnnouncement)
	}

	return c.Status(fiber.StatusOK).JSON(respAnnouncements)
}

func (ctrl *controller) SendMail(c *fiber.Ctx) error {
	var body dto.MailRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := ctrl.Service.SendMail(body.Subject, body.Body, body.ToMail); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "status ok",
	})
}

func (ctrl *controller) GetPaymentToken(c *fiber.Ctx) error {
	var body dto.PaymentGetReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request " + err.Error(),
		})
	}

	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request " + err.Error(),
		})
	}

	merchantID := strconv.Itoa(ctrl.ConfigManager.GetMerchantID())
	merchantKey := []byte(ctrl.ConfigManager.GetMerchantKey())
	merchantSalt := []byte(ctrl.ConfigManager.GetMerchantSalt())

	merchantOid := randomkeygen.NewKeygen(64).GenerateRandomKey()
	fmt.Println(merchantOid)
	email := c.Locals("email").(string)
	paymentAmount := strconv.Itoa(int(dto.DuesPrice * 100))
	userName := body.UserName
	userAddress := body.UserAddress
	userPhone := body.UserPhone
	merchantOkUrl := ctrl.ConfigManager.GetOkUrl()
	merchantFailUrl := ctrl.ConfigManager.GetFailUrl()

	userBasket := body.UserBasket
	userBasketJSON, _ := json.Marshal(userBasket)
	userBasketEncoded := base64.StdEncoding.EncodeToString(userBasketJSON)

	userIP := c.IP()
	timeOutLimit := "30"
	debugOn := body.DebugOn
	testMode := body.TestMode
	noInstallment := "1"
	maxInstallment := "0"
	currency := "TL"

	hashStr := strings.Join([]string{
		merchantID, userIP, merchantOid, email, paymentAmount,
		userBasketEncoded, noInstallment, maxInstallment, currency, testMode,
	}, "")
	h := hmac.New(sha256.New, merchantKey)
	h.Write([]byte(hashStr))
	h.Write(merchantSalt)
	paytrToken := base64.StdEncoding.EncodeToString(h.Sum(nil))

	req := dto.PaymentSendReq{
		MerchantId:     merchantID,
		UserIP:         userIP,
		MerchantOID:    merchantOid,
		Email:          email,
		PaymentAmount:  paymentAmount,
		PaytrToken:     paytrToken,
		UserBasket:     userBasketEncoded,
		DebugOn:        debugOn,
		NoInstallment:  noInstallment,
		MaxInstallment: maxInstallment,
		UserName:       userName,
		UserAddress:    userAddress,
		UserPhone:      userPhone,
		OkURL:          merchantOkUrl,
		FailURL:        merchantFailUrl,
		TimeOutLimit:   timeOutLimit,
		Currency:       currency,
		TestMode:       testMode,
	}

	token, err := ctrl.Service.GetPaymentToken(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}

func (ctrl *controller) PaymentCallback(c *fiber.Ctx) error {
	merchantKey := ctrl.ConfigManager.GetMerchantKey()
	merchantSalt := ctrl.ConfigManager.GetMerchantSalt()

	merchantOid := c.FormValue("merchant_oid")
	status := c.FormValue("status")
	totalAmount := c.FormValue("total_amount")
	receivedHash := c.FormValue("hash")

	hashStr := merchantOid + merchantSalt + status + totalAmount
	h := hmac.New(sha256.New, []byte(merchantKey))
	h.Write([]byte(hashStr))
	expectedHash := base64.StdEncoding.EncodeToString(h.Sum(nil))

	if receivedHash != expectedHash {
		fmt.Println("bad hash")
		return c.SendString("PAYTR notification failed: bad hash")
	}

	if status == "success" {
		err := ctrl.Service.PaymentCallback(merchantOid)
		if err != nil {
			fmt.Println("service error")
			return c.SendString(err.Error())
		}
		fmt.Println("payment done")
	} else {
		fmt.Println("payment is not done")
		return c.SendString(status)
	}

	return c.SendString("OK")
}
