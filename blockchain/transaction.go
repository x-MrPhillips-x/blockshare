package blockchain

type Transaction struct {
	PickUp  LatLng `json:"pickUp"`
	DropOff string `json:"dropOff"`
	Price   int    `json:"price"`
}

type LatLng struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}
