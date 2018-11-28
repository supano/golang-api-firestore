package main

import (
	"context"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"google.golang.org/api/option"
)

func main() {

	// firebase cloud firestore
	ctx := context.Background()
	opt := option.WithCredentialsFile("./golang-75d13-firebase-adminsdk-k4zs4-40015f5c5e.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error when create new app : %v", err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error when initial new firestore obj : %v", err)
	}

	e := echo.New()
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		log.Printf("%s\n", resBody)
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.POST("/api/login", func(c echo.Context) error {
		u := new(User)
		c.Bind(u)
		query := client.Collection("users").Where("email", "==", u.Email).Where("password", "==", u.Password).Limit(1).Documents(ctx)
		doc, err := query.Next()
		if err == nil && doc != nil {
			return c.JSON(http.StatusOK, Message{Text: "Welcome"})
		}
		return c.JSON(http.StatusUnauthorized, Message{Text: "Wrong Email or Password"})
	})

	e.Logger.Fatal(e.Start(":1323"))

	// _, _, err = client.Collection("users").Add(ctx, map[string]interface{}{
	// 	"email":    "supano1995@gmail.com",
	// 	"password": "123456",
	// })
	// if err != nil {
	// 	log.Fatalf("error when add new data : %v", err)
	// }

	// iter := client.Collection("users").Documents(ctx)
	// for {
	// 	doc, err := iter.Next()
	// 	if err == iterator.Done {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatalf("Failed to iterate: %v", err)
	// 	}
	// 	fmt.Println(doc.Data())
	// }

	// defer client.Close()
}

type Message struct {
	Text string `json:"text"`
}

type User struct {
	Email    string `json:"email" form:"email" query:"email"`
	Password string `json:"password" form:"password" query:"password"`
}
