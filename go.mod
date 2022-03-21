module github.com/crossplane-contrib/provider-ibm-cloud

go 1.16

require (
	github.com/IBM-Cloud/bluemix-go v0.0.0-20201019071904-51caa09553fb
	github.com/IBM/cloudant-go-sdk v0.0.34
	github.com/IBM/eventstreams-go-sdk v1.1.0
	github.com/IBM/experimental-go-sdk v0.0.0-20210112204617-192fc5b15655
	github.com/IBM/go-sdk-core v1.1.0
	github.com/IBM/go-sdk-core/v4 v4.10.0
	github.com/IBM/ibm-cos-sdk-go v1.7.0
	github.com/IBM/platform-services-go-sdk v0.17.18
	github.com/crossplane/crossplane-runtime v0.15.0
	github.com/crossplane/crossplane-tools v0.0.0-20201007233256-88b291e145bb
	github.com/go-openapi/strfmt v0.20.1
	github.com/golang-jwt/jwt/v4 v4.4.0 // indirect
	github.com/google/go-cmp v0.5.5
	github.com/jeremywohl/flatten v1.0.1
	github.com/pkg/errors v0.9.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	sigs.k8s.io/controller-runtime v0.9.2
	sigs.k8s.io/controller-tools v0.2.4
)

replace github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt/v4 v4.3.0

replace github.com/form3tech-oss/jwt-go => github.com/golang-jwt/jwt/v4 v4.2.0
