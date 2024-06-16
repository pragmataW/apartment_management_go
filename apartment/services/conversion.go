package services

import "github.com/pragmataW/apartment_management/models"

type Apartment struct {
	FlatNo       int
	OwnerName    string
	OwnerSurname string
	Mail         string
	Password     string
	DuesCount    int
}

func (ar *Apartment) ToApartmentModel() models.Apartment {
	apartment := models.Apartment{
		FlatNo:       ar.FlatNo,
		OwnerName:    ar.OwnerName,
		OwnerSurname: ar.OwnerSurname,
		Mail:         ar.Mail,
		Password:     ar.Password,
		DuesCount:    ar.DuesCount,
	}
	return apartment
}

func (ar *Apartment) ToApartmentServiceObject(apartment models.Apartment){
	ar.FlatNo = apartment.FlatNo
	ar.OwnerName = apartment.OwnerName
	ar.OwnerSurname = apartment.OwnerSurname
	ar.Mail = apartment.Mail
	ar.Password = apartment.Password
	ar.DuesCount = apartment.DuesCount
}

//announcement

type Announcement struct {
	AnnouncementID int
	Title          string
	Content        string
}

func (an *Announcement) ToAnnouncementServiceObject(announcement models.Announcement) {
	an.AnnouncementID = announcement.AnnouncementID
	an.Title = announcement.Title
	an.Content = announcement.Content
}

func (an Announcement) ToAnnouncementModel() models.Announcement{
	return models.Announcement{
		AnnouncementID: an.AnnouncementID,
		Title: an.Title,
		Content: an.Content,
	}
}