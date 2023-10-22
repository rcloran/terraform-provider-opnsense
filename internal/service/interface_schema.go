package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/browningluke/opnsense-go/pkg/diagnostics"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-opnsense/internal/tools"
)

type InterfaceDataSourceModel struct {
	Device     types.String `tfsdk:"device"`
	Media      types.String `tfsdk:"media"`
	MediaRaw   types.String `tfsdk:"media_raw"`
	MacAddr    types.String `tfsdk:"macaddr"`
	IsPhysical types.Bool   `tfsdk:"is_physical"`
	Mtu        types.Int64  `tfsdk:"mtu"`
	Status     types.String `tfsdk:"status"`

	Flags          types.Set `tfsdk:"flags"`
	Capabilities   types.Set `tfsdk:"capabilities"`
	Options        types.Set `tfsdk:"options"`
	SupportedMedia types.Set `tfsdk:"supported_media"`
	Groups         types.Set `tfsdk:"groups"`

	Ipv4 types.List `tfsdk:"ipv4"`
	Ipv6 types.List `tfsdk:"ipv6"`
}

type Ipv4Model struct {
	Ipaddr     types.String `tfsdk:"ipaddr"`
	SubnetBits types.Int64  `tfsdk:"subnetbits"`
	Tunnel     types.Bool   `tfsdk:"tunnel"`
}

type Ipv6Model struct {
	Ipaddr     types.String `tfsdk:"ipaddr"`
	SubnetBits types.Int64  `tfsdk:"subnetbits"`
	Tunnel     types.Bool   `tfsdk:"tunnel"`
	Autoconf   types.Bool   `tfsdk:"autoconf"`
	Deprecated types.Bool   `tfsdk:"deprecated"`
	LinkLocal  types.Bool   `tfsdk:"link_local"`
	Tentative  types.Bool   `tfsdk:"tentative"`
}

var ipv4AttrTypes = map[string]attr.Type{
	"ipaddr":     types.StringType,
	"subnetbits": types.Int64Type,
	"tunnel":     types.BoolType,
}

var ipv6AttrTypes = map[string]attr.Type{
	"ipaddr":     types.StringType,
	"subnetbits": types.Int64Type,
	"tunnel":     types.BoolType,
	"autoconf":   types.BoolType,
	"deprecated": types.BoolType,
	"link_local": types.BoolType,
	"tentative":  types.BoolType,
}

func InterfaceDataSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Interfaces describe network interfaces.",

		Attributes: map[string]schema.Attribute{
			"device": schema.StringAttribute{
				MarkdownDescription: "Device name.",
				Required:            true,
			},
			"media": schema.StringAttribute{
				MarkdownDescription: "Media type.",
				Computed:            true,
			},
			"media_raw":   schema.StringAttribute{Computed: true},
			"macaddr":     schema.StringAttribute{Computed: true},
			"is_physical": schema.BoolAttribute{Computed: true},
			"mtu":         schema.Int64Attribute{Computed: true},
			"status":      schema.StringAttribute{Computed: true},

			"flags":           schema.SetAttribute{Computed: true, ElementType: types.StringType},
			"capabilities":    schema.SetAttribute{Computed: true, ElementType: types.StringType},
			"options":         schema.SetAttribute{Computed: true, ElementType: types.StringType},
			"supported_media": schema.SetAttribute{Computed: true, ElementType: types.StringType},
			"groups":          schema.SetAttribute{Computed: true, ElementType: types.StringType},
			"ipv4": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"ipaddr":     schema.StringAttribute{Computed: true},
					"subnetbits": schema.Int64Attribute{Computed: true},
					"tunnel":     schema.BoolAttribute{Computed: true},
				},
			}},
			"ipv6": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"ipaddr":     schema.StringAttribute{Computed: true},
					"subnetbits": schema.Int64Attribute{Computed: true},
					"tunnel":     schema.BoolAttribute{Computed: true},
					"autoconf":   schema.BoolAttribute{Computed: true},
					"deprecated": schema.BoolAttribute{Computed: true},
					"link_local": schema.BoolAttribute{Computed: true},
					"tentative":  schema.BoolAttribute{Computed: true},
				},
			}},
		},
	}
}

func convertInterfaceConfigStructToSchema(ctx context.Context, data *map[string]diagnostics.InterfaceConfig, device string) (*InterfaceDataSourceModel, error) {
	var d, ok = (*data)[device]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Err: cannot find interface %s", device))
	}
	model := &InterfaceDataSourceModel{
		Device:         types.StringValue(d.Device),
		Media:          types.StringValue(d.Media),
		MediaRaw:       types.StringValue(d.MediaRaw),
		MacAddr:        types.StringValue(d.MacAddr),
		IsPhysical:     types.BoolValue(d.IsPhysical),
		Mtu:            tools.StringToInt64Null(d.Mtu),
		Status:         types.StringValue(d.Status),
		Flags:          tools.StringSliceToSet(d.Flags),
		Capabilities:   tools.StringSliceToSet(d.Capabilities),
		Options:        tools.StringSliceToSet(d.Options),
		SupportedMedia: tools.StringSliceToSet(d.SupportedMedia),
		Groups:         tools.StringSliceToSet(d.Groups),
	}

	var ipv4s []Ipv4Model
	for _, elem := range d.Ipv4 {
		ipv4 := Ipv4Model{
			Ipaddr:     types.StringValue(elem.IpAddr),
			SubnetBits: types.Int64Value(elem.SubnetBits),
			Tunnel:     types.BoolValue(elem.Tunnel),
		}
		ipv4s = append(ipv4s, ipv4)
	}
	var ipv6s []Ipv6Model
	for _, elem := range d.Ipv6 {
		ipv6 := Ipv6Model{
			Ipaddr:     types.StringValue(elem.IpAddr),
			SubnetBits: types.Int64Value(elem.SubnetBits),
			Tunnel:     types.BoolValue(elem.Tunnel),
			Autoconf:   types.BoolValue(elem.Autoconf),
			Deprecated: types.BoolValue(elem.Deprecated),
			LinkLocal:  types.BoolValue(elem.LinkLocal),
			Tentative:  types.BoolValue(elem.Tentative),
		}
		ipv6s = append(ipv6s, ipv6)
	}
	model.Ipv4, _ = types.ListValueFrom(ctx, types.ObjectType{}.WithAttributeTypes(ipv4AttrTypes), ipv4s)
	model.Ipv6, _ = types.ListValueFrom(ctx, types.ObjectType{}.WithAttributeTypes(ipv6AttrTypes), ipv6s)
	return model, nil
}
