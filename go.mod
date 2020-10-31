module github.com/crossplane-contrib/provider-ibm-cloud

go 1.13

require (
	github.com/IBM-Cloud/bluemix-go v0.0.0-20201019071904-51caa09553fb
	github.com/IBM/go-sdk-core v1.1.0
	github.com/IBM/platform-services-go-sdk v0.14.0
	github.com/crossplane/crossplane-runtime v0.10.0
	github.com/crossplane/crossplane-tools v0.0.0-20201007233256-88b291e145bb
	github.com/go-openapi/strfmt v0.19.5
	github.com/google/go-cmp v0.5.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
	sigs.k8s.io/controller-runtime v0.6.2
	sigs.k8s.io/controller-tools v0.2.4
)
