package models

type District struct {
	ID           int    `json:"id"`
	SubDistricts string `json:"sub_districts_name"`
	DistrictId   int    `json:"districts_id"`
	District     string `json:"district"`
	ProvinceId   int    `json:"province_id"`
	Province     string `json:"province"`
}
