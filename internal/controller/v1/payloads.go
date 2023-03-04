package v1

type GetBeersResponse struct {
	Body []BeersPayload
}

type BeersPayload struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// --------------------------------------------

type SaveBeerPayload struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// --------------------------------------------

type GetUserBeerPreferencesResponse struct {
	Body []UserBeerPreferencesPayload
}

type UserBeerPreferencesPayload struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	DrunkTheBeerBefore bool   `json:"drunk_before"`
	GotDrunk           bool   `json:"got_drunk"`
	LastTime           string `json:"last_time"`
	Rating             int    `json:"rating"`
	Comment            string `json:"comment"`
}

// --------------------------------------------
