package models

type Apartment struct {
    FlatNo       int    `gorm:"primaryKey;column:flat_no"`
    OwnerName    string `gorm:"column:owner_name"`
    OwnerSurname string `gorm:"column:owner_surname"`
    Mail         string `gorm:"column:mail;unique"`
    Password     string `gorm:"column:password"`
    DuesCount    int    `gorm:"column:dues_count"`
}

func (Apartment) TableName() string {
    return "apartments"
}

type Announcement struct {
	AnnouncementID int    `gorm:"primaryKey;column:announcement_id;autoIncrement"`
	Title          string `gorm:"column:title;not null"`
	Content        string `gorm:"column:content;not null"`
}

func (Announcement) TableName() string{
	return "announcements"
}

type Merchant struct {
	MerchantID string `gorm:"primaryKey"`
	Email      string `gorm:"column:email;not null"`
}

func (Merchant) TableName() string {
	return "merchants"
}
