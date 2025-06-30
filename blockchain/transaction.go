package blockchain

type Transaction struct {
	PickUp  LatLng `json:"pickUp"`
	DropOff string `json:"dropOff"`
	Price   int    `json:"price"`
}

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
