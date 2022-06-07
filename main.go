package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const fileName = "categories.txt"

// Categories is a model.
type Category struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func main() {

	file, err := os.Open(fileName)
	if err != nil {
		file, err = os.Create(fileName)
		if err != nil {
			panic(err)
		}
	}

	defer file.Close()

	route := gin.Default()

	// route path eg:
	route.POST("/category", CreateCategory)
	route.GET("/categories/", GetAllCategories)
	route.GET("/category/:id", GetCategory)
	route.PUT("/category/:id", UpdateCategory)
	route.DELETE("/category/:id", DeleteCategory)

	// listening at :8085
	route.Run(":8085")
}

func WriteFile() *os.File {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	return file
}

func AppendDataToFile() *os.File {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	return file
}

func ReadFile() []byte {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return data
}

func ReadFilelineByLine(category Category, categories []Category) []Category {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Info("error reading file %s", err)
			break
		}
		newBuff := bytes.NewBuffer(line)
		dec := gob.NewDecoder(newBuff)
		dec.Decode(&category)
		categories = append(categories, category)
	}

	file.Close()
	return categories
}

func removeIndex(categories []Category, index int) []Category {
	return append(categories[:index], categories[index+1:]...)
}

func CreateCategory(c *gin.Context) {
	var category Category
	c.BindJSON(&category)

	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(category)
	if err != nil {
		panic(err)
	}

	file := AppendDataToFile()

	_, err = file.Write(network.Bytes())
	if err != nil {
		panic(err)
	}

	lineBreak := "\n"
	_, err = file.Write([]byte(lineBreak))
	if err != nil {
		panic(err)
	}

	file.Close()

	c.JSON(200, category)
}

func GetAllCategories(c *gin.Context) {
	var categories []Category
	var category Category

	categories = ReadFilelineByLine(category, categories)

	c.JSON(200, categories)
}

func GetCategory(c *gin.Context) {
	id := c.Params.ByName("id")
	var categories []Category
	var category Category

	categories = ReadFilelineByLine(category, categories)

	categoryID, _ := strconv.Atoi(id)
	var check bool
	for _, value := range categories {
		if categoryID == value.ID {
			category = value
			check = true
			break
		}
	}

	if check {
		c.JSON(200, category)
	} else {
		c.JSON(404, "record not found")
	}

}

func UpdateCategory(c *gin.Context) {
	id := c.Params.ByName("id")
	var categories []Category
	var category Category
	var newCategory Category

	c.BindJSON(&newCategory)

	categories = ReadFilelineByLine(category, categories)

	categoryID, _ := strconv.Atoi(id)
	check := false

	for i, prev := range categories {
		if categoryID == prev.ID {
			categories[i] = newCategory
			check = true
			break
		}
	}

	if check {
		uf := WriteFile()
		for _, ctgry := range categories {
			var network bytes.Buffer
			enc := gob.NewEncoder(&network)
			err := enc.Encode(ctgry)
			if err != nil {
				panic(err)
			}

			_, err = uf.Write(network.Bytes())
			if err != nil {
				panic(err)
			}

			lineBreak := "\n"
			_, err = uf.Write([]byte(lineBreak))
			if err != nil {
				panic(err)
			}
		}
		uf.Close()
		c.JSON(200, newCategory)
	} else {
		c.JSON(404, "record not found to update")
	}

}

func DeleteCategory(c *gin.Context) {
	id := c.Params.ByName("id")
	var categories []Category
	var category Category
	var newCategory Category

	c.BindJSON(&newCategory)

	categories = ReadFilelineByLine(category, categories)

	categoryID, _ := strconv.Atoi(id)
	var deleteID int
	var check bool

	for index, prev := range categories {
		if categoryID == prev.ID {
			deleteID = index
			check = true
			break
		}
	}

	if check {
		categories = removeIndex(categories, deleteID)
		uf := WriteFile()
		for _, ctgry := range categories {
			var network bytes.Buffer
			enc := gob.NewEncoder(&network)
			err := enc.Encode(ctgry)
			if err != nil {
				panic(err)
			}

			_, err = uf.Write(network.Bytes())
			if err != nil {
				panic(err)
			}

			lineBreak := "\n"
			_, err = uf.Write([]byte(lineBreak))
			if err != nil {
				panic(err)
			}
		}
		uf.Close()

		c.JSON(200, "record deleted successfully!")
	} else {
		c.JSON(404, "record not found")
	}
}
