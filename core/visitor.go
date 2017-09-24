package core

// Visitor is passed a Secret
type Visitor interface {
	Visit(secret Secret)
}

// Walker is passed a Visitor
type Walker interface {
	Walk(visitor Visitor)
}
