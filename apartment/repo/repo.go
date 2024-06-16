package repo

import (
	"errors"

	"github.com/pragmataW/apartment_management/dto"
	"github.com/pragmataW/apartment_management/models"
	"gorm.io/gorm"
)

func (r repo) CreateFlat(flatNo int) error {
	var tmp models.Apartment
	result := r.db.Where("flat_no = ?", flatNo).First(&tmp)
	if result.RowsAffected > 0 {
		return dto.FlatAlreadyExists{Message: "flat already exists"}
	}

	apartment := models.Apartment{
		FlatNo: flatNo,
	}

	result = r.db.Create(&apartment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r repo) UpdateFlatOwner(apartment models.Apartment) error {
	result := r.db.Save(&apartment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r repo) DeleteFlat(flatNo int) error {
	result := r.db.Where("flat_no = ?", flatNo).Delete(&models.Apartment{})
	if result.RowsAffected == 0 {
		return dto.ThereIsNoFlat{}
	}

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r repo) GetAllInfoAboutFlat(flatNo int) (models.Apartment, error) {
	flat := models.Apartment{}
	result := r.db.Where("flat_no = ?", flatNo).Take(&flat)
	if result.Error != nil {
		return models.Apartment{}, result.Error
	}

	return flat, nil
}

func (r repo) GetAllInfoAboutAllFlats() ([]models.Apartment, error) {
	var flatList []models.Apartment
	result := r.db.Select("flat_no", "owner_name", "owner_surname", "mail", "dues_count").Find(&flatList)
	if result.Error != nil {
		return nil, result.Error
	}

	return flatList, nil
}

func (r repo) AddDues(flatNo int) error {
	result := r.db.Model(&models.Apartment{}).Where("flat_no = ?", flatNo).UpdateColumn("dues_count", gorm.Expr("dues_count + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r repo) DeleteDues(flatNo int) error {
	existingDues, err := r.GetDuesCount(flatNo)
	if err != nil {
		return err
	}

	if existingDues <= 0 {
		return dto.ThereIsNoDues{}
	}

	result := r.db.Model(&models.Apartment{}).Where("flat_no = ?", flatNo).UpdateColumn("dues_count", gorm.Expr("dues_count - ?", 1))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r repo) AddDuesForAll() error {
	stmt := "UPDATE apartments SET dues_count = dues_count + 1"
	result := r.db.Exec(stmt)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r repo) GetDuesCount(flatNo int) (int, error) {
	flat := models.Apartment{}
	result := r.db.Where("flat_no = ?", flatNo).Select("dues_count").First(&flat)
	if result.Error != nil {
		return -1, result.Error
	}
	return flat.DuesCount, nil
}

func (r repo) DeleteDuesByEmail(email string) error {
	var flat models.Apartment
	result := r.db.Where("mail = ?", email).First(&flat)
	if result.Error != nil {
		return result.Error
	}

	if flat.DuesCount <= 0 {
		return dto.ThereIsNoDues{}
	}

	result = r.db.Model(&models.Apartment{}).Where("mail = ?", email).UpdateColumn("dues_count", gorm.Expr("dues_count - ?", 1))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r repo) GetPasswordAndFlatNoByEmail(email string) (string, int, error) {
	var apartment struct {
		Password string
		FlatNo   int
	}
	result := r.db.Model(&models.Apartment{}).Select("password, flat_no").Where("mail = ?", email).First(&apartment)
	if result.Error != nil {
		return "", 0, result.Error
	}

	return apartment.Password, apartment.FlatNo, nil
}

func (r repo) GetAllAnnouncements() ([]models.Announcement, error) {
	var announcements []models.Announcement
	result := r.db.Find(&announcements)
	if result.Error != nil {
		return nil, result.Error
	}
	return announcements, nil
}

func (r repo) AddAnnouncement(announcement models.Announcement) error {
	result := r.db.Create(&announcement)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r repo) AddMerchant(uuid string, email string) error {
	newMerchant := models.Merchant{
		MerchantID: uuid,
		Email:      email,
	}

	result := r.db.Create(&newMerchant)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r repo) GetEmailFromMerchant(merchantOID string) (string, error) {
	var merchant models.Merchant
	result := r.db.First(&merchant, "merchant_id = ?", merchantOID)
	if result.Error != nil {
		return "", result.Error
	}

	if result.RowsAffected == 0 {
		return "", errors.New("there is no merchant")
	}

	return merchant.Email, nil
}
