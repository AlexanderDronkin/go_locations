package main

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/utf8string"
)

// go make cities
func loadCities() {
	var cityFull []interface{}
	var cityShort []interface{}

	// go steal some data from cdek
	if cdekApi.Auth().token == "" {
		panic("Не могу авторизоваться на сервере CDEK API")
	}

	cdekRegions := cdekApi.regionList()
	zoneByCode := map[int]string{}

	cursor, _ := mongodb.collection("tariffZones").Find(mongodb.context, bson.D{})
	for cursor.TryNext(mongodb.context) {
		var zone Zone
		cursor.Decode(&zone)
		for _, row := range cdekRegions {
			if zone.FIAS == row.FIAS {
				zoneByCode[row.Code] = zone.Code
				break
			}
		}
	}

	page := 0
	for {
		cityList := cdekApi.cityList(10000, page)
		if len(cityList) == 0 {
			break
		}
		for _, city := range cityList {
			zone, found := zoneByCode[city.RegionCode]
			if len(city.FIAS) != 36 || !found {
				continue
			}
			short, full := city.adapt()
			full.Zone = zone
			cityFull = append(cityFull, full)
			cityShort = append(cityShort, short)
		}
		page++
		time.Sleep(time.Second / 10)
	}

	mongodb.collection("cityFull", true).InsertMany(mongodb.context, cityFull)
	mongodb.collection("cityShort", true).InsertMany(mongodb.context, cityShort)
}

// go write cities in files
func writeCities() {
	os.MkdirAll(DIR+"city/fias", 0766)

	cursor, _ := mongodb.collection("cityFull").Find(mongodb.context, bson.D{})
	for cursor.TryNext(mongodb.context) {
		var city CityFull
		cursor.Decode(&city)
		if len(city.FIAS) == 36 {
			fput("city/fias/"+city.FIAS, city.toJson())
		}
	}

	var cities []CityShort
	abc := map[string][]CityShort{}
	cursor, _ = mongodb.collection("cityShort").Find(mongodb.context, bson.D{})
	for cursor.TryNext(mongodb.context) {
		var city CityShort
		cursor.Decode(&city)
		if len(city.FIAS) == 36 {
			name := utf8string.NewString(strings.ToLower(strings.TrimSpace(city.Name))).Slice(0, 2)
			cities = append(cities, city)
			abc[name] = append(abc[name], city)
		}
	}
	citiesJson, _ := json.Marshal(cities)
	fput("city/index.json", citiesJson)
	fput_br("city/index.json", citiesJson)

	for name, value := range abc {
		abcJson, _ := json.Marshal(value)
		fput("city/"+name, abcJson)
		fput_br("city/"+name, abcJson)
	}

	fgc("city", time.Hour)
}

// go get tariff zones from bitrix hb
func loadZones() {
	const (
		sqlFix = `#_PRIVATE_#`
		sql    = `#_PRIVATE_#`
	)

	var zones []interface{}

	// charset fix
	mysql.use("taber").Exec(sqlFix)
	result, err := mysql.use("taber").Query(sql)
	if err != nil {
		panic(err.Error())
	}
	for result.Next() {
		row := ZoneMysql{}
		err = result.Scan(&row.Code, &row.FIAS, &row.KLADR, &row.Name, &row.CourierPrice, &row.CourierFreeFrom, &row.PvzPrice, &row.PvzFreeFrom)
		if err != nil {
			panic(err.Error())
		}
		zones = append(zones, row.adapt())
	}

	mongodb.collection("tariffZones", true).InsertMany(mongodb.context, zones)
}

// go get pvz list
func loadPvz() {
	const (
		dropView   = `#_PRIVATE_#`
		createView = `#_PRIVATE_#`
		pvz        = `#_PRIVATE_#`
	)
	var pvzList []interface{}

	db := mysql.use("taber")
	db.Exec(dropView)
	db.Exec(createView)
	result, err := db.Query(pvz)
	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		row := PvzMysql{}
		err = result.Scan(
			&row.CityFIAS,
			&row.GfCity,
			&row.Code,
			&row.DeliveryID,
			&row.DeliveryCode,
			&row.DeliveryClassName,
			&row.DeliveryDaysMin,
			&row.DeliveryDaysMax,
			&row.Name,
			&row.Address,
			&row.Tel,
			&row.Schedule,
			&row.Wayto,
			&row.Latitude,
			&row.Longitude,
			&row.Options,
		)
		if err != nil {
			panic(err.Error())
		}
		pvz := row.adapt()
		if len(pvz.Options.ZipCode) > 0 {
			pvz.CityPostal = pvz.Options.ZipCode
		}
		if !pvz.Options.CardPaymentAvailable {
			pvz.PaymentCard = false
		}
		if pvz.Options.IsCashForbidden {
			pvz.PaymentCash = false
		}
		// go magik
		if pvz.DeliveryCode == "#_PRIVATE_#" {
			pvz.DeliveryType = "#_PRIVATE_#"
		}
		if pvz.DeliveryCode == "#_PRIVATE_#" ||
			pvz.DeliveryCode == "#_PRIVATE_#" ||
			pvz.DeliveryCode == "#_PRIVATE_#" {
			pvz.DeliveryType = "#_PRIVATE_#"
		}
		if pvz.DeliveryID == 54 {
			pvz.DeliveryType = "#_PRIVATE_#"
		}
		pvzList = append(pvzList, pvz)
	}

	mongodb.collection("pvz", true).InsertMany(mongodb.context, pvzList)
	mongodb.index("pvz", "#_PRIVATE_#")
	mongodb.index("pvz", "#_PRIVATE_#")
	mongodb.index("pvz", "#_PRIVATE_#")
	db.Exec(dropView)
}

func writePvz() {
	os.MkdirAll(DIR+"shops/pvz", 0766)
	os.MkdirAll(DIR+"shops/psp", 0766)
	os.MkdirAll(DIR+"shops/letual", 0766)
	os.MkdirAll(DIR+"shops/pickup", 0766)

	pvzAll := []Pvz{}
	pvzCity := map[string][]Pvz{}

	cursor, _ := mongodb.collection("pvz").Find(mongodb.context, bson.D{})
	for cursor.TryNext(mongodb.context) {
		var pvz Pvz
		cursor.Decode(&pvz)

		if len(pvz.CityFIAS) == 36 {
			pvzAll = append(pvzAll, pvz)
			pvzCity[pvz.CityFIAS] = append(pvzCity[pvz.CityFIAS], pvz)
		}
	}

	pvzAllJson, _ := json.Marshal(pvzAll)
	fput("shops/index.json", pvzAllJson)
	fput_br("shops/index.json", pvzAllJson)

	for fias, cityPvzList := range pvzCity {
		pvzTyped := map[string][]Pvz{}
		cityPvzListJson, _ := json.Marshal(cityPvzList)

		fput("shops/"+fias, cityPvzListJson)
		fput_br("shops/"+fias, cityPvzListJson)

		for _, cityPvz := range cityPvzList {
			pvzTyped[cityPvz.DeliveryType] = append(pvzTyped[cityPvz.DeliveryType], cityPvz)
		}
		for typed, pvzTypedList := range pvzTyped {
			cityPvzTypedJson, _ := json.Marshal(pvzTypedList)
			fput("shops/"+typed+"/"+fias, cityPvzTypedJson)
			fput_br("shops/"+typed+"/"+fias, cityPvzTypedJson)
		}
	}

	fgc("shops", time.Hour)
}
