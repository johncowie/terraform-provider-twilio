package autopilot

import (
	"fmt"
	"strings"
	"time"

	"github.com/RJPearson94/terraform-provider-twilio/twilio/common"
	"github.com/RJPearson94/terraform-provider-twilio/twilio/utils"
	"github.com/RJPearson94/twilio-sdk-go/service/autopilot/v1/assistant/webhook"
	"github.com/RJPearson94/twilio-sdk-go/service/autopilot/v1/assistant/webhooks"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const webhookEventsSeperator = " "

func resourceAutopilotWebhook() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutopilotWebhookCreate,
		Read:   resourceAutopilotWebhookRead,
		Update: resourceAutopilotWebhookUpdate,
		Delete: resourceAutopilotWebhookDelete,
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
			"unique_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"events": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"webhook_url": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"webhook_method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"GET",
					"POST",
				}, false),
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

func resourceAutopilotWebhookCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*common.TwilioClient).Autopilot

	createInput := &webhooks.CreateWebhookInput{
		UniqueName:    d.Get("unique_name").(string),
		WebhookURL:    d.Get("webhook_url").(string),
		Events:        utils.ConvertSliceToSeperatedString(d.Get("events").([]interface{}), webhookEventsSeperator),
		WebhookMethod: utils.OptionalString(d, "webhook_method"),
	}

	createResult, err := client.Assistant(d.Get("assistant_sid").(string)).Webhooks.Create(createInput)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to create autopilot webhook: %s", err.Error())
	}

	d.SetId(createResult.Sid)
	return resourceAutopilotWebhookRead(d, meta)
}

func resourceAutopilotWebhookRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*common.TwilioClient).Autopilot

	getResponse, err := client.Assistant(d.Get("assistant_sid").(string)).Webhook(d.Id()).Fetch()
	if err != nil {
		if utils.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Failed to read autopilot webhook: %s", err.Error())
	}

	d.Set("sid", getResponse.Sid)
	d.Set("account_sid", getResponse.AccountSid)
	d.Set("assistant_sid", getResponse.AssistantSid)
	d.Set("unique_name", getResponse.UniqueName)
	d.Set("webhook_url", getResponse.WebhookURL)
	d.Set("webhook_method", getResponse.WebhookMethod)
	d.Set("events", strings.Split(getResponse.Events, webhookEventsSeperator))
	d.Set("date_created", getResponse.DateCreated.Format(time.RFC3339))

	if getResponse.DateUpdated != nil {
		d.Set("date_updated", getResponse.DateUpdated.Format(time.RFC3339))
	}

	d.Set("url", getResponse.URL)
	return nil
}

func resourceAutopilotWebhookUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*common.TwilioClient).Autopilot

	updateInput := &webhook.UpdateWebhookInput{
		UniqueName:    utils.OptionalString(d, "unique_name"),
		WebhookURL:    utils.OptionalString(d, "webhook_url"),
		Events:        utils.OptionalSeperatedString(d, "events", webhookEventsSeperator),
		WebhookMethod: utils.OptionalString(d, "webhook_method"),
	}

	updateResp, err := client.Assistant(d.Get("assistant_sid").(string)).Webhook(d.Id()).Update(updateInput)
	if err != nil {
		return fmt.Errorf("Failed to update autopilot webhook: %s", err.Error())
	}

	d.SetId(updateResp.Sid)
	return resourceAutopilotWebhookRead(d, meta)
}

func resourceAutopilotWebhookDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*common.TwilioClient).Autopilot

	if err := client.Assistant(d.Get("assistant_sid").(string)).Webhook(d.Id()).Delete(); err != nil {
		return fmt.Errorf("Failed to delete autopilot webhook: %s", err.Error())
	}
	d.SetId("")
	return nil
}