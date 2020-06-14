package main

import (
	"gopkg.in/mgo.v2/bson"
)

const messageFetchSize = 10

type Poem struct {
	ID      bson.ObjectId `bson:"_id" json:"id"`
	Title   string        `bson:"title" json:"title"`
	Author  string        `bson:"author" json:"author"`
	Content string        `bson:"content" json:"content"`
}
