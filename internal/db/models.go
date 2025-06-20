package db_manager


type DBConfig struct {
	Addr		string
	User		string
	Password	string
	Database	string
}

type User struct {
	tableName struct{} `pg:"users_temp"`

	Id			int64
	Name		string
	Password 	string
}