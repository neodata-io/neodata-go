package policy

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/neodata-io/neodata-go/config"
)

type PolicyManager struct {
	e *casbin.Enforcer
}

func NewPolicyManager(cfg *config.AppConfig) (*PolicyManager, error) {
	// TODO: implement caching or singleton to prevent initiated multiple times
	enforcer, err := InitializeCasbin(cfg)
	if err != nil {
		return nil, fmt.Errorf("error adding policy: %v", err)
	}
	return &PolicyManager{
		e: enforcer,
	}, nil
}

// AddPolicyForUser adds a specific policy for a user with an effect (e.g., allow or deny)
func (pm *PolicyManager) AddPolicyForUser(user string, resource string, action string, effect string) error {
	// Add policy with the subject (user), object (resource), action, and effect (allow or deny)
	_, err := pm.e.AddPolicy(user, resource, action, effect)
	if err != nil {
		return fmt.Errorf("error adding policy: %v", err)
	}
	return nil
}

func (pm *PolicyManager) AddPoliciesForUser(userID string, policies [][]string) error {
	for _, policy := range policies {
		if len(policy) != 3 {
			return fmt.Errorf("invalid policy format: %v", policy)
		}
		err := pm.AddPolicyForUser(userID, policy[0], policy[1], policy[2])
		if err != nil {
			return err
		}
	}
	return nil
}

// GetFilteredPolicy retrieves policies for a specific user.
// index specifies the field to filter on; userID is the value to match.
func (pm *PolicyManager) GetFilteredPolicy(index int, userID string) ([][]string, error) {
	// Use Casbin's GetFilteredPolicy method to get policies for the given user.
	policies, err := pm.e.GetFilteredPolicy(index, userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving permission for a specific user: %v", err)
	}

	// Check if any policies were returned.
	if len(policies) == 0 {
		return nil, fmt.Errorf("no policies found for user: %s", userID)
	}

	return policies, err
}

// HasPolicyForUser checks if a specific policy exists for a user
func (pm *PolicyManager) HasPolicyForUser(userID string, resource string, action string, effect string) (bool, error) {
	exists, err := pm.e.HasPolicy(userID, resource, action, effect)
	if err != nil {
		return false, fmt.Errorf("error checking policy permission for user %s: %v", userID, err)
	}
	return exists, nil
}

// CanUserLogin checks if the user is allowed to execute the "login" action.
func (pm *PolicyManager) CanUserLogin(userID string) (bool, error) {
	// Use Casbin's Enforce method to check if the user can execute the login action.
	allowed, err := pm.e.Enforce(userID, "login", "execute")
	if err != nil {
		return false, fmt.Errorf("error checking login permission for user %s: %v", userID, err)
	}

	if !allowed {
		return false, fmt.Errorf("user %s is not allowed to log in", userID)
	}

	return true, nil
}

// RemovePolicyForUser removes a specific policy for a user (subject, object, action, effect)
func (pm *PolicyManager) RemovePolicyForUser(user string, resource string, action string, effect string) error {
	// Remove a specific policy that matches all four fields
	removed, err := pm.e.RemovePolicy(user, resource, action, effect)
	if err != nil {
		return fmt.Errorf("error removing policy: %v", err)
	}

	if !removed {
		return fmt.Errorf("policy not found for user: %s, resource: %s, action: %s", user, resource, action)
	}

	return nil
}

// RemoveAllPoliciesForUser removes all policies for a specific user
func (pm *PolicyManager) RemoveAllPoliciesForUser(user string) error {
	// Remove all policies where the subject (user) matches
	removed, err := pm.e.RemoveFilteredPolicy(0, user)
	if err != nil {
		return fmt.Errorf("error removing policies for user: %v", err)
	}

	if !removed {
		return fmt.Errorf("no policies found for user: %s", user)
	}

	return nil
}

// AddMultiplePolicies adds multiple policies for multiple users in one call
func (pm *PolicyManager) AddMultiplePolicies(policies [][]string) error {
	ok, err := pm.e.AddPolicies(policies)
	if err != nil || !ok {
		return fmt.Errorf("error adding multiple policies: %v", err)
	}
	return nil
}

// RemoveMultiplePolicies removes multiple policies in one call
func (pm *PolicyManager) RemoveMultiplePolicies(policies [][]string) error {
	ok, err := pm.e.RemovePolicies(policies)
	if err != nil || !ok {
		return fmt.Errorf("error removing multiple policies: %v", err)
	}
	return nil
}

// CanUserPerformAction checks if a user is allowed to perform a specific action on a resource
func (pm *PolicyManager) CanUserPerformAction(user string, resource string, action string) (bool, error) {
	allowed, err := pm.e.Enforce(user, resource, action)
	if err != nil {
		return false, fmt.Errorf("error enforcing policy: %v", err)
	}
	return allowed, nil
}

// ResetPolicies clears all policies in the system (use with caution)
func (pm *PolicyManager) ResetPolicies() {
	pm.e.ClearPolicy()
}

func (pm *PolicyManager) ReloadPolicies() error {
	if err := pm.e.LoadPolicy(); err != nil {
		return fmt.Errorf("failed to reload policies: %w", err)
	}
	return nil
}
