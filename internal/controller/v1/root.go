package v1

import (
	"encoding/json"
	"errors"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/viktorcitaku/beer-app/internal/beerapi"
	"github.com/viktorcitaku/beer-app/internal/cache"
	"github.com/viktorcitaku/beer-app/internal/repository"
)

type Controller struct {
	beerApiClient                 *beerapi.Client
	userProfileCrudRepository     repository.UserProfileCrudRepository
	userPreferencesCrudRepository repository.UserPreferencesCrudRepository
	cacheProvider                 cache.Provider
}

func New(
	beerApiClient *beerapi.Client,
	userProfileCrudRepository repository.UserProfileCrudRepository,
	userPreferencesCrudRepository repository.UserPreferencesCrudRepository,
	cacheProvider cache.Provider,
) *Controller {
	return &Controller{
		beerApiClient:                 beerApiClient,
		userProfileCrudRepository:     userProfileCrudRepository,
		userPreferencesCrudRepository: userPreferencesCrudRepository,
		cacheProvider:                 cacheProvider,
	}
}

func (c *Controller) HelloWorld(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Hello World"))
}

func (c *Controller) SaveUserProfiles(w http.ResponseWriter, r *http.Request) {
	var err error
	err = r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	email := r.Form["email"][0]
	err = c.cacheProvider.GetRedisCache().Set(r.Context(), email, time.Now(), 0).Err()
	if err != nil {
		log.Printf("Error setting session: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = c.userProfileCrudRepository.Save(&repository.UserProfile{
		Email:      email,
		LastUpdate: time.Now(),
	})
	if err != nil {
		log.Printf("Error saving user profile: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Initialize a new cookie containing the string "Hello world!" and some
	// non-default attributes.
	cookie := http.Cookie{
		Name:  "BEER_SESSION",
		Value: email,
		Path:  "/",
	}

	// Use the http.SetCookie() function to send the cookie to the client.
	// Behind the scenes, this adds a `Set-Cookie` header to the response
	// containing the necessary cookie data.
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) GetBeers(w http.ResponseWriter, r *http.Request) {
	err := c.sessionValidation(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pairingFood, bitterness, fermentation := extractQueryParams(r)
	beers, err := c.beerApiClient.GetBeers(pairingFood, bitterness, fermentation)
	if err != nil {
		log.Printf("Error getting beers: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var beerPayloads []BeersPayload
	for _, beer := range beers {
		beerPayloads = append(beerPayloads, BeersPayload{
			ID:   beer.ID,
			Name: beer.Name,
		})
	}

	err = json.NewEncoder(w).Encode(beerPayloads)
	if err != nil {
		log.Printf("Error encoding beers: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func extractQueryParams(r *http.Request) (string, string, string) {
	var finalPairingFood string
	pairingFood := r.URL.Query()["pairing_food"]
	if len(pairingFood) > 0 {
		finalPairingFood = pairingFood[0]
	}

	var finalBitterness string
	bitterness := r.URL.Query()["bitterness"]
	if len(bitterness) > 0 {
		finalBitterness = bitterness[0]
	}

	var finalFermentation string
	fermentation := r.URL.Query()["fermentation"]
	if len(fermentation) > 0 {
		finalFermentation = fermentation[0]
	}

	return finalPairingFood, finalBitterness, finalFermentation
}

func (c *Controller) SaveBeers(w http.ResponseWriter, r *http.Request) {
	err := c.sessionValidation(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie, _ := r.Cookie("BEER_SESSION")
	_, err = c.cacheProvider.GetRedisCache().Get(r.Context(), cookie.Value).Result()
	if err != nil {
		log.Printf("Error getting session: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var saveBeerPayload SaveBeerPayload
	err = json.NewDecoder(r.Body).Decode(&saveBeerPayload)
	if err != nil {
		log.Printf("Error decoding request: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userPreferences, err := c.userPreferencesCrudRepository.FindByBeerId(uint(saveBeerPayload.ID))
	if err != nil {
		log.Printf("Error getting user preferences: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userPreferences != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	userPreferences = &repository.UserPreferences{
		Id:        GenerateUniqueID(saveBeerPayload.ID, cookie.Value),
		BeerId:    uint(saveBeerPayload.ID),
		BeerName:  saveBeerPayload.Name,
		UserEmail: cookie.Value,
	}

	err = c.userPreferencesCrudRepository.Save(userPreferences)
	if err != nil {
		log.Printf("Error saving user preferences: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GenerateUniqueID(beerID int, userEmail string) uint {
	h := fnv.New32a()
	_, _ = io.WriteString(h, userEmail)
	uniqueID := int(h.Sum32()) + beerID
	return uint(uniqueID)
}

func (c *Controller) GetUserPreferences(w http.ResponseWriter, r *http.Request) {
	err := c.sessionValidation(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie, _ := r.Cookie("BEER_SESSION")
	_, err = c.cacheProvider.GetRedisCache().Get(r.Context(), cookie.Value).Result()
	if err != nil {
		log.Printf("Error getting session: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userPreferences, err := c.userPreferencesCrudRepository.FindByUserEmail(cookie.Value)
	if err != nil {
		log.Printf("Error getting user preferences: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userBeerPreferencesPayload []UserBeerPreferencesPayload
	for _, userPreference := range userPreferences {
		userBeerPreferencesPayload = append(userBeerPreferencesPayload, UserBeerPreferencesPayload{
			ID:                 int(userPreference.Id),
			Name:               userPreference.BeerName,
			DrunkTheBeerBefore: userPreference.DrunkTheBeerBefore,
			GotDrunk:           userPreference.GotDrunk,
			LastTime:           userPreference.LastTime,
			Rating:             int(userPreference.Rating),
			Comment:            userPreference.Comment,
		})
	}

	err = json.NewEncoder(w).Encode(userBeerPreferencesPayload)
	if err != nil {
		log.Printf("Error encoding user preferences: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) SaveUserPreferences(w http.ResponseWriter, r *http.Request) {
	var err error
	err = c.sessionValidation(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie, _ := r.Cookie("BEER_SESSION")
	_, err = c.cacheProvider.GetRedisCache().Get(r.Context(), cookie.Value).Result()
	if err != nil {
		log.Printf("Error getting session: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userBeerPreferencesPayloads []UserBeerPreferencesPayload
	err = json.Unmarshal(bytes, &userBeerPreferencesPayloads)
	if err != nil {
		log.Printf("Error decoding request: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, userBeerPreferencesPayload := range userBeerPreferencesPayloads {
		var userPreferences *repository.UserPreferences
		userPreferences, err = c.userPreferencesCrudRepository.FindById(uint(userBeerPreferencesPayload.ID))
		if err != nil {
			log.Printf("Error getting user preferences: %v\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userPreferences.BeerName = userBeerPreferencesPayload.Name
		userPreferences.DrunkTheBeerBefore = userBeerPreferencesPayload.DrunkTheBeerBefore
		userPreferences.GotDrunk = userBeerPreferencesPayload.GotDrunk
		userPreferences.LastTime = userBeerPreferencesPayload.LastTime
		userPreferences.Rating = uint(userBeerPreferencesPayload.Rating)
		userPreferences.Comment = userBeerPreferencesPayload.Comment

		err = c.userPreferencesCrudRepository.Save(userPreferences)
		if err != nil {
			log.Printf("Error saving user preferences: %v\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) sessionValidation(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("BEER_SESSION")
	if err != nil {
		log.Printf("Error getting cookie: %v\n", err.Error())
		return errors.New("error getting cookie")
	}

	_, err = c.cacheProvider.GetRedisCache().Get(r.Context(), cookie.Value).Result()
	if err != nil {
		log.Printf("Error session mismatch: %v\n", err.Error())
		return errors.New("error session mismatch")
	}

	http.SetCookie(w, cookie)

	return nil
}
