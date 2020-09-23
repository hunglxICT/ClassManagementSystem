//Model/UserModel.go

package Models

import "mime/multipart"

type Account struct {
	Id			int		`form:"Id"`
	Username	string	`form:"Username"`
	Password	string	`form:"Password"`
	Firstname	string	`form:"Firstname"`
	Lastname	string	`form:"Lastname"`
	Email		string	`form:"Email"`
	Tel			string	`form:"Tel"`
	Avatar		string	`form:"Avatar"`
	Role		int		`form:"Role"`
	Status		int		`form:"Status"`
	Create_date	int64	`form:"create_date"`
}

func (b *Account) TableName() string {
	return "Account"
}


type Class struct {
	Id			int		`form:"id"`
	Teacherid	int		`form:"teacherid"`
	Classname	string	`form:"classname"`
	Status		int		`form:"status"`
	Create_date	int64	`form:"create_date"`
}

func (b *Class) TableName() string {
	return "Class"
}


type Enroll struct {
	Id			int		`form:"id"`
	Studentid	int		`form:"studentid"`
	Classid		int		`form:"classid"`
	Status		int		`form:"status"`
	Create_date	int64	`form:"create_date"`
}

func (b *Enroll) TableName() string {
	return "Enroll"
}


type Exercises struct {
	Id			int		`form:"id"`
	Classid		int		`form:"classid"`
	Title		string	`form:"title"`
	Description	string	`form:"description"`
	Link		string	`form:"link"`
	Status		int		`form:"status"`
	Create_date	int64	`form:"create_date"`
}

func (b *Exercises) TableName() string {
	return "Exercises"
}


type Messenger struct {
	Id			int		`form:"id"`
	Senderid	int		`form:"senderid"`
	Receiverid	int		`form:"receiverid"`
	Status		int		`form:"status"`
	Message		string	`form:"message"`
	Create_date	int64	`form:"create_date"`
}

func (b *Messenger) TableName() string {
	return "Messenger"
}


type Submission struct {
	Id			int		`form:"id"`
	Studentid	int		`form:"studentid"`
	Exerciseid 	int		`form:"exerciseid"`
	Description	string	`form:"description"`
	Link		string	`form:"link"`
	Create_date	int64	`form:"create_date"`
}

func (b *Submission) TableName() string {
	return "Submission"
}

type MultipartForm struct {
	File		*multipart.File	`form:"file"`
	Data		*multipart.Form	`form:"data"`
}