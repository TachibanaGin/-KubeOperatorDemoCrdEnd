package model

import "time"

type User struct{
	Id 			int
	Username 	string
	Password 	string
	Create_time time.Time
}
