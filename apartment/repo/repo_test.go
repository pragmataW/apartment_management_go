package repo

import (
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/pragmataW/apartment_management/dto"
	"github.com/pragmataW/apartment_management/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dsn = "host=localhost user=postgres password=123wsedrf dbname=Apartments port=5432 sslmode=disable"
)

func setupDb(model interface{}) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")

	migrator := db.Migrator()
	migrator.CreateTable(model)

	return db
}

func TestCreateFlat(t *testing.T) {
	db := setupDb(models.Apartment{})

	repo := NewRepo(db)

	flatNo := 124

	err := repo.CreateFlat(flatNo)
	assert.NoError(t, err)

	var apartment models.Apartment
	result := db.First(&apartment, "flat_no = ?", flatNo)
	assert.NoError(t, result.Error)
	assert.Equal(t, flatNo, apartment.FlatNo)
}

func TestCreateFlatForError(t *testing.T) {
	db := setupDb(models.Apartment{})

	repo := NewRepo(db)

	flatNo := 124

	err := repo.CreateFlat(flatNo)
	assert.NoError(t, err)

	err = repo.CreateFlat(flatNo)
	assert.Error(t, err)
	assert.IsType(t, dto.FlatAlreadyExists{}, err)
}

func TestUpdateFlatOwner(t *testing.T) {
	db := setupDb(models.Apartment{})

	repo := NewRepo(db)
	expected := models.Apartment{
		FlatNo:       1,
		OwnerName:    "Yusuf",
		OwnerSurname: "Çiftçi",
		Mail:         "yciftci@gmail.com",
		Password:     "123",
		DuesCount:    0,
	}

	err := repo.CreateFlat(1)
	assert.NoError(t, err)

	err = repo.UpdateFlatOwner(expected)
	assert.NoError(t, err)

	var apartment models.Apartment
	result := db.First(&apartment, "flat_no = ?", 1)
	assert.NoError(t, result.Error)

	assert.Equal(t, expected, apartment)
}

func TestDeleteFlat(t *testing.T) {
	db := setupDb(models.Apartment{})
	repo := NewRepo(db)

	apartment := models.Apartment{
		FlatNo: 1,
	}

	result := db.Create(&apartment)
	assert.NoError(t, result.Error)

	err := repo.DeleteFlat(1)
	assert.NoError(t, err)

	var tmp models.Apartment
	result = db.First(&tmp, "flat_no = ?", 1)
	assert.Error(t, result.Error)
	assert.Equal(t, result.RowsAffected, int64(0))
}

func TestDeleteFlatForError(t *testing.T) {
	db := setupDb(models.Apartment{})
	repo := NewRepo(db)

	err := repo.DeleteFlat(32)
	assert.Error(t, err)
	assert.IsType(t, dto.ThereIsNoFlat{}, err)
}

func TestGetAllInfoAboutFlat(t *testing.T) {
	db := setupDb(models.Apartment{})
	repo := NewRepo(db)

	flatNo := 1
	ownerName := "yusuf"
	ownerSurname := "çiftçi"
	ownerMail := "yciftci@gmail.com"
	ownerPassword := "123"
	duesCount := 0

	apartment := models.Apartment{
		FlatNo:       flatNo,
		OwnerName:    ownerName,
		OwnerSurname: ownerSurname,
		Mail:         ownerMail,
		Password:     ownerPassword,
		DuesCount:    duesCount,
	}

	result := db.Create(&apartment)
	assert.NoError(t, result.Error)

	expected, err := repo.GetAllInfoAboutFlat(flatNo)
	assert.NoError(t, err)
	assert.Equal(t, expected, apartment)
}

func TestGetAllInfoAboutAllFlats(t *testing.T) {
	db := setupDb(models.Apartment{})
	repo := NewRepo(db)

	apartments := []models.Apartment{
		{FlatNo: 1, OwnerName: "yusuf", OwnerSurname: "çiftçi", Mail: "yciftci@gmail.com", Password: "123", DuesCount: 0},
		{FlatNo: 2, OwnerName: "ali", OwnerSurname: "veli", Mail: "aliveli@gmail.com", Password: "456", DuesCount: 1},
	}

	for _, apartment := range apartments {
		result := db.Create(&apartment)
		assert.NoError(t, result.Error)
	}

	resultApartments, err := repo.GetAllInfoAboutAllFlats()
	assert.NoError(t, err)

	assert.Len(t, resultApartments, len(apartments))
	for i, expected := range apartments {
		assert.Equal(t, expected.FlatNo, resultApartments[i].FlatNo)
		assert.Equal(t, expected.OwnerName, resultApartments[i].OwnerName)
		assert.Equal(t, expected.OwnerSurname, resultApartments[i].OwnerSurname)
		assert.Equal(t, expected.Mail, resultApartments[i].Mail)
		assert.Equal(t, expected.DuesCount, resultApartments[i].DuesCount)
	}
}

