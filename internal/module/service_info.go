package models

type ServiceInfo struct {
	tableName      struct{} `pg:"service_info"`
	Id             int64
	Service        string
	Server         string
	Contact_type   string
	Contact_first  string
	Contact_second string
	Contact_third  string
}
