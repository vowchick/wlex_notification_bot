package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

func check(err error) int {
	if err != nil {
		fmt.Println("smth wrong")
		return -1
	}
	return 0
}

func checkForNew(data []map[string]string) {
	jsonFile, err := os.Open("listings.json")
	if check(err) != 0 {
		return
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if check(err) != 0 {
		return
	}

	var oldData []map[string]string
	if len(byteValue) != 0 {
		err = json.Unmarshal(byteValue, &oldData)
		if check(err) != 0 {
			return
		}
	}
	for i, el := range data {
		has := false
		for _, el2 := range oldData {
			if (el["sellerWallet"] == el2["sellerWallet"]) &&
				(el["emailLogin"] == el2["emailLogin"]) &&
				(el["emailPassword"] == el2["emailPassword"]) &&
				(el["discordPassword"] == el2["discordPassword"]) &&
				(el["discordLogin"] == el2["discordLogin"] &&
					el["projectId"] == el2["projectId"]) {
				has = true
			}
		}
		if !has {
			myLog("sellerWallet: " + el["sellerWallet"] + "\n" +
				"emailLogin: " + el["emailLogin"] + "\n" +
				"emailPassword: " + el["emailPassword"] + "\n" +
				"discordLogin: " + el["discordLogin"] + "\n" +
				"discordPassword: " + el["discordPassword"] + "\n" +
				"projectId: " + el["projectId"] + "\n" +
				"listingId: " + strconv.Itoa(i) + "\n")
		}
	}
	file, _ := json.MarshalIndent(data, "", " ")

	_ = ioutil.WriteFile("listings.json", file, 0644)
}

func readFromFirebase() {
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: "https://wlex-bot-default-rtdb.firebaseio.com/",
	}
	// Fetch the service account key JSON file contents
	opt := option.WithCredentialsFile("path_to_private_key.json")

	// Initialize the app with a service account, granting admin privileges
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		fmt.Println("Error initializing app:", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		fmt.Println("Error initializing database client:", err)
	}

	// As an admin, the app has access to read and write all data, regradless of Security Rules
	ref := client.NewRef("/sellRequests")
	var data []map[string]string
	if err := ref.Get(ctx, &data); err != nil {
		fmt.Println("Error reading from database:", err)
	}
	checkForNew(data)
}