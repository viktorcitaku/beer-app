package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserProfileRepository struct {
	ctx context.Context
	db  *sql.DB
}

func NewUserProfileRepository(
	ctx context.Context,
	db *sql.DB,
) *UserProfileRepository {
	return &UserProfileRepository{
		ctx: ctx,
		db:  db,
	}
}

func (u *UserProfileRepository) Save(user *UserProfile) error {
	var err error
	tx, err := u.db.BeginTx(u.ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT * FROM beer.user_profile WHERE email = '%s'", user.Email)
	row := tx.QueryRowContext(u.ctx, query)
	if err = row.Err(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("Error querying user profile: %v\n", err.Error())
		return err
	}

	insert := "INSERT INTO beer.user_profile(email, last_update) VALUES ($1, $2)"
	_, err = tx.ExecContext(u.ctx, insert, user.Email, user.LastUpdate)
	if err != nil {
		log.Printf("Error inserting user profile: %v\n", err.Error())
		if err = tx.Rollback(); err != nil {
			log.Printf("Error rolling back transaction: %v\n", err.Error())
			return err
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v\n", err.Error())
		return err
	}

	return nil
}

func (u *UserProfileRepository) FindByEmail(email string) (*UserProfile, error) {
	query := fmt.Sprintf("SELECT * FROM beer.user_profile WHERE email = '%s'", email)
	row := u.db.QueryRow(query)
	if err := row.Err(); err != nil {
		return nil, err
	}

	userProfile := &UserProfile{}
	if err := row.Scan(&userProfile.Email, &userProfile.LastUpdate); err != nil {
		return nil, err
	}

	return userProfile, nil
}

type UserPreferencesRepository struct {
	ctx    context.Context
	client *mongo.Client
}

func NewUserPreferencesRepository(
	ctx context.Context,
	client *mongo.Client,
) *UserPreferencesRepository {
	return &UserPreferencesRepository{
		ctx:    ctx,
		client: client,
	}
}

func (u *UserPreferencesRepository) Save(preferences *UserPreferences) error {
	coll := u.client.Database("beer_db").Collection("beer")
	if coll == nil {
		return errors.New("empty collection")
	}

	updateOptions := options.Replace().SetUpsert(true)
	filter := bson.D{{"_id", preferences.Id}}
	_, err := coll.ReplaceOne(u.ctx, filter, preferences, updateOptions)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserPreferencesRepository) FindById(id uint) (*UserPreferences, error) {
	coll := u.client.Database("beer_db").Collection("beer")
	if coll == nil {
		return nil, errors.New("empty collection")
	}

	var err error
	filter := bson.D{{"_id", id}}
	cur, err := coll.Find(u.ctx, filter)
	if err != nil {
		return nil, err
	}

	var ups []*UserPreferences
	if err = cur.All(u.ctx, &ups); err != nil {
		return nil, err
	}

	if size := len(ups); size > 1 || size == 0 {
		return nil, errors.New("zero or more than one element found")
	}

	return ups[0], nil
}

func (u *UserPreferencesRepository) FindByBeerId(beerId uint) (*UserPreferences, error) {
	coll := u.client.Database("beer_db").Collection("beer")
	if coll == nil {
		return nil, errors.New("empty collection")
	}

	var err error
	filter := bson.D{{"beer_id", beerId}}
	cur, err := coll.Find(u.ctx, filter)
	if err != nil {
		return nil, err
	}

	var ups []*UserPreferences
	if err = cur.All(u.ctx, &ups); err != nil {
		return nil, err
	}

	size := len(ups)
	if size > 1 {
		return nil, errors.New("more than one element found")
	}

	if size == 0 {
		return nil, nil
	}

	return ups[0], nil
}

func (u *UserPreferencesRepository) FindByUserEmail(email string) ([]*UserPreferences, error) {
	coll := u.client.Database("beer_db").Collection("beer")
	if coll == nil {
		return nil, errors.New("empty collection")
	}

	var err error
	filter := bson.D{{"user_email", email}}
	cur, err := coll.Find(u.ctx, filter)
	if err != nil {
		return nil, err
	}

	var ups []*UserPreferences
	if err = cur.All(u.ctx, &ups); err != nil {
		return nil, err
	}

	return ups, nil
}
