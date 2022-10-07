package model

import (
	"fmt"
	"time"
)

type User struct {
	Id int `gorm:"primaryKey"`
	Name string `gorm:"unique"`
	Phone string `gorm:"unique"`
	Sex string 
	Add_time int
	Status int
}

func (user User) TableName() string {
	return "im_user"
}

//注册用户
func (this *User) Register(name string, phone string, sex string) (int, error, int){
	this = &User{
		Name: name,
		Phone: phone,
		Sex: sex,
		Add_time: int(time.Now().Unix()),
		Status: 1,
	}
	result := Db.Create(this)
	//主键id， error， 行数
	return this.Id, result.Error, int(result.RowsAffected)
}

func (this *User) GetInfoById(id int) *User{
	// Db.Where(&User{Name:"sc2"}).First(this)
	Db.First(this, id)
	return this
}

func (this *User) GetLists() {
	var users []User
	Db.Where("id > ? AND name != ?", 0, "abc").Find(&users)
	fmt.Println(users)
}