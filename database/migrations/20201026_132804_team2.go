package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Team2_20201026_132804 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Team2_20201026_132804{}
	m.Created = "20201026_132804"

	migration.Register("Team2_20201026_132804", m)
}

// Run the migrations
func (m *Team2_20201026_132804) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE team2(`id` int(11) DEFAULT NULL,`tname` varchar(128) NOT NULL,`owner` varchar(128) NOT NULL,`create_at` int(11) DEFAULT NULL)")
}

// Reverse the migrations
func (m *Team2_20201026_132804) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `team2`")
}
