package tfe

import (
	"fmt"
	"log"

	tfe "github.com/HappyPathway/go-tfe"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTFERegistryModule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTFERegistryModuleCreate,
		Read:   resourceTFERegistryModuleRead,
		Update: resourceTFERegistryModuleUpdate,
		Delete: resourceTFERegistryModuleDelete,

		Schema: map[string]*schema.Schema{
			"organization": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"oauth_token": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"repo": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"module_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Computed: true,
			},
			"module_provider": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Computed: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Computed: true,
			},
		},
	}
}

func resourceTFERegistryModuleCreate(d *schema.ResourceData, meta interface{}) error {
	tfeClient := meta.(*tfe.Client)

	// Get the organization name.
	oauth_token := d.Get("oauth_token").(string)
	repo := d.Get("repo").(string)

	vcs_repo := tfe.RegistryModuleVCSRepo{
		IDENTIFIER:     repo,
		OAUTH_TOKEN_ID: oauth_token,
	}

	options := tfe.RegistryModuleCreateOptions{
		VCSRepo: &vcs_repo,
	}

	rm, err := tfeClient.RegistryModules.Create(ctx, options)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	d.SetId(rm.ID)

	return nil
}

func resourceTFERegistryModuleRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceTFERegistryModuleUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceTFERegistryModuleDelete(d *schema.ResourceData, meta interface{}) error {
	tfeClient := meta.(*tfe.Client)

	// Get the organization name.

	organization := d.Get("organization").(string)
	module_id := d.Get("module_id").(string)
	provider := d.Get("provider").(string)

	// Create a new options struct.
	log.Printf("[DEBUG] Delete Registry Module Organization: %s Module: %s Provider: %s ", organization, module_id, provider)

	vcs, err := tfeClient.RegistryModules.Delete(ctx, organization, module_id, provider)
	if err != nil {
		return fmt.Errorf("Error deleting the  module %s: %v", module_id, err)
	}
	d.SetId(vcs.ID)

	return nil
}
