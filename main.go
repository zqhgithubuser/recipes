package main

import (
    "encoding/json"
    "github.com/gin-gonic/gin"
    "github.com/rs/xid"
    "io/ioutil"
    "net/http"
    "time"
)

var recipes []Recipe

type Recipe struct {
    ID           string    `json:"id"`
    Name         string    `json:"name"`
    Tags         []string  `json:"tags"`
    Ingredients  []string  `json:"ingredients"`
    Instructions []string  `json:"instructions"`
    PublishedAt  time.Time `json:"published_at"`
}

func NewRecipeHandler(c *gin.Context) {
    var recipe Recipe
    if err := c.ShouldBindJSON(&recipe); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    recipe.ID = xid.New().String()
    recipe.PublishedAt = time.Now()

    recipes = append(recipes, recipe)
    c.JSON(http.StatusOK, recipe)
}

func ListRecipesHandler(c *gin.Context) {
    c.JSON(http.StatusOK, recipes)
}

func UpdateRecipeHandler(c *gin.Context) {
    id := c.Param("id")
    var recipe Recipe
    if err := c.ShouldBindJSON(&recipe); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    index := -1
    for i := 0; i < len(recipes); i++ {
        if recipes[i].ID == id {
            index = i
        }
    }

    if index == -1 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
        return
    }

    recipes[index] = recipe
    c.JSON(http.StatusOK, recipe)
}

func DeleteRecipeHandler(c *gin.Context) {
    id := c.Param("id")

    index := -1
    for i := 0; i < len(recipes); i++ {
        if recipes[i].ID == id {
            index = i
        }
    }

    if index == -1 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
        return
    }

    recipes = append(recipes[:index], recipes[index+1:]...)
    c.JSON(http.StatusOK, gin.H{"message": "Recipe has been deleted"})
}

func GetRecipeHandler(c *gin.Context) {
    id := c.Param("id")
    for i := 0; i < len(recipes); i++ {
        if recipes[i].ID == id {
            c.JSON(http.StatusOK, recipes[i])
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
}

func init() {
    recipes = make([]Recipe, 0)
    file, _ := ioutil.ReadFile("recipes.json")
    _ = json.Unmarshal(file, &recipes)
}

func main() {
    server := gin.Default()
    server.POST("/recipes", NewRecipeHandler)
    server.GET("/recipes", ListRecipesHandler)
    server.PUT("/recipes/:id", UpdateRecipeHandler)
    server.DELETE("/recipes/:id", DeleteRecipeHandler)
    server.GET("/recipes/:id", GetRecipeHandler)
    server.Run()
}
