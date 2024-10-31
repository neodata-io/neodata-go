package interfaces

// PolicyManager defines an interface for managing policies, allowing for multiple implementations.
type PolicyManager interface {
	// Policy Management
	AddPolicyForUser(user string, resource string, action string, effect string) error
	AddPoliciesForUser(userID string, policies [][]string) error
	AddMultiplePolicies(policies [][]string) error
	RemovePolicyForUser(user string, resource string, action string, effect string) error
	RemoveAllPoliciesForUser(user string) error
	RemoveMultiplePolicies(policies [][]string) error
	ResetPolicies()
	ReloadPolicies() error

	// Policy Checks
	GetFilteredPolicy(index int, userID string) ([][]string, error)
	HasPolicyForUser(userID string, resource string, action string, effect string) (bool, error)
	CanUserLogin(userID string) (bool, error)
	CanUserPerformAction(user string, resource string, action string) (bool, error)
}
