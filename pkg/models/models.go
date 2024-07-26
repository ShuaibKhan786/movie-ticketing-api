package models


type SignInCredentials struct {
	Role string `json:"role"`
	Provider string `json:"provider"`
	RedirectedURL string `json:"redirected_url"`
}








type Hall struct {
	Name        string `json:"name"`
	Manager     string `json:"manager"`
	Contact     string `json:"contact"`
	Location      Location       `json:"location"`
	SeatLayout    SeatLayout     `json:"seatlayout"`
	OperationTime OperationTime  `json:"operationtime"`
}

type ActualHall struct {
	Name        string 
	Manager     string 
	Contact     string 
	AdminId		int64
}

type Location struct {
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}


type OperationTime struct {
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
}

type SeatLayout struct {
	MaxCapacity int64   `json:"max_capacity"`
	Rows        int    `json:"rows"`
	Columns     int    `json:"columns"`
	Types       string `json:"types"`
	Layout      string `json:"layout"`
}

type UserDetails struct {
	Email string
	EmailVerified bool
	Profile Profile
}

type Profile struct {
	Name string `json:"name"`
	PosterUrl string `json:"profile_pic_url"`
}










type Show struct {
	Movie           Movie           `json:"movie"`
	Cast            Cast            `json:"cast"`
	MovieShowTiming []ShowDate `json:"movie_show_timing"`
}

type Movie struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int64  `json:"duration"`
	Genre       string `json:"genre"`
	ReleaseDate string `json:"release_date"`
}

type Cast struct {
	Actors   []CastBlueprint `json:"actors"`
	Actress  []CastBlueprint `json:"actress"`
	Directors []CastBlueprint `json:"directors"`
	Producers []CastBlueprint `json:"producers"`
}

type CastBlueprint struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Alias string `json:"alias"`
}

type ShowDate struct {
	Date  string `json:"show_date"`
	Timing []string `json:"show_timing"`
}
