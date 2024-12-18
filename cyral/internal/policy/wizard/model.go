package wizard

import (
	"context"

	methods "buf.build/gen/go/cyral/policy/grpc/go/policy/v1/policyv1grpc"
	msg "buf.build/gen/go/cyral/policy/protocolbuffers/go/policy/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
)

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
	updateSchema(wizardList, rd)
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

func updateSchema(wizards []*msg.PolicyWizard, rd *schema.ResourceData) {
	wizardList := make([]any, 0, len(wizards))
	for _, wiz := range wizards {
		wizardList = append(wizardList, wizardToMap(wiz))
	}
	rd.Set("wizards", wizardList)
	rd.SetId("cyral-wizard-list")
}
