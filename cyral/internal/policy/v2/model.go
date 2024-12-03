package policyv2

import (
	"context"
	"fmt"
	"time"

	methods "buf.build/gen/go/cyral/policy/grpc/go/policy/v1/policyv1grpc"
	msg "buf.build/gen/go/cyral/policy/protocolbuffers/go/policy/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ChangeInfo represents information about changes to the policy
type ChangeInfo struct {
	Actor     string `json:"actor,omitempty"`
	ActorType string `json:"actorType,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// changeInfoToMap converts ChangeInfo to a map
func changeInfoToMap(c *msg.ChangeInfo) map[string]interface{} {
	return map[string]interface{}{
		"actor":      c.GetActor(),
		"actor_type": c.GetActorType().String(),
		"timestamp":  c.GetTimestamp().AsTime().Format(time.RFC3339),
	}
}

// scopeToMap converts Scope to a list of maps
func scopeToMap(s *msg.Scope) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"repo_ids": s.GetRepoIds(),
		},
	}
}

// updateSchema writes the policy data to the schema
func updateSchema(p *msg.Policy, ptype msg.PolicyType, d *schema.ResourceData) error {
	if err := d.Set("id", p.GetId()); err != nil {
		return fmt.Errorf("error setting 'id' field: %w", err)
	}
	if err := d.Set("name", p.GetName()); err != nil {
		return fmt.Errorf("error setting 'name' field: %w", err)
	}
	if err := d.Set("description", p.GetDescription()); err != nil {
		return fmt.Errorf("error setting 'description' field: %w", err)
	}
	if err := d.Set("enabled", p.GetEnabled()); err != nil {
		return fmt.Errorf("error setting 'enabled' field: %w", err)
	}
	if err := d.Set("tags", p.GetTags()); err != nil {
		return fmt.Errorf("error setting 'tags' field: %w", err)
	}
	if err := d.Set("valid_from", timestampFromProtobuf(p.GetValidFrom())); err != nil {
		return fmt.Errorf("error setting 'valid_from' field: %w", err)
	}
	if err := d.Set("valid_until", timestampFromProtobuf(p.GetValidUntil())); err != nil {
		return fmt.Errorf("error setting 'valid_until' field: %w", err)
	}
	if err := d.Set("document", p.GetDocument()); err != nil {
		return fmt.Errorf("error setting 'document' field: %w", err)
	}
	// Use the changeInfoToMap method to set the last_updated and created fields
	if err := d.Set("last_updated", changeInfoToMap(p.GetLastUpdated())); err != nil {
		return fmt.Errorf("error setting 'last_updated' field: %w", err)
	}
	if err := d.Set("created", changeInfoToMap(p.GetCreated())); err != nil {
		return fmt.Errorf("error setting 'created' field: %w", err)
	}
	if err := d.Set("enforced", p.GetEnforced()); err != nil {
		return fmt.Errorf("error setting 'enforced' field: %w", err)
	}
	// policy types have aliases, so we don't want to set the policy type
	// except if the new value is not an alias for the old one.
	if msg.PolicyType_value[d.Get("type").(string)] != int32(ptype) {
		if err := d.Set("type", ptype.String()); err != nil {
			return fmt.Errorf("error setting 'type' field: %w", err)
		}
	}
	if p.GetScope() != nil {
		if err := d.Set("scope", scopeToMap(p.GetScope())); err != nil {
			return fmt.Errorf("error setting 'scope' field: %w", err)
		}
	}
	d.SetId(p.GetId())
	return nil
}

func timestampFromResourceData(key string, d *schema.ResourceData) (*timestamppb.Timestamp, error) {
	if v, ok := d.GetOk(key); ok {
		ts := v.(string)
		if ts == "" {
			return nil, nil
		}
		if t, err := time.Parse(time.RFC3339, ts); err != nil {
			return nil, fmt.Errorf("invalid valid_from value: %s", ts)
		} else {
			return timestamppb.New(t), nil
		}
	}
	return nil, nil
}

func timestampFromProtobuf(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}
	return ts.AsTime().Format(time.RFC3339)
}

func policyAndTypeFromSchema(d *schema.ResourceData) (*msg.Policy, msg.PolicyType, error) {
	ptypeString := d.Get("type").(string)
	ptype := msg.PolicyType(msg.PolicyType_value[ptypeString])
	if ptype == msg.PolicyType_POLICY_TYPE_UNSPECIFIED {
		return nil, msg.PolicyType_POLICY_TYPE_UNSPECIFIED, fmt.Errorf(
			"invalid policy type: %s", ptypeString,
		)
	}
	p := &msg.Policy{
		Id:          d.Get("id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Enabled:     d.Get("enabled").(bool),
		Tags:        utils.ConvertFromInterfaceList[string](d.Get("tags").([]interface{})),
		Document:    d.Get("document").(string),
		Enforced:    d.Get("enforced").(bool),
	}
	var err error
	if p.ValidFrom, err = timestampFromResourceData("valid_from", d); err != nil {
		return nil, msg.PolicyType_POLICY_TYPE_UNSPECIFIED, nil
	}
	if p.ValidUntil, err = timestampFromResourceData("valid_until", d); err != nil {
		return nil, msg.PolicyType_POLICY_TYPE_UNSPECIFIED, nil
	}

	if v, ok := d.GetOk("scope"); ok {
		p.Scope = scopeFromInterface(v.([]interface{}))
	}
	return p, msg.PolicyType(ptype), nil
}

// scopeFromInterface converts the map to a Scope struct
func scopeFromInterface(s []interface{}) *msg.Scope {
	if len(s) == 0 || s[0] == nil {
		return nil
	}
	m := s[0].(map[string]interface{})
	scope := msg.Scope{
		RepoIds: utils.ConvertFromInterfaceList[string](m["repo_ids"].([]interface{})),
	}
	return &scope
}

func createPolicy(ctx context.Context, cl *client.Client, rd *schema.ResourceData) error {
	p, ptype, err := policyAndTypeFromSchema(rd)
	if err != nil {
		return err
	}
	req := &msg.CreatePolicyRequest{
		Type:   ptype,
		Policy: p,
	}
	grpcClient := methods.NewPolicyServiceClient(cl.GRPCClient())
	resp, err := grpcClient.CreatePolicy(ctx, req)
	if err != nil {
		return err
	}
	rd.SetId(resp.GetId())
	return nil
}

func readPolicy(ctx context.Context, cl *client.Client, rd *schema.ResourceData) error {
	p, ptype, err := policyAndTypeFromSchema(rd)
	if err != nil {
		return err
	}
	req := &msg.ReadPolicyRequest{
		Id:   p.GetId(),
		Type: ptype,
	}
	grpcClient := methods.NewPolicyServiceClient(cl.GRPCClient())
	resp, err := grpcClient.ReadPolicy(ctx, req)
	if err != nil {
		return err
	}
	return updateSchema(resp.GetPolicy(), ptype, rd)
}

func updatePolicy(ctx context.Context, cl *client.Client, rd *schema.ResourceData) error {
	p, ptype, err := policyAndTypeFromSchema(rd)
	if err != nil {
		return err
	}
	req := &msg.UpdatePolicyRequest{
		Id:     p.GetId(),
		Type:   ptype,
		Policy: p,
	}
	grpcClient := methods.NewPolicyServiceClient(cl.GRPCClient())
	_, err = grpcClient.UpdatePolicy(ctx, req)
	return err
}

func deletePolicy(ctx context.Context, cl *client.Client, rd *schema.ResourceData) error {
	p, ptype, err := policyAndTypeFromSchema(rd)
	if err != nil {
		return err
	}
	req := &msg.DeletePolicyRequest{
		Id:   p.GetId(),
		Type: ptype,
	}
	grpcClient := methods.NewPolicyServiceClient(cl.GRPCClient())
	_, err = grpcClient.DeletePolicy(ctx, req)
	return err
}
