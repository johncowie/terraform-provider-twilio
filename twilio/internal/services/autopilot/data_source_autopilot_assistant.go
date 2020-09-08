package autopilot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/RJPearson94/terraform-provider-twilio/twilio/common"
	"github.com/RJPearson94/terraform-provider-twilio/twilio/utils"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
)

func dataSourceAutopilotAssistant() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAutopilotAssistantRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"sid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"account_sid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_model_build_sid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"friendly_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"unique_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"callback_events": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"callback_url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_queries": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"development_stage": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"needs_model_build": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"defaults": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"stylesheet": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_updated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAutopilotAssistantRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*common.TwilioClient).Autopilot
	ctx, cancel := context.WithTimeout(meta.(*common.TwilioClient).StopContext, d.Timeout(schema.TimeoutRead))
	defer cancel()

	sid := d.Get("sid").(string)
	getResponse, err := client.Assistant(sid).FetchWithContext(ctx)
	if err != nil {
		if utils.IsNotFoundError(err) {
			return fmt.Errorf("[ERROR] Assistant with sid (%s) was not found", sid)
		}
		return fmt.Errorf("[ERROR] Failed to read autopilot assistant: %s", err.Error())
	}

	d.SetId(getResponse.Sid)
	d.Set("sid", getResponse.Sid)
	d.Set("account_sid", getResponse.AccountSid)
	d.Set("latest_model_build_sid", getResponse.LatestModelBuildSid)
	d.Set("unique_name", getResponse.UniqueName)
	d.Set("friendly_name", getResponse.FriendlyName)

	if getResponse.CallbackEvents != nil {
		d.Set("callback_events", strings.Split(*getResponse.CallbackEvents, " "))
	}

	d.Set("callback_url", getResponse.CallbackURL)
	d.Set("log_queries", getResponse.LogQueries)
	d.Set("development_stage", getResponse.DevelopmentStage)
	d.Set("needs_model_build", getResponse.NeedsModelBuild)
	d.Set("date_created", getResponse.DateCreated.Format(time.RFC3339))

	if getResponse.DateUpdated != nil {
		d.Set("date_updated", getResponse.DateUpdated.Format(time.RFC3339))
	}

	d.Set("url", getResponse.URL)

	getDefaultsResponse, err := client.Assistant(d.Id()).Defaults().FetchWithContext(ctx)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to read autopilot assistant defaults: %s", err.Error())
	}
	defaultsJSON, err := structure.FlattenJsonToString(getDefaultsResponse.Data)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to flatten defaults json to string: %s", err.Error())
	}
	d.Set("defaults", defaultsJSON)

	getStyleSheetResponse, err := client.Assistant(d.Id()).StyleSheet().FetchWithContext(ctx)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to read autopilot assistant stylesheet: %s", err.Error())
	}
	styleSheetJSON, err := structure.FlattenJsonToString(getStyleSheetResponse.Data)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to flatten stylesheet json to string: %s", err.Error())
	}
	d.Set("stylesheet", styleSheetJSON)

	return nil
}
