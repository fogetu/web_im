package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Team_20201026_132802 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Team_20201026_132802{}
	m.Created = "20201026_132802"

	migration.Register("Team_20201026_132802", m)
}

// Run the migrations
func (m *Team_20201026_132802) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE team(`id` int(11) DEFAULT NULL,`tname` varchar(128) NOT NULL,`owner` varchar(128) NOT NULL,`create_at` int(11) DEFAULT NULL)")
}

// Reverse the migrations
func (m *Team_20201026_132802) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `team`")
}
