package helper

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func QueueSidValidation() schema.SchemaValidateFunc {
	return validation.StringMatch(regexp.MustCompile("^QU[0-9a-fA-F]{32}$"), "")
}