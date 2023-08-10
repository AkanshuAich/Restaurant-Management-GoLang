package controller

import "github.com/gin-gonic/gin"

func GetMenus() gin.HandlerFunc{
	return func(c *gin.Context){

	}
}

func GetMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		menuId := c.Param("menu_id")
		var food models.Menu

		err := foodCollection.FindOne(ctx, bson.M{"menu_id": foodID}).Decode(&menu)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the menu"})
		}
		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		var menu models.Menu
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(menu)
		if validationErr !=nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error})
			return
		}

		menu.Created_at, _=time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_at, _=time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		result, insertErr := menuCollection.InsertOne(ctx, menu)
		if insertErr !=nil{
			msg := fmt.Sprintf("Menu item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return start.After(time.Now()) && end.After(start)
}

func UpdateMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu module.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":berr.Error()})
			return
		}

		menuId := c.Param("menu_id")
		filter := bson.M{"menu_id":menuId}

		var updateObj primitive.D

		if menu.Start_Date != nil && menu.End_Date!=nil{
			if !inTimeSpan(*menu.Start_Date, *menu.End_Date, time.Now()){
				msg := "kindly retype the time"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				defer cancel()
				return
			}

			updateObj = append(updateObj, bson.E{"start_date", menu.Start_Date})
			updateObj = append(updateObj, bson.E{"end_date", menu.End_Date})

			if menu.Name != ""{
				updateObj = append(updateObj, bson.E{"name", menu.Name})
			}
			if menu.Category != ""{
				updateObj = append(updateObj, bson.E{"name", menu.Category})
			}

			menu.Updated_at, _ =time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updateObj = append(updateObj, bson.E{"updated_at", menu.Updated_at})

			upsert:=true

			opt := options.UpdateOptions{
				Upsert : &upsert,
			}

			result, err := menu.Collection.UpdateOne(
				ctx,
				filter,
				bson.D{
					{"$set",updateObj},
				},
				&opt,
			)

			if err!= nil {
				msg := "Menu update failed"
				c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			}

			defer cancel()
			c.JSON(http.StatusOK, result)
		}

	}
}

