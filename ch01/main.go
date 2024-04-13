package main

import (
	"encoding/xml"

	"github.com/gin-gonic/gin"
)

func IndexHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hello world",
	})
}

func HelloHandler(c *gin.Context) {
	name := c.Params.ByName("name")
	c.JSON(200, gin.H{
		"message": "hello " + name,
	})
}

type Person struct {
	XMLName   xml.Name `xml:"person"`
	FirstName string   `xml:"firstName,attr"`
	LastName  string   `xml:"lastName,attr"`
}

func XmlHandler(c *gin.Context) {
	c.XML(200, Person{FirstName: "Mohamed",
		LastName: "Labouardy"})
}

func main() {
	router := gin.Default()
	router.GET("/hello/:name", HelloHandler)
	router.GET("/xml", XmlHandler)
	router.GET("/", IndexHandler)
	router.Run()
}
