package pixela

import (
	"context"
	"fmt"
	"time"

	pixela "github.com/ebc-2in2crc/pixela4go"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGraph() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGraphCreate,
		ReadContext:   resourceGraphRead,
		UpdateContext: resourceGraphUpdate,
		DeleteContext: resourceGraphDelete,
		/*
		   {
		     "id": "test-graph",
		     "name": "graph-name",
		     "unit": "commit",
		     "type": "int",
		     "color": "shibafu",
		     "timezone": "Asia/Tokyo",
		     "purgeCacheURLs": [
		       "https://camo.githubusercontent.com/xxx/xxxx"
		     ],
		     "selfSufficient": "increment",
		     "isSecret": false,
		     "publishOptionalData": true
		   }
		*/
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"graph_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"unit": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"color": {
				Type:     schema.TypeString,
				Required: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "UTC",
			},
			//"purge_cache_urls": {
			//	Type:     schema.TypeList,
			//	Optional: true,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeString,
			//	},
			//},
			"self_sufficient": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "none",
			},
			"is_secret": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"publish_optional_data": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func flattenGraph(graph pixela.GraphDefinition) interface{} {
	g := make(map[string]interface{})
	g["graph_id"] = graph.ID
	g["name"] = graph.Name
	g["unit"] = graph.Unit
	g["type"] = graph.Type
	g["color"] = graph.Color
	g["timezone"] = graph.TimeZone
	// below fields is used camel case in json tags.
	//g["purge_cache_urls"] = graph.PurgeCacheURLs
	g["self_sufficient"] = graph.SelfSufficient
	g["is_secret"] = graph.IsSecret
	g["publish_optional_data"] = graph.PublishOptionalData

	return g
}

func resourceGraphCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := m.(*pixela.Client)

	id := d.Get("graph_id").(string)
	name := d.Get("name").(string)
	unit := d.Get("unit").(string)
	gtype := d.Get("type").(string)
	// TODO: validate color type
	color := d.Get("color").(string)
	timezone := d.Get("timezone").(string)
	selfSufficient := d.Get("self_sufficient").(string)
	is := d.Get("is_secret").(bool)
	pod := d.Get("publish_optional_data").(bool)

	result, err := client.Graph().Create(&pixela.GraphCreateInput{
		ID:                  String(id),
		Name:                String(name),
		Unit:                String(unit),
		Type:                String(gtype),
		Color:               String(color),
		TimeZone:            String(timezone),
		SelfSufficient:      String(selfSufficient),
		IsSecret:            Bool(is),
		PublishOptionalData: Bool(pod),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if !result.IsSuccess {
		return diag.FromErr(fmt.Errorf(result.Message))
	}
	d.SetId(id)
	return diags
}

func resourceGraphRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := m.(*pixela.Client)

	// FIXME: pixela api cannot get a graph.
	result, err := client.Graph().GetAll()
	if err != nil {
		return diag.FromErr(err)
	}

	var g pixela.GraphDefinition
	var found bool
	gid := d.Id()
	for _, graph := range result.Graphs {
		if graph.ID == gid {
			g = graph
			found = true
			break
		}
	}
	if !found {
		return diag.FromErr(fmt.Errorf("cannot find graph %q", gid))
	}
	if err := d.Set("graph_id", g.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", g.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("type", g.Type); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("unit", g.Unit); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("color", g.Color); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("timezone", g.TimeZone); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("self_sufficient", g.SelfSufficient); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_secret", g.IsSecret); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("publish_optional_data", g.PublishOptionalData); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceGraphUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*pixela.Client)

	if d.HasChange("graph_id") {
		return diag.Errorf("cannot support update graph_id")
	}
	if d.HasChange("type") {
		return diag.Errorf("cannot support update type")
	}

	if d.HasChange("name") || d.HasChange("unit") ||
		d.HasChange("color") || d.HasChange("timezone") ||
		d.HasChange("self_sufficient") || d.HasChange("is_secret") ||
		d.HasChange("publish_optional_data") {
		id := d.Get("graph_id").(string)
		name := d.Get("name").(string)
		unit := d.Get("unit").(string)
		// TODO: validate color type
		color := d.Get("color").(string)
		timezone := d.Get("timezone").(string)
		selfSufficient := d.Get("self_sufficient").(string)
		is := d.Get("is_secret").(bool)
		pod := d.Get("publish_optional_data").(bool)

		_, err := c.Graph().Update(&pixela.GraphUpdateInput{
			ID:       String(id),
			Name:     String(name),
			Unit:     String(unit),
			Color:    String(color),
			TimeZone: String(timezone),
			// PurgeCacheURLs:      nil,
			SelfSufficient:      String(selfSufficient),
			IsSecret:            Bool(is),
			PublishOptionalData: Bool(pod),
		})
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceGraphRead(ctx, d, m)
}

func resourceGraphDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*pixela.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	graphID := d.Id()

	r, err := c.Graph().Delete(&pixela.GraphDeleteInput{ID: String(graphID)})
	if err != nil {
		return diag.FromErr(err)
	}
	if !r.IsSuccess {
		return diag.Errorf("destroy failed %q", r.Message)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}