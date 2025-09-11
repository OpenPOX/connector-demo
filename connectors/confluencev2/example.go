package confluencev2

import (
	"fmt"
)

func Example() {
	var (
		clientID     = "your_client_id"
		clientSecret = "your_client_secret"
		state        = "your_customer_id"
	)
	authorizationURL := GenerateAuthorizationURL(state)
	fmt.Println(authorizationURL)
	code := "the_code_from_callback"
	// verify state
	oauth2 := NewOAuth2(clientID, clientSecret, "localhost:8080/AtlassianCallback")
	token, err := oauth2.AuthorizationCode(code)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	resources, err := oauth2.AccessibleResources(token.AccessToken)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	apiclient := NewAPI(token.AccessToken, resources[0].ID)
	pages, err := apiclient.GetPages(nil)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(pages)
	if len(pages.Results) > 0 {
		fmt.Println(pages.Results[0].Body)
	}
	// ...
	newToken, err := oauth2.RefreshToken(token.RefreshToken)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(newToken)
}
