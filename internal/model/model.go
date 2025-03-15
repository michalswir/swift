package model

type SwiftData struct {
	CountryISO2   string `bson:"countryIso2" json:"countryIso2"`
	SwiftCode     string `bson:"swiftCode" json:"swiftCode"`
	CodeType      string `bson:"codeType" json:"codeType"`
	BankName      string `bson:"bankName" json:"bankName"`
	Address       string `bson:"address" json:"address"`
	TownName      string `bson:"townName" json:"townName"`
	CountryName   string `bson:"countryName" json:"countryName"`
	TimeZone      string `bson:"timeZone" json:"timeZone"`
	IsHeadquarter bool   `bson:"isHeadquarter" json:"isHeadquarter"`
}
