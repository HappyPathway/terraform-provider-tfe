package tfe

import (
	"fmt"
	"log"

	tfe "github.com/HappyPathway/go-tfe"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceTFEVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceTFEVariableCreate,
		Read:   resourceTFEVariableRead,
		Update: resourceTFEVariableUpdate,
		Delete: resourceTFEVariableDelete,

		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"value": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"category": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						string(tfe.CategoryEnv),
						string(tfe.CategoryTerraform),
					},
					false,
				),
			},

			"hcl": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"sensitive": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"workspace_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceTFEVariableCreate(d *schema.ResourceData, meta interface{}) error {
	tfeClient := meta.(*tfe.Client)

	// Get key, category, workspace and organization.
	key := d.Get("key").(string)
	category := d.Get("category").(string)
	workspace, organization := unpackWorkspaceID(d.Get("workspace_id").(string))

	// Get the workspace.
	ws, err := tfeClient.Workspaces.Read(ctx, organization, workspace)
	if err != nil {
		return fmt.Errorf(
			"Error retrieving workspace %s from organization %s: %v", workspace, organization, err)
	}

	// Create a new options struct.
	options := tfe.VariableCreateOptions{
		Key:       tfe.String(key),
		Value:     tfe.String(d.Get("value").(string)),
		Category:  tfe.Category(tfe.CategoryType(category)),
		HCL:       tfe.Bool(d.Get("hcl").(bool)),
		Sensitive: tfe.Bool(d.Get("sensitive").(bool)),
		Workspace: ws,
	}

	log.Printf("[DEBUG] Create %s variable: %s", category, key)
	variable, err := tfeClient.Variables.Create(ctx, options)
	if err != nil {
		return fmt.Errorf("Error creating %s variable %s: %v", category, key, err)
	}

	d.SetId(variable.ID)

	return resourceTFEVariableRead(d, meta)
}

func resourceTFEVariableRead(d *schema.ResourceData, meta interface{}) error {
	tfeClient := meta.(*tfe.Client)

	// Get workspace and organization.
	workspace, organization := unpackWorkspaceID(d.Get("workspace_id").(string))

	// Create a new options struct.
	options := tfe.VariableListOptions{
		Organization: tfe.String(organization),
		Workspace:    tfe.String(workspace),
	}

	log.Printf("[DEBUG] List variables of workspace: %s", workspace)
	variables, err := tfeClient.Variables.List(ctx, options)
	if err != nil {
		if err == tfe.ErrResourceNotFound {
			log.Printf("[DEBUG] Variable %s does no longer exist", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error listing variables of workspace %s: %v", d.Id(), err)
	}

	var variable *tfe.Variable
	for _, v := range variables {
		if v.ID == d.Id() {
			variable = v
			break
		}
	}

	if variable == nil {
		log.Printf("[DEBUG] Variable %s does no longer exist", d.Id())
		d.SetId("")
		return nil
	}

	// Update config.
	d.Set("key", variable.Key)
	d.Set("category", string(variable.Category))
	d.Set("hcl", variable.HCL)
	d.Set("sensitive", variable.Sensitive)

	// Only set the value if its not sensitive, as otherwise it will be empty.
	if !variable.Sensitive {
		d.Set("value", variable.Value)
	}

	return nil
}

func resourceTFEVariableUpdate(d *schema.ResourceData, meta interface{}) error {
	tfeClient := meta.(*tfe.Client)

	// Create a new options struct.
	options := tfe.VariableUpdateOptions{
		Key:       tfe.String(d.Get("key").(string)),
		Value:     tfe.String(d.Get("value").(string)),
		HCL:       tfe.Bool(d.Get("hcl").(bool)),
		Sensitive: tfe.Bool(d.Get("sensitive").(bool)),
	}

	log.Printf("[DEBUG] Update variable: %s", d.Id())
	_, err := tfeClient.Variables.Update(ctx, d.Id(), options)
	if err != nil {
		return fmt.Errorf("Error updating variable %s: %v", d.Id(), err)
	}

	return resourceTFEVariableRead(d, meta)
}

func resourceTFEVariableDelete(d *schema.ResourceData, meta interface{}) error {
	tfeClient := meta.(*tfe.Client)

	log.Printf("[DEBUG] Delete variable: %s", d.Id())
	err := tfeClient.Variables.Delete(ctx, d.Id())
	if err != nil {
		if err == tfe.ErrResourceNotFound {
			return nil
		}
		return fmt.Errorf("Error deleting variable%s: %v", d.Id(), err)
	}

	return nil
}
