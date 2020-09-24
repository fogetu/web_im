package main

import (
	"github.com/astaxie/beego/migration"
)

// bee generate migration user -driver=mysql -fields="name:string,age:int"
// bee migrate -driver=mysql  -conn="username:password@tcp(127.0.0.1:3306)/db"

// DO NOT MODIFY
type User_20200308_205046 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &User_20200308_205046{}
	m.Created = "20200308_205046"

	migration.Register("User_20200308_205046", m)
}

// Run the migrations
func (m *User_20200308_205046) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE user(`id` int(11) NOT NULL AUTO_INCREMENT,`name` varchar(128) NOT NULL,`age` int(11) DEFAULT NULL,PRIMARY KEY (`id`))")
}

// Reverse the migrations
func (m *User_20200308_205046) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `user`")
}
