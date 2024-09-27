package main

import (
	"strconv"

	"log"

	"net/http"

	"github.com/gin-gonic/gin"

	"example.com/personweb/models"
)

func main() {
	err := models.ConnectDatabase();
	checkErr(err);

	r := gin.Default();

	r.GET("person", getPersons);
	r.GET("person/:id", getPersonById);
	r.POST("person", addPerson);
	r.PUT("person/:id", updatePerson);
	r.DELETE("person/:id", deletePerson);

	r.Run()
}

func getPersons(c *gin.Context) {

	persons, err := models.GetPersons(1005);
	checkErr(err)

	if persons == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": persons})
	}
}

func getPersonById(c *gin.Context) {

	id := c.Param("id")

	person, err := models.GetPersonById(id)
	checkErr(err)
	// If the name is blank we can assume nothing is found
	if person.FirstName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": person})
	}
}

func addPerson(c *gin.Context) {

	var json models.Person

	if err := c.ShouldBindJSON(&json); err != nil { 
		// Attempts to bind the JSON body of the request to the json variable
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := models.AddPerson(json) // Write to database

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func updatePerson(c *gin.Context) {

	var json models.Person

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	personId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.UpdatePerson(json, personId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func deletePerson(c *gin.Context) {

	personId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.DeletePerson(personId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}