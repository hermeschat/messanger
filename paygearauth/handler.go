package paygearauth

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
)

var (
	ServiceName           string
	ApplicationServiceURL string
	HttpRequestAuthToken  string
)

func Init(serviceName string, appServiceURL string, httpRequestAuthToken string) {
	ServiceName = serviceName
	ApplicationServiceURL = appServiceURL
	HttpRequestAuthToken = httpRequestAuthToken
}

// GetAuthentication check and get auth from given token
func GetAuthentication(token string, accountID string) (*Identity, error) {
	if token == "" || len(token) < 7 {
		return nil, UnauthorizedError{}
	}
	if token[:7] == "bearer " {
		token = token[7:]
	}
	// get jwt data
	jwtDataStr := strings.Split(token, ".")[1]
	bytes, err := base64Decode(jwtDataStr)
	if err != nil {
		log.Println(err)
		return nil, UnauthorizedError{}
	}
	properties := make(map[string]interface{})
	err = json.Unmarshal(bytes, &properties)
	if err != nil {
		log.Println(err)
		return nil, UnauthorizedError{}
	}
	// find secret key of application
	secretKey := ""
	// clientID := ""
	if app, ok := properties["app"]; ok {
		applicationID := app.(string)
		_, secretKey, err = GetApplicationAuthKeys(applicationID)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Invalid JWT")
	}
	claims, err := ValidateJWT(token, secretKey)
	if err != nil {
		//log.Println(err)
		return nil, err
	}
	serviceAvailable := false
	var service IdentityService
	if svc, ok := claims["svc"]; ok {
		if services, ok := svc.(map[string]interface{}); ok {
			if _service, ok := services[ServiceName]; ok {
				if serviceMap, ok := _service.(map[string]interface{}); ok {
					serviceAvailable = true
					perm, ok := serviceMap["perm"]
					var permissions int64 = 0
					if ok {
						permissions = int64(perm.(float64))
					}
					service = IdentityService{
						Permissions: permissions,
					}
				}
			}
		}
	}
	_ = serviceAvailable
	// nameID := claims["nameid"].(string)
	userID := claims["id"].(string)

	var rolesMap map[string]bool
	if _roles, ok := claims["role"]; ok {
		rolesMap = make(map[string]bool)
		if roles, ok := _roles.([]interface{}); ok {
			for _, role := range roles {
				rolesMap[role.(string)] = true
			}
		}
	}

	//if !serviceAvailable {
	//	// check if is rostam
	//	if _, ok := rolesMap["rostam"]; !ok {
	//		return nil, errors.New("is not rostam")
	//	}
	//}

	merchantRoles := map[string][]string{}
	if _merchantRoles, ok := claims["merchant_roles"]; ok {
		if merchantRolesMapFace, ok := _merchantRoles.(map[string]interface{}); ok {
			for merchantRolesKey, _mRoles := range merchantRolesMapFace {
				if mRoles, ok := _mRoles.([]interface{}); ok {
					merchantRoles[merchantRolesKey] = []string{}
					for _, _role := range mRoles {
						merchantRoles[merchantRolesKey] = append(merchantRoles[merchantRolesKey], _role.(string))
					}
				}
			}
		}
	}

	var AppID string
	AppID = claims["app"].(string)
	if accountID != "" && userID != accountID {
		accountAccessGranted := false
		exists := false
		if rolesMap != nil {
			if _, ok := rolesMap["zeus"]; ok {
				exists = true
			} else if _, ok := rolesMap["rostam"]; ok {
				exists = true
			}
		}
		if exists {
			accountAccessGranted = true
		}
		if !accountAccessGranted {
			if _, ok := merchantRoles[accountID]; ok {
				accountAccessGranted = true
			} else {
				return nil, ForbiddenError{}
			}
		} else {
			return nil, ForbiddenError{}
		}
	}
	identity := Identity{Token: token, ID: userID, Roles: rolesMap, MerchantRoles: merchantRoles, Service: service, AppId: AppID}
	return &identity, nil
}
