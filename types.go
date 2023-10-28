package main

import (
	"database/sql"
	"encoding/json"
	"strings"
)

// delivery tariffs zone
type Zone struct {
	Code            string  `json:"code" bson:"_id"`
	FIAS            string  `json:"FIAS" bson:"FIAS"`
	KLADR           string  `json:"KLADR,omitempty" bson:"KLADR"`
	Name            string  `json:"Name" bson:"Name"`
	CourierPrice    float64 `json:"CourierPrice" bson:"CourierPrice"`
	CourierFreeFrom float64 `json:"CourierFreeFrom" bson:"CourierFreeFrom"`
	PvzPrice        float64 `json:"PvzPrice" bson:"PvzPrice"`
	PvzFreeFrom     float64 `json:"PvzFreeFrom" bson:"PvzFreeFrom"`
}

func (zone Zone) toJson() []byte { return toJson(zone) }

type ZoneMysql struct {
	Code            sql.RawBytes
	FIAS            sql.RawBytes
	KLADR           sql.RawBytes
	Name            sql.RawBytes
	CourierPrice    sql.NullInt64
	CourierFreeFrom sql.NullInt64
	PvzPrice        sql.NullInt64
	PvzFreeFrom     sql.NullInt64
}

func (row ZoneMysql) adapt() Zone {
	return Zone{
		Code:            string(row.Code),
		FIAS:            string(row.FIAS),
		KLADR:           string(row.KLADR),
		Name:            string(row.Name),
		CourierPrice:    float64(row.CourierPrice.Int64),
		CourierFreeFrom: float64(row.CourierFreeFrom.Int64),
		PvzPrice:        float64(row.PvzPrice.Int64),
		PvzFreeFrom:     float64(row.PvzFreeFrom.Int64),
	}
}

// pvz import
type Pvz struct {
	Code            string  `json:"Code" bson:"Code"`
	CityFIAS        string  `json:"CityFIAS" bson:"CityFIAS"`
	CityPostal      string  `json:"CityPostal,omitempty" bson:"CityPostal"`
	GfCity          string  `json:"GfCity" bson:"GfCity"`
	DeliveryID      int     `json:"DeliveryID" bson:"DeliveryID"`
	DeliveryCode    string  `json:"DeliveryCode,omitempty" bson:"DeliveryCode"`
	DeliveryType    string  `json:"DeliveryType" bson:"DeliveryType"`
	DeliveryDaysMin int     `json:"DeliveryDaysMin" bson:"DeliveryDaysMin"`
	DeliveryDaysMax int     `json:"DeliveryDaysMax" bson:"DeliveryDaysMax"`
	Name            string  `json:"Name" bson:"Name"`
	Address         string  `json:"Address" bson:"Address"`
	Tel             string  `json:"Tel" bson:"Tel"`
	Schedule        string  `json:"Schedule" bson:"Schedule"`
	Wayto           string  `json:"Wayto" bson:"Wayto"`
	Latitude        float64 `json:"Latitude,omitempty" bson:"Latitude"`
	Longitude       float64 `json:"Longitude,omitempty" bson:"Longitude"`
	PaymentCard     bool    `json:"PaymentCard" bson:"PaymentCard"`
	PaymentCash     bool    `json:"PaymentCash" bson:"PaymentCash"`
	Options         PvzOptions
}

func (pvz Pvz) toJson() []byte { return toJson(pvz) }

type PvzOptions struct {
	IsCashForbidden      bool                   `json:"IsCashForbidden"`
	CardPaymentAvailable bool                   `json:"CardPaymentAvailable"`
	ZipCode              string                 `json:"ZipCode"`
	Address              map[string]interface{} `json:"Address"`
}

type PvzMysql struct {
	CityFIAS          sql.RawBytes
	GfCity            sql.RawBytes
	Code              sql.RawBytes
	DeliveryID        sql.NullInt64
	DeliveryCode      sql.RawBytes
	DeliveryClassName sql.RawBytes
	DeliveryDaysMin   sql.NullInt64
	DeliveryDaysMax   sql.NullInt64
	Name              sql.RawBytes
	Address           sql.RawBytes
	Tel               sql.RawBytes
	Schedule          sql.RawBytes
	Wayto             sql.RawBytes
	Latitude          sql.NullFloat64
	Longitude         sql.NullFloat64
	Options           sql.RawBytes
}

