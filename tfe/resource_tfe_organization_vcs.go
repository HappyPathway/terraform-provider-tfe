package tfe

import (
	"fmt"
	"log"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTFEOrganizationVCS() *schema.Resource {
	return &schema.Resource{
		Create: resourceTFEOrganizationVCSCreate,
		Read:   resourceTFEOrganizationVCSRead,
		Update: resourceTFEOrganizationVCSUpdate,
		Delete: resourceTFEOrganizationVCSDelete,

		Schema: map[string]*schema.Schema{
			"organization": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"github_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"github_secret": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"service_provider": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"http_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"api_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vcs_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Computed: true,
			},
			"callback_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Computed: true,
			},
		},
	}
}

func resourceTFEOrganizationVCSCreate(d *schema.ResourceData, meta interface{}) error {
	tfeClient := meta.(*tfe.Client)

	// Get the organization name.
	organization := d.Get("organization").(string)
	key := d.Get("github_key").(string)
	secret := d.Get("github_secret").(string)
	http_url := d.Get("http_url").(string)
	api_url := d.Get("api_url").(string)
	service_provider := d.Get("service_provider")

	// Create a new options struct.
	if service_provider == "github" {
		options := tfe.OAuthClientCreateOptions{
			ServiceProvider: tfe.ServiceProvider(tfe.ServiceProviderGithub),
			APIURL:          tfe.String(api_url),
			HTTPURL:         tfe.String(http_url),
			Key:             tfe.String(key),
			Secret:          tfe.String(secret),
		}
		log.Printf("[DEBUG] Create new VCS Connection for Org: %s", organization)
		vcs, err := tfeClient.OAuthClients.Create(ctx, organization, options)
		if err != nil {
			return fmt.Errorf("Error creating the new organization %s: %v", organization, err)
		}

		d.SetId(vcs.ID)
		d.Set("vcs_id", vcs.ID)
		d.Set("callback_url", vcs.CallbackURL)

	}
	return nil
}

func resourceTFEOrganizationVCSRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceTFEOrganizationVCSUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceTFEOrganizationVCSDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
