package models

type SignInCredentials struct {
	Role          string `json:"role"`
	Provider      string `json:"provider"`
	RedirectedURL string `json:"redirected_url"`
}

type Hall struct {
	Name          string        `json:"name"`
	Manager       string        `json:"manager"`
	Contact       string        `json:"contact"`
	Location      Location      `json:"location"`
	OperationTime OperationTime `json:"operationtime"`
}

type ActualHall struct {
	Name    string
	Manager string
	Contact string
	AdminId int64
}

type Location struct {
	Address    string  `json:"address"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	PostalCode string  `json:"postal_code"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

type OperationTime struct {
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
}

type SeatLayout struct {
	SeatTypes []SeatType `json:"seat_types"`
}

type SeatType struct {
	Name            *string   `json:"name"`
	Price           *int      `json:"price"`
	SeatRow         *int      `json:"seat_row"`
	SeatColumn      *int      `json:"seat_column"`
	SeatMatrix      *string   `json:"seat_matrix"`
	OrderFromScreen *int      `json:"order_from_screen"`
	RowName         []string `json:"row_names"`
}

type SeatRowNameUpdate struct {
	RowName *string `json:"row_name"`
}

type UserDetails struct {
	Email         string
	EmailVerified bool
	Profile       Profile
}

type Profile struct {
	Name      string `json:"name"`
	PosterUrl string `json:"profile_pic_url"`
}

type Show struct {
	Status          string     `json:"status"`
	Movie           Movie      `json:"movie"`
	Cast            Cast       `json:"cast"`
	MovieShowTiming []ShowDate `json:"movie_show_timing"`
}

type Movie struct {
	Id           int64  `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Duration     int64  `json:"duration"`
	Genre        string `json:"genre"`
	ReleaseDate  string `json:"release_date"`
	PortraitUrl  string `json:"portrait_url"`
	LandscapeUrl string `json:"landscape_url"`
}

type Cast struct {
	Actors    []CastBlueprint `json:"actors"`
	Actress   []CastBlueprint `json:"actress"`
	Directors []CastBlueprint `json:"directors"`
	Producers []CastBlueprint `json:"producers"`
}

type CastBlueprint struct {
	Id        int64   `json:"id"`
	Name      string  `json:"name"`
	Alias     string  `json:"alias"`
	PosterUrl *string `json:"poster"`
	// PosterUrl is a pointer because cast posters are optional;
	// using a pointer allows handling NULL values returned from the database.
}

type ShowDate struct {
	Date   string   `json:"show_date"`
	Timing []Timing `json:"show_timing"`
}
type Timing struct {
	Time         string `json:"time"`
	TicketStatus bool   `json:"ticket_status"`
	//if TicketStatus field is set to true
	//it means ticket is avilable for booking of that timing
	PreExpiry  int `json:"pre_expiry_secs"`
	PostExpiry int `json:"post_expiry_secs"`
}
