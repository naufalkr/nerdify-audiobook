package utils

import "time"

// Subscription plans
const (
	PlanBasic      = "Basic"
	PlanPremium    = "Premium"
	PlanEnterprise = "Enterprise"
)

// SubscriptionPlan defines the properties for each subscription plan
type SubscriptionPlan struct {
	Name        string
	MaxUsers    int
	Description string
	Features    []string
}

// GetSubscriptionPlans returns all available subscription plans
func GetSubscriptionPlans() map[string]SubscriptionPlan {
	return map[string]SubscriptionPlan{
		PlanBasic: {
			Name:        PlanBasic,
			MaxUsers:    50,
			Description: "Basic plan for small teams",
			Features: []string{
				"Up to 50 users",
				"Basic features",
				"Standard support",
			},
		},
		PlanPremium: {
			Name:        PlanPremium,
			MaxUsers:    100,
			Description: "Premium plan for medium-sized organizations",
			Features: []string{
				"Up to 100 users",
				"All Basic features",
				"Advanced analytics",
				"Premium support",
			},
		},
		PlanEnterprise: {
			Name:        PlanEnterprise,
			MaxUsers:    500,
			Description: "Enterprise plan for large organizations",
			Features: []string{
				"Up to 500 users",
				"All Premium features",
				"Custom integrations",
				"Dedicated support",
				"SSO authentication",
			},
		},
	}
}

// GetMaxUsersByPlan returns the maximum number of users allowed for a given plan
func GetMaxUsersByPlan(planName string) int {
	plans := GetSubscriptionPlans()
	if plan, exists := plans[planName]; exists {
		return plan.MaxUsers
	}
	// Default to Basic if plan not found
	return plans[PlanBasic].MaxUsers
}

// GetValidPlans returns a list of valid subscription plan names
func GetValidPlans() []string {
	return []string{PlanBasic, PlanPremium, PlanEnterprise}
}

// CalculateOneYearSubscription calculates the start and end dates for a one-year subscription
func CalculateOneYearSubscription(customStartDate string) (startDate, endDate time.Time) {
	// Default to current time if no custom start date
	startDate = time.Now()

	// Try to parse custom start date if provided
	if customStartDate != "" {
		if parsed, err := time.Parse("2006-01-02", customStartDate); err == nil {
			startDate = parsed
		}
	}

	// End date is one year from start date
	endDate = startDate.AddDate(1, 0, 0)

	return startDate, endDate
}
