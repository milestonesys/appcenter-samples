package enums

import "fmt"

type CredentialsFlowType int

const (
	LoginForm CredentialsFlowType = iota + 1
	ClientCredentialsFlow
)

var (
	credentialsFlowTypeMap = map[string]CredentialsFlowType{
		"LoginForm":             LoginForm,
		"ClientCredentialsFlow": ClientCredentialsFlow,
	}
)

func (c CredentialsFlowType) String() string {
	return [...]string{"LoginForm", "ClientCredentialsFlow"}[c-1]
}

func ParseCredentialsFlowType(str string) (CredentialsFlowType, error) {
	c, ok := credentialsFlowTypeMap[str]
	if !ok {
		return 0, fmt.Errorf("invalid CredentialsFlowType: %s", str)
	}
	return c, nil
}

func GetCredentialsFlowTypes() []string {
	return []string{
		LoginForm.String(),
		ClientCredentialsFlow.String(),
	}
}
