package repository

import "time"

type UserProfile struct {
	Email      string
	LastUpdate time.Time
}

type UserProfileCrudRepository interface {
	Save(user *UserProfile) error
	FindByEmail(email string) (*UserProfile, error)
}

type UserPreferences struct {
	Id                 uint   `bson:"_id,minsize"`
	BeerId             uint   `bson:"beer_id,minsize"`
	BeerName           string `bson:"beer_name"`
	UserEmail          string `bson:"user_email"`
	DrunkTheBeerBefore bool   `bson:"drunk_the_beer_before"`
	GotDrunk           bool   `bson:"got_drunk"`
	LastTime           string `bson:"last_time"`
	Rating             uint   `bson:"rating,minsize"`
	Comment            string `bson:"comment"`
}

type UserPreferencesCrudRepository interface {
	Save(preferences *UserPreferences) error
	FindById(id uint) (*UserPreferences, error)
	FindByBeerId(id uint) (*UserPreferences, error)
	FindByUserEmail(email string) ([]*UserPreferences, error)
}
