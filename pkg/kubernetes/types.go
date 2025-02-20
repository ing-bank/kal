package kubernetes

// ApiVerbs is the official list of API Verbs accepted by Kubernetes API
var ApiVerbs = []string{
	"create",
	"get",
	"list",
	"watch",
	"update",
	"patch",
	"delete",
	"deletecollection",
	"impersonate",
	"bind",
	"approve",
	"escalate",
}
