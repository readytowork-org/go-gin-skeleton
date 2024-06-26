package services

import (
	"boilerplate-api/infrastructure"
	"context"
	"fmt"

	"cloud.google.com/go/billing/budgets/apiv1/budgetspb"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/googleapis/type/money"
)

// GCPBillingService -> handles the gcp billing related functions
type GCPBillingService struct {
	logger     infrastructure.Logger
	gcpBilling infrastructure.GCPBilling
	env        infrastructure.Env
}

// NewGCPBillingService -> initilization for the GCPBilling struct
func NewGCPBillingService(
	logger infrastructure.Logger,
	gcpBilling infrastructure.GCPBilling,
	env infrastructure.Env,
) GCPBillingService {
	return GCPBillingService{
		logger:     logger,
		gcpBilling: gcpBilling,
		env:        env,
	}
}

// Get Billing info for certain date
func (s GCPBillingService) GetBillingInfo(
	ctx context.Context,
) (*cloudbilling.ProjectBillingInfo, error) {
	projectName := fmt.Sprintf("projects/%s", s.env.ProjectName)
	billingInfo, err := s.gcpBilling.BillingClient.Projects.GetBillingInfo(projectName).Do()

	return billingInfo, err
}

// Get Billing info for certain date
func (s GCPBillingService) GetExistingBudgetList(
	ctx context.Context,
) (*budgetspb.Budget, error) {
	var budgetList []*budgetspb.Budget
	var err error
	parentId := fmt.Sprintf("billingAccounts/%s", s.env.BillingAccountId)
	req := budgetspb.ListBudgetsRequest{
		Parent: parentId,
	}

	budgetsIter := s.gcpBilling.BudgetClient.ListBudgets(ctx, &req)
	for {
		budget, err := budgetsIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			s.logger.Zap.Errorf("failed to retrieve budget: %v", err)
		}
		budgetList = append(budgetList, budget)
	}
	if len(budgetList) > 0 {
		return budgetList[0], err
	}

	return nil, err
}

func (s GCPBillingService) GetBudgetCreateUpdateRequest() *budgetspb.Budget {
	projectName := fmt.Sprintf("projects/%s", s.env.ProjectName)

	budget := &budgetspb.Budget{
		DisplayName: "Project Budget",
		Name:        s.env.BudgetDisplayName,
		BudgetFilter: &budgetspb.Filter{
			CreditTypesTreatment: budgetspb.Filter_INCLUDE_ALL_CREDITS,
			Projects:             []string{projectName},
		},
		Amount: &budgetspb.BudgetAmount{
			BudgetAmount: &budgetspb.BudgetAmount_SpecifiedAmount{
				SpecifiedAmount: &money.Money{
					Units:        s.env.BudgetAmount,
					Nanos:        0,
					CurrencyCode: "JPY",
				},
			},
		},
		ThresholdRules: []*budgetspb.ThresholdRule{
			{
				ThresholdPercent: 0.25,
				SpendBasis:       budgetspb.ThresholdRule_CURRENT_SPEND,
			},
			{
				ThresholdPercent: 0.50,
				SpendBasis:       budgetspb.ThresholdRule_CURRENT_SPEND,
			},
			{
				ThresholdPercent: 0.75,
				SpendBasis:       budgetspb.ThresholdRule_CURRENT_SPEND,
			},
			{
				ThresholdPercent: 1,
				SpendBasis:       budgetspb.ThresholdRule_CURRENT_SPEND,
			},
		},
	}
	return budget
}

func (s GCPBillingService) CreateBudget(ctx context.Context) (*budgetspb.Budget, error) {
	parentId := fmt.Sprintf("billingAccounts/%s", s.env.BillingAccountId)
	budget := s.GetBudgetCreateUpdateRequest()
	createRequest := budgetspb.CreateBudgetRequest{
		Parent: parentId,
		Budget: budget,
	}

	billingInfo, err := s.gcpBilling.BudgetClient.CreateBudget(ctx, &createRequest)
	if err != nil {
		s.logger.Zap.Errorf("failed to create budget: %v", err)
	}

	return billingInfo, err
}

func (s GCPBillingService) CreateOrUpdateBudget(ctx context.Context) (*budgetspb.Budget, error) {
	budget, _ := s.GetExistingBudgetList(ctx)
	if budget == nil {
		return s.CreateBudget(ctx)
	} else {
		return s.EditBudget(ctx, budget)
	}

}

func (s GCPBillingService) EditBudget(ctx context.Context, budget *budgetspb.Budget) (*budgetspb.Budget, error) {
	// Modify the budget configuration here
	envBudget := s.GetBudgetCreateUpdateRequest()

	budget.Amount.BudgetAmount = envBudget.Amount.BudgetAmount
	budget.ThresholdRules = envBudget.ThresholdRules

	editRequest := budgetspb.UpdateBudgetRequest{
		Budget: budget,
	}

	billingInfo, err := s.gcpBilling.BudgetClient.UpdateBudget(ctx, &editRequest)
	if err != nil {
		s.logger.Zap.Errorf("failed to retrieve budget: %v", err)
	}

	return billingInfo, err
}
