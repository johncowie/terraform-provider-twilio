package autopilot

import (
	"fmt"
	"time"

	"github.com/RJPearson94/terraform-provider-twilio/twilio/common"
	"github.com/RJPearson94/terraform-provider-twilio/twilio/utils"
	"github.com/RJPearson94/twilio-sdk-go/service/autopilot/v1/assistant/task/fields"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAutopilotTaskField() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutopilotTaskFieldCreate,
		Read:   resourceAutopilotTaskFieldRead,
		Delete: resourceAutopilotTaskFieldDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"sid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_sid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"assistant_sid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"task_sid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"unique_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"field_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAutopilotTaskFieldCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*common.TwilioClient).Autopilot

	createInput := &fields.CreateFieldInput{
		UniqueName: d.Get("unique_name").(string),
		FieldType:  d.Get("field_type").(string),
	}

	createResult, err := client.Assistant(d.Get("assistant_sid").(string)).Task(d.Get("task_sid").(string)).Fields.Create(createInput)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to create autopilot task field: %s", err.Error())
	}

	d.SetId(createResult.Sid)
	return resourceAutopilotTaskFieldRead(d, meta)
}

func resourceAutopilotTaskFieldRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*common.TwilioClient).Autopilot

	getResponse, err := client.Assistant(d.Get("assistant_sid").(string)).Task(d.Get("task_sid").(string)).Field(d.Id()).Fetch()
	if err != nil {
		if utils.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Failed to read autopilot task field: %s", err.Error())
	}

	d.Set("sid", getResponse.Sid)
	d.Set("account_sid", getResponse.AccountSid)
	d.Set("assistant_sid", getResponse.AssistantSid)
	d.Set("task_sid", getResponse.TaskSid)
	d.Set("unique_name", getResponse.UniqueName)
	d.Set("field_type", getResponse.FieldType)
	d.Set("date_created", getResponse.DateCreated.Format(time.RFC3339))

	if getResponse.DateUpdated != nil {
		d.Set("date_updated", getResponse.DateUpdated.Format(time.RFC3339))
	}

	d.Set("url", getResponse.URL)
	return nil
}

func resourceAutopilotTaskFieldDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*common.TwilioClient).Autopilot

	if err := client.Assistant(d.Get("assistant_sid").(string)).Task(d.Get("task_sid").(string)).Field(d.Id()).Delete(); err != nil {
		return fmt.Errorf("Failed to delete autopilot task field: %s", err.Error())
	}
	d.SetId("")
	return nil
}