package beerapi

type beer struct {
	ID               int     `json:"id,omitempty"`
	Name             string  `json:"name,omitempty"`
	Tagline          string  `json:"tagline,omitempty"`
	FirstBrewed      string  `json:"first_brewed,omitempty"`
	Description      string  `json:"description,omitempty"`
	ImageURL         string  `json:"image_url,omitempty"`
	Abv              float64 `json:"abv,omitempty"`
	Ibu              float64 `json:"ibu,omitempty"`
	TargetFg         float64 `json:"target_fg,omitempty"`
	TargetOg         float64 `json:"target_og,omitempty"`
	Ebc              float64 `json:"ebc,omitempty"`
	Srm              float64 `json:"srm,omitempty"`
	Ph               float64 `json:"ph,omitempty"`
	AttenuationLevel float64 `json:"attenuation_level,omitempty"`
	Volume           struct {
		Value int    `json:"value,omitempty"`
		Unit  string `json:"unit,omitempty"`
	} `json:"volume,omitempty"`
	BoilVolume struct {
		Value int    `json:"value,omitempty"`
		Unit  string `json:"unit,omitempty"`
	} `json:"boil_volume,omitempty"`
	Method struct {
		MashTemp []struct {
			Temp struct {
				Value int    `json:"value,omitempty"`
				Unit  string `json:"unit,omitempty"`
			} `json:"temp,omitempty"`
			Duration int `json:"duration,omitempty"`
		} `json:"mash_temp,omitempty"`
		Fermentation struct {
			Temp struct {
				Value float64 `json:"value,omitempty"`
				Unit  string  `json:"unit,omitempty"`
			} `json:"temp,omitempty"`
		} `json:"fermentation,omitempty"`
		Twist interface{} `json:"twist,omitempty"`
	} `json:"method,omitempty"`
	Ingredients struct {
		Malt []struct {
			Name   string `json:"name,omitempty"`
			Amount struct {
				Value float64 `json:"value,omitempty"`
				Unit  string  `json:"unit,omitempty"`
			} `json:"amount,omitempty"`
		} `json:"malt,omitempty"`
		Hops []struct {
			Name   string `json:"name,omitempty"`
			Amount struct {
				Value float64 `json:"value,omitempty"`
				Unit  string  `json:"unit,omitempty"`
			} `json:"amount,omitempty"`
			Add       string `json:"add,omitempty"`
			Attribute string `json:"attribute,omitempty"`
		} `json:"hops,omitempty"`
		Yeast string `json:"yeast,omitempty"`
	} `json:"ingredients,omitempty"`
	FoodPairing   []string `json:"food_pairing,omitempty"`
	BrewersTips   string   `json:"brewers_tips,omitempty"`
	ContributedBy string   `json:"contributed_by,omitempty"`
}
