package model

type Inventory struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	StockKey string `json:"stockKey" bson:"stockKey"`
	Count    int64  `json:"count" bson:"count"`
}
