package sdk

import (
	"fmt"
	"strings"
)

type WIFType string

const (
	WIFTypeAWS   WIFType = "AWS"
	WIFTypeAzure WIFType = "AZURE"
	WIFTypeGCP   WIFType = "GCP"
	WIFTypeOIDC  WIFType = "OIDC"
)

func ToWIFTypeType(s string) (WIFType, error) {
	switch strings.ToUpper(s) {
	case string(WIFTypeAWS):
		return WIFTypeAWS, nil
	case string(WIFTypeAzure):
		return WIFTypeAzure, nil
	case string(WIFTypeGCP):
		return WIFTypeGCP, nil
	case string(WIFTypeOIDC):
		return WIFTypeOIDC, nil
	default:
		return "", fmt.Errorf("invalid WIF type: %s", s)
	}
}

type SecondaryRolesOption string

const (
	SecondaryRolesOptionDefault SecondaryRolesOption = "DEFAULT"
	SecondaryRolesOptionNone    SecondaryRolesOption = "NONE"
	SecondaryRolesOptionAll     SecondaryRolesOption = "ALL"
)

func ToSecondaryRolesOption(s string) (SecondaryRolesOption, error) {
	switch strings.ToUpper(s) {
	case string(SecondaryRolesOptionDefault):
		return SecondaryRolesOptionDefault, nil
	case string(SecondaryRolesOptionNone):
		return SecondaryRolesOptionNone, nil
	case string(SecondaryRolesOptionAll):
		return SecondaryRolesOptionAll, nil
	default:
		return "", fmt.Errorf("invalid secondary roles option: %s", s)
	}
}

var ValidSecondaryRolesOptionsString = []string{
	string(SecondaryRolesOptionDefault),
	string(SecondaryRolesOptionNone),
	string(SecondaryRolesOptionAll),
}

type UserType string

const (
	UserTypePerson        UserType = "PERSON"
	UserTypeService       UserType = "SERVICE"
	UserTypeLegacyService UserType = "LEGACY_SERVICE"
)

func ToUserType(s string) (UserType, error) {
	switch strings.ToUpper(s) {
	case string(UserTypePerson):
		return UserTypePerson, nil
	case string(UserTypeService):
		return UserTypeService, nil
	case string(UserTypeLegacyService):
		return UserTypeLegacyService, nil
	default:
		return "", fmt.Errorf("invalid user type: %s", s)
	}
}

var AllUserTypes = []UserType{
	UserTypePerson,
	UserTypeService,
	UserTypeLegacyService,
}

var AcceptableUserTypes = map[UserType][]string{
	UserTypePerson:        {"", string(UserTypePerson)},
	UserTypeService:       {string(UserTypeService)},
	UserTypeLegacyService: {string(UserTypeLegacyService)},
}
