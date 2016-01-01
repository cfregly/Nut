package specification

import (
	"strings"
)

func ParentName(id string) string {
	// docker FROM entry: org/repo:version -> org-repo_version
	orgConverted := strings.Replace(id, "/", "-", 1)
	repoConverted := strings.Replace(orgConverted, ":", "_", 1)
	return repoConverted
}