func (row PvzMysql) adapt() Pvz {
	pvzOptions := PvzOptions{}
	json.Unmarshal([]byte(row.Options), &pvzOptions)

	return Pvz{
		Code:            string(row.Code),
		CityFIAS:        string(row.CityFIAS),
		GfCity:          string(row.GfCity),
		DeliveryID:      int(row.DeliveryID.Int64),
		DeliveryCode:    string(row.DeliveryCode),
		DeliveryType:    "pvz",
		DeliveryDaysMin: int(row.DeliveryDaysMin.Int64),
		DeliveryDaysMax: int(row.DeliveryDaysMax.Int64),
		Name:            string(row.Name),
		Address:         string(row.Address),
		Tel:             string(row.Tel),
		Schedule:        string(row.Schedule),
		Wayto:           string(row.Wayto),
		Latitude:        float64(row.Latitude.Float64),
		Longitude:       float64(row.Longitude.Float64),
		PaymentCard:     true,
		PaymentCash:     true,
		Options:         pvzOptions,
	}
}

// new city data
type CityShort struct {
	FIAS    string `json:"FIAS" bson:"_id"`
	Name    string `json:"Name" bson:"Name"`
	Region  string `json:"Region" bson:"Region"`
	Courier bool   `json:"Courier,omitempty" bson:"Courier,omitempty"`
	Pvz     bool   `json:"Pvz,omitempty" bson:"Pvz,omitempty"`
	Psp     bool   `json:"Psp,omitempty" bson:"Psp,omitempty"`
	Pickup  bool   `json:"Pickup,omitempty" bson:"Pickup,omitempty"`
}

func (city CityShort) toJson() []byte { return toJson(city) }

type CityFull struct {
	FIAS      string  `json:"FIAS" bson:"_id"`
	KLADR     string  `json:"KLADR,omitempty" bson:"KLADR"`
	Zone      string  `json:"Zone" bson:"Zone"`
	Name      string  `json:"Name" bson:"Name"`
	Region    string  `json:"Region" bson:"Region"`
	SubRegion string  `json:"SubRegion,omitempty" bson:"SubRegion"`
	Longitude float64 `json:"Longitude,omitempty" bson:"Longitude"`
	Latitude  float64 `json:"Latitude,omitempty" bson:"Latitude"`
}

func (city CityFull) toJson() []byte { return toJson(city) }

// from cdek api
type cdekAuth struct {
	Token   string `json:"access_token"`
	Ttype   string `json:"token_type"`
	Expires int    `json:"expires_in"`
	Scopes  string `json:"scope"`
	Id      string `json:"jti"`
}

type cdekRegion struct {
	Code  int    `json:"region_code"`
	Name  string `json:"region"`
	FIAS  string `json:"fias_region_guid"`
	KLADR string `json:"kladr_region_code"`
}

type cdekCity struct {
	Code         int     `json:"code"`
	Name         string  `json:"city"`
	FIAS         string  `json:"fias_guid"`
	KLADR        string  `json:"kladr_code"`
	RegionCode   int     `json:"region_code"`
	Region       string  `json:"region"`
	SubRegion    string  `json:"sub_region"`
	Longitude    float64 `json:"longitude"`
	Latitude     float64 `json:"latitude"`
	PaymentLimit float64 `json:"payment_limit"`
}

func (city cdekCity) adapt() (CityShort, CityFull) {
	regionName := city.Region
	if len(city.SubRegion) > 0 {
		regionName = strings.Join([]string{city.Region, city.SubRegion}, ", ")
	}

	return CityShort{
			FIAS:   city.FIAS,
			Name:   city.Name,
			Region: regionName,
		},
		CityFull{
			FIAS:      city.FIAS,
			Name:      city.Name,
			Region:    city.Region,
			SubRegion: city.SubRegion,
			Longitude: float64(city.Longitude),
			Latitude:  float64(city.Latitude),
		}
}