func TestGetDuesCount(t *testing.T) {
	db := setupDb(models.Apartment{})
	repo := NewRepo(db)

	flatNo := 1
	apartment := models.Apartment{
		FlatNo:       flatNo,
		OwnerName:    "Yusuf",
		OwnerSurname: "Çiftçi",
		Mail:         "yciftci@gmail.com",
		Password:     "123",
		DuesCount:    0,
	}

	result := db.Create(&apartment)
	assert.NoError(t, result.Error)

	duesCount, err := repo.GetDuesCount(flatNo)
	assert.NoError(t, err)
	assert.Equal(t, apartment.DuesCount, duesCount)
}

func TestAddDues(t *testing.T) {
	db := setupDb(models.Apartment{})
	repo := NewRepo(db)

	flatNo := 1
	duesCount := 0
	apartment := models.Apartment{
		FlatNo:       flatNo,
		OwnerName:    "Yusuf",
		OwnerSurname: "Çiftçi",
		Mail:         "yciftci@gmail.com",
		Password:     "123",
		DuesCount:    duesCount,
	}

	result := db.Create(&apartment)
	assert.NoError(t, result.Error)

	err := repo.AddDues(flatNo)
	assert.NoError(t, err)

	var actual models.Apartment
	result = db.First(&actual, "flat_no = ?", flatNo)
	assert.NoError(t, result.Error)
	assert.Equal(t, apartment.DuesCount+1, actual.DuesCount)
}

func TestAddDuesForAll(t *testing.T) {
	db := setupDb(models.Apartment{})
	repo := NewRepo(db)

	flatNo1 := 1
	flatNo2 := 2
	duesCount := 0
	apartment1 := models.Apartment{
		FlatNo:       flatNo1,
		OwnerName:    "Yusuf",
		OwnerSurname: "Çiftçi",
		Mail:         "yciftci@gmail.com",
		Password:     "123",
		DuesCount:    duesCount,
	}
	apartment2 := models.Apartment{
		FlatNo:       flatNo2,
		OwnerName:    "Ahmet",
		OwnerSurname: "Demir",
		Mail:         "ademir@gmail.com",
		Password:     "456",
		DuesCount:    duesCount,
	}

	result := db.Create(&apartment1)
	assert.NoError(t, result.Error)
	result = db.Create(&apartment2)
	assert.NoError(t, result.Error)

	err := repo.AddDuesForAll()
	assert.NoError(t, err)

	var actual1, actual2 models.Apartment
	result = db.First(&actual1, "flat_no = ?", flatNo1)
	assert.NoError(t, result.Error)
	result = db.First(&actual2, "flat_no = ?", flatNo2)
	assert.NoError(t, result.Error)

	assert.Equal(t, apartment1.DuesCount+1, actual1.DuesCount)
	assert.Equal(t, apartment2.DuesCount+1, actual2.DuesCount)
}

func TestDeleteDues(t *testing.T) {
	db := setupDb(models.Apartment{})
	repo := NewRepo(db)

	flatNo := 1
	duesCount := 1
	apartment := models.Apartment{
		FlatNo:       flatNo,
		OwnerName:    "Yusuf",
		OwnerSurname: "Çiftçi",
		Mail:         "yciftci@gmail.com",
		Password:     "123",
		DuesCount:    duesCount,
	}

	result := db.Create(&apartment)
	assert.NoError(t, result.Error)

	err := repo.DeleteDues(flatNo)
	assert.NoError(t, err)

	var actual models.Apartment
	result = db.First(&actual, "flat_no = ?", flatNo)
	assert.NoError(t, result.Error)
	assert.Equal(t, apartment.DuesCount-1, actual.DuesCount)
}

func TestDeleteDuesForError(t *testing.T) {
	db := setupDb(models.Apartment{})
	repo := NewRepo(db)

	flatNo := 1
	duesCount := 0
	apartment := models.Apartment{
		FlatNo:       flatNo,
		OwnerName:    "Yusuf",
		OwnerSurname: "Çiftçi",
		Mail:         "yciftci@gmail.com",
		Password:     "123",
		DuesCount:    duesCount,
	}

	result := db.Create(&apartment)
	assert.NoError(t, result.Error)

	err := repo.DeleteDues(flatNo)
	assert.Error(t, err)
	assert.IsType(t, dto.ThereIsNoDues{}, err)
}

