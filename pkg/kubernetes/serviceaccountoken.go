package kubernetes

import (
	"reflect"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/projectdiscovery/gologger"
)

// ServiceAccountTokenClaims is the structure used to unmarshal Kubernetes JWT Claims
type ServiceAccountTokenClaims struct {
	Kubernetes struct {
		Namespace string
		Node      struct {
			Name string
		}
		Pod struct {
			Name string
		}
		ServiceAccount struct {
			Name string
		}
	} `json:"kubernetes.io"`
	jwt.RegisteredClaims
}

// GrabNamespaceFromToken returns the namespace stored in a Kubernetes JWT Claim
func GrabNamespaceFromToken(sa string) (ns string) {
	ns = "default"
	if !isJWT(sa) {
		gologger.Warning().Msg("service account token is probably not a JWT")
		// always return the default namespace
		return
	}

	token, err := jwt.ParseWithClaims(sa, &ServiceAccountTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("a"), nil
	})

	if err != nil && reflect.TypeOf(err) == reflect.TypeOf(jwt.ErrTokenSignatureInvalid) {
		gologger.Error().Msgf("could not parse service account token claims. err: %v\n", err)
		return
	}

	claims := token.Claims.(*ServiceAccountTokenClaims)
	ns = claims.Kubernetes.Namespace

	if ns == "" {
		token, _ = jwt.Parse(sa, func(t *jwt.Token) (interface{}, error) {
			return []byte("a"), nil
		})
		ns = token.Claims.(jwt.MapClaims)["kubernetes.io/serviceaccount/namespace"].(string)
	}

	gologger.Debug().Msgf("found service account token namespace: %v\n", ns)
	return
}

func isJWT(tk string) bool {
	splitted := strings.Split(tk, ".")
	if len(splitted) < 3 && strings.Count(tk, ".") != 2 {
		return false
	}

	return true
}
