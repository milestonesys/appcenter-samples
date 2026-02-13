package appcenter

import (
	"fmt"
	"os"

	"apigateway-webserver/src/pkg/constants/enums"
)

func ReadCredentialsFromEnv(username, password string, flowType enums.CredentialsFlowType) (string, string, error) {
	if flowType == enums.ClientCredentialsFlow {
		// For credentials to be available for the app the clientCredentialsFlow must be defined in app-definition.yaml like this:
		//
		// credentials:
		//  clientCredentialsFlow:
		//    clientName: "apigateway-sample-service"
		//    clientScopes: [ "managementserver" ]

		clientID, available := os.LookupEnv("CCF_CLIENT_ID")
		if !available {
			return "", "", fmt.Errorf("environment variable CCF_CLIENT_ID not set")
		}

		clientSecret, available := os.LookupEnv("CCF_CLIENT_SECRET")
		if !available {
			return "", "", fmt.Errorf("environment variable CCF_CLIENT_SECRET not set")
		}
		return clientID, clientSecret, nil
	}
	return username, password, nil
}