func TestGetAllAnnouncements(t *testing.T) {
	db := setupDb(models.Announcement{})
	repo := NewRepo(db)

	announcements := []models.Announcement{
		{Title: "Announcement 1", Content: "Content 1"},
		{Title: "Announcement 2", Content: "Content 2"},
	}

	for _, announcement := range announcements {
		result := db.Create(&announcement)
		assert.NoError(t, result.Error)
	}

	actual, err := repo.GetAllAnnouncements()
	assert.NoError(t, err)

	assert.Len(t, actual, len(announcements))

	for i, announcement := range actual {
		assert.Equal(t, announcements[i].Title, announcement.Title)
		assert.Equal(t, announcements[i].Content, announcement.Content)
	}
}
func TestGetPasswordAndFlatNoByEmail(t *testing.T) {
	// Test veritabanı kurulumunu yap
	db := setupDb(models.Apartment{})
	repo := NewRepo(db)

	// Test verilerini oluştur
	email := "test@example.com"
	expectedPassword := "securepassword"
	expectedFlatNo := 101

	apartment := models.Apartment{
		Mail:     email,
		Password: expectedPassword,
		FlatNo:   expectedFlatNo,
	}

	// Test verisini veritabanına ekle
	result := db.Create(&apartment)
	assert.NoError(t, result.Error)

	// Fonksiyonu test et
	actualPassword, actualFlatNo, err := repo.GetPasswordAndFlatNoByEmail(email)
	assert.NoError(t, err)

	// Beklenen sonuçları doğrula
	assert.Equal(t, expectedPassword, actualPassword)
	assert.Equal(t, expectedFlatNo, actualFlatNo)
}


func TestAddAnnouncement(t *testing.T) {
	db := setupDb(models.Announcement{})
	repo := NewRepo(db)

	announcemement := models.Announcement{
		Title:   "Title1",
		Content: "content1",
	}

	err := repo.AddAnnouncement(announcemement)
	assert.NoError(t, err)

	var actual models.Announcement
	result := db.First(&actual, "title = ?", announcemement.Title)
	assert.NoError(t, result.Error)
	assert.Equal(t, announcemement.Title, actual.Title)
	assert.Equal(t, announcemement.Content, actual.Content)
}

func TestAddMerchant(t *testing.T) {
	// Veritabanı bağlantısını kur
	db := setupDb(models.Merchant{})
	repo := NewRepo(db)

	// Eklenecek satıcının e-posta adresi
	email := "ornek@email.com"

	// Satıcıyı eklemeyi dene
	err := repo.AddMerchant(uuid.New().String(), email)

	// Hata olup olmadığını kontrol et
	assert.NoError(t, err)

	// Satıcının veritabanına başarıyla eklendiğini doğrula
	var count int64
	result := db.Model(&models.Merchant{}).Where("email = ?", email).Count(&count)
	assert.NoError(t, result.Error)
	assert.Equal(t, int64(1), count, "Expected merchant to be added to the database")
}

func TestGetEmailFromMerchant(t *testing.T) {
    // Veritabanı bağlantısını kur
    db := setupDb(models.Merchant{})
    repo := NewRepo(db)

    // Eklenecek satıcının bilgileri
    email := "ornek@email.com"
    newMerchant := models.Merchant{
        MerchantID: uuid.NewString(),
		Email: email,
    }
    result := db.Create(&newMerchant)
    assert.NoError(t, result.Error)

    merchantOID := newMerchant.MerchantID

    retrievedEmail, err := repo.GetEmailFromMerchant(merchantOID)

    assert.NoError(t, err)
    assert.Equal(t, email, retrievedEmail)
}

func TestDeleteDuesByEmail(t *testing.T) {
    db := setupDb(models.Apartment{})
    repo := NewRepo(db)

    email := "ornek@email.com"
    apartment := models.Apartment{
        Mail: email,
        DuesCount: 1,
    }
    result := db.Create(&apartment)
    assert.NoError(t, result.Error)

    err := repo.DeleteDuesByEmail(email)
    
    assert.NoError(t, err)

    var updatedApartment models.Apartment
    db.Where("mail = ?", email).First(&updatedApartment)
    assert.Equal(t, 0, updatedApartment.DuesCount)
}

