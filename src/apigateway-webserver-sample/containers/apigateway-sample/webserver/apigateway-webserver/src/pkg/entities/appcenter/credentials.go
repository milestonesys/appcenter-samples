package appcenter

import (
	"fmt"
	"io"
	"os"

	"apigateway-webserver/src/pkg/constants/enums"
)

func ReadCredentialsFlowFiles(username, password string, flowType enums.CredentialsFlowType) (string, string, error) {
	if flowType == enums.ClientCredentialsFlow {
		// The username should be the same username we set in the values.yaml for the helm charts configurations
		// Here is a copy of the values file to show how the secret is shared with the service
		// clientCredentialsFlow:
		//   secret: "app-client"
		//   clientID: "09da6440-e308-31a0-963e-1af823e76a33"
		//   clientName: "apigateway-sample-service"
		//   ...

		// At the deployment.yaml file you can see how the secret is mounted and mapped to the /etc/app-client/ directory
		var err error
		clientID, err := readFileContent("/etc/app-client/client-id")
		if err != nil {
			return "", "", err
		}

		clientSecret, err := readFileContent("/etc/app-client/client-secret")
		if err != nil {
			return "", "", err
		}
		return clientID, clientSecret, nil
	}
	return username, password, nil
}

func readFileContent(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return string(content), nil
}
