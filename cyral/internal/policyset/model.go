package policyset

import (
	"context"
	"fmt"
	"time"

	methods "buf.build/gen/go/cyral/policy/grpc/go/policy/v1/policyv1grpc"
	msg "buf.build/gen/go/cyral/policy/protocolbuffers/go/policy/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
)

// ToMap converts PolicySetPolicy to a map
func policySetPolicyToMap(p *msg.PolicySetPolicy) map[string]interface{} {
	return map[string]interface{}{
		"type": p.GetType().String(),
		"id":   p.GetId(),
	}
}

func policiesToMaps(policies []*msg.PolicySetPolicy) []map[string]interface{} {
	var result []map[string]interface{}
	for _, policy := range policies {
		result = append(result, policySetPolicyToMap(policy))
	}
	return result
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

// updatePolicySetSchema writes the policy set data to the schema
func updatePolicySetSchema(ps *msg.PolicySet, d *schema.ResourceData) error {
	if err := d.Set("id", ps.GetId()); err != nil {
		return fmt.Errorf("error setting 'id' field: %w", err)
	}
	if err := d.Set("wizard_id", ps.GetWizardId()); err != nil {
		return fmt.Errorf("error setting 'id' field: %w", err)
	}
	if err := d.Set("name", ps.GetName()); err != nil {
		return fmt.Errorf("error setting 'name' field: %w", err)
	}
	if err := d.Set("description", ps.GetDescription()); err != nil {
		return fmt.Errorf("error setting 'description' field: %w", err)
	}
	if err := d.Set("enabled", ps.GetEnabled()); err != nil {
		return fmt.Errorf("error setting 'enabled' field: %w", err)
	}
	if err := d.Set("tags", ps.GetTags()); err != nil {
		return fmt.Errorf("error setting 'tags' field: %w", err)
	}
	if err := d.Set("wizard_parameters", ps.GetWizardParameters()); err != nil {
		return fmt.Errorf("error setting 'document' field: %w", err)
	}

	if err := d.Set("policies", policiesToMaps(ps.GetPolicies())); err != nil {
		return fmt.Errorf("error setting 'policies' field: %w", err)
	}
	if ps.GetScope() != nil {
		if err := d.Set("scope", scopeToMap(ps.GetScope())); err != nil {
			return fmt.Errorf("error setting 'scope' field: %w", err)
		}
	}
	// Use the changeInfoToMap method to set the last_updated and created fields
	if err := d.Set("last_updated", changeInfoToMap(ps.GetLastUpdated())); err != nil {
		return fmt.Errorf("error setting 'last_updated' field: %w", err)
	}
	if err := d.Set("created", changeInfoToMap(ps.GetCreated())); err != nil {
		return fmt.Errorf("error setting 'created' field: %w", err)
	}
	d.SetId(ps.GetId())
	return nil
}

func policySetFromSchema(d *schema.ResourceData) *msg.PolicySet {
	p := &msg.PolicySet{
		Id:               d.Get("id").(string),
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Enabled:          d.Get("enabled").(bool),
		Tags:             utils.ConvertFromInterfaceList[string](d.Get("tags").([]interface{})),
		WizardId:         d.Get("wizard_id").(string),
		WizardParameters: d.Get("wizard_parameters").(string),
	}

	if v, ok := d.GetOk("scope"); ok {
		p.Scope = scopeFromInterface(v.([]interface{}))
	}
	return p
}

// scopeFromInterface converts the map to a Scope struct
func scopeFromInterface(s []interface{}) *msg.Scope {
	if len(s) == 0 || s[0] == nil {
		return nil
	}
	m := s[0].(map[string]interface{})
	return &msg.Scope{
		RepoIds: utils.ConvertFromInterfaceList[string](m["repo_ids"].([]interface{})),
	}
}

func createPolicySet(ctx context.Context, cl *client.Client, rd *schema.ResourceData) error {
	ps := policySetFromSchema(rd)
	req := &msg.CreatePolicySetRequest{
		PolicySet: ps,
	}
	grpcClient := methods.NewPolicyWizardServiceClient(cl.GRPCClient())
	resp, err := grpcClient.CreatePolicySet(ctx, req)
	if err != nil {
		return err
	}
	rd.SetId(resp.GetPolicySet().GetId())
	return nil
}

func readPolicySet(ctx context.Context, cl *client.Client, rd *schema.ResourceData) error {
	req := &msg.ReadPolicySetRequest{
		Id: rd.Get("id").(string),
	}
	grpcClient := methods.NewPolicyWizardServiceClient(cl.GRPCClient())
	resp, err := grpcClient.ReadPolicySet(ctx, req)
	if err != nil {
		return err
	}
	return updatePolicySetSchema(resp.GetPolicySet(), rd)
}

func updatePolicySet(ctx context.Context, cl *client.Client, rd *schema.ResourceData) error {
	ps := policySetFromSchema(rd)
	req := &msg.UpdatePolicySetRequest{
		Id:        ps.GetId(),
		PolicySet: ps,
	}
	grpcClient := methods.NewPolicyWizardServiceClient(cl.GRPCClient())
	resp, err := grpcClient.UpdatePolicySet(ctx, req)
	if err != nil {
		return err
	}
	return updatePolicySetSchema(resp.GetPolicySet(), rd)
}

func deletePolicySet(ctx context.Context, cl *client.Client, rd *schema.ResourceData) error {
	req := &msg.DeletePolicySetRequest{
		Id: rd.Get("id").(string),
	}
	grpcClient := methods.NewPolicyWizardServiceClient(cl.GRPCClient())
	_, err := grpcClient.DeletePolicySet(ctx, req)
	return err
}

func readPolicyWizards(ctx context.Context, cl *client.Client, rd *schema.ResourceData) error {
	var wizardList []*msg.PolicyWizard

	wizId := rd.Get("wizard_id").(string)
	grpcClient := methods.NewPolicyWizardServiceClient(cl.GRPCClient())
	if wizId != "" {
		req := &msg.ReadPolicyWizardRequest{
			Id: wizId,
		}
		resp, err := grpcClient.ReadPolicyWizard(ctx, req)
		if err != nil && status.Code(err) != codes.NotFound {
			return err
		}
		if status.Code(err) != codes.NotFound {
			wizardList = []*msg.PolicyWizard{resp.GetPolicyWizard()}
		}
	} else {
		req := &msg.ListPolicyWizardsRequest{}
		resp, err := grpcClient.ListPolicyWizards(ctx, req)
		if err != nil {
			return err
		}
		wizardList = resp.GetPolicyWizards()
	}
	updatePolicyWizardsSchema(wizardList, rd)
	return nil
}

func wizardToMap(wiz *msg.PolicyWizard) map[string]any {
	return map[string]any{
		"id":               wiz.GetId(),
		"name":             wiz.GetName(),
		"description":      wiz.GetDescription(),
		"parameter_schema": wiz.GetParameterSchema(),
		"tags": func() []any {
			tags := make([]any, 0, len(wiz.GetTags()))
			for _, t := range wiz.GetTags() {
				tags = append(tags, t)
			}
			return tags
		}(),
	}
}

func updatePolicyWizardsSchema(wizards []*msg.PolicyWizard, rd *schema.ResourceData) {
	wizardList := make([]any, 0, len(wizards))
	for _, wiz := range wizards {
		wizardList = append(wizardList, wizardToMap(wiz))
	}
	rd.Set("wizards", wizardList)
	rd.SetId("cyral-wizard-list")
}
