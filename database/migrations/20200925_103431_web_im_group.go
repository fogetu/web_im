package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type WebImGroup_20200925_103431 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &WebImGroup_20200925_103431{}
	m.Created = "20200925_103431"

	migration.Register("WebImGroup_20200925_103431", m)
}

// Run the migrations
func (m *WebImGroup_20200925_103431) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update

}

// Reverse the migrations
func (m *WebImGroup_20200925_103431) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
