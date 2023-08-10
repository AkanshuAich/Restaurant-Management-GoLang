package controller

import (
	"context"
	"fmt"
	"golang-restaurant-management/database"
	"golang-restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/mongo"
)
var foodCollection *mongo.Collection = database.OpenCollection(databse.Client, "food")
var validate = validator.New()
func GetFoods() gin.HandlerFunc{
	return func(c *gin.Context){

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		ecordPerPage, err := strvconv.Atoi(c.Query("recordPerPage"))
		recordPerPage := 
		if err!=nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strvconv.Atoi(c.Query("page"))
		if err !=nil || page <1{
			page=1
		}

		startIndex := (page-1)*recordPerPage
		startIndex, err = strvconv.Atoi(c.Query("startIndex"))

		matchStage : bson.D{{"$match", bson.D{{}}}}
		groupStage : bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}},{"total_count", bson.D{{"$sum, 1"}}},{"data", bson.D{{"$push","$$ROOT"}}}}}}
		projectStage := bson.D{
			{
				"$project", bson.D{
					{"_id",0},
					{"total_count", 1},
					{"food_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
				}}}

		result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})

		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError, gin.h{"error":"error occurred while listing food items"})
		}
		var allFoods []bson.M
		if err = result.All(ctx, &allFoods); err != nil{
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allFoods[0])
		
	}
}

func GetFood() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		foodId := c.Param("food_id")
		var food models.Food

		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodID}).Decode(&food)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the food item"})
		}
		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu
		var food models.Food

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(food)
		if validationErr !=nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id})
		defer cancel()
		if err != nil {
			msg := fmt.Sprintf("menu was not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		food.Created_at, _=time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_at, _=time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID()
		food.Food_id = food.ID.Hex()
		var num = toFixed(*food.Price, 2)
		food.Price = &num

		result, insertErr := foodCollection.InsertOne(ctx, menu)
		if insertErr !=nil{
			msg := fmt.Sprintf("Food item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func round(num float64) int {
return int(num + math.Copysign(0.5, num))

}

func toFixed(num float64, precision int) float64 {
	output := math.pow(10, float64(precision))
	return float64(round(num*output))/output
}

func UpdateFood() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu
		var food models.Food

		foodId := c.Param("food_id")

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error", err.Error()})
			return
		}

		var updateObj primitive.D

		if food.Name != nil {
			updateObj = append(updateObj, bson.E("name", food.Name))

		}

		if food.Price != nil {
			updateObj = append(updateObj, bson.E("price", food.Price))

		}

		if food.Food_image != nil {

			updateObj = append(updateObj, bson.E("food_image", food.Food_image))
		}

		if food.Menu_id != nil {
			err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id})
			defer cancel()
			if err!=nil{
				msg:= fmt.Sprintf("message:Menu was not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
				return
			}
			updateObj = append(updateObj, bson.E{"menu", food.Price})

		}

		food.Updated_at, _ =time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", menu.Updated_at})

		upsert := true
		filter := bson.M{"food_id": foodID}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := foodCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err!= nil {
			msg := fmt.Sprint("foot item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}

		
		c.JSON(http.StatusOK, result)

	}
}

