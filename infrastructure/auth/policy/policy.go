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
	if effect == "allow" {
		_, err := pm.e.AddPolicy(user, resource, action)
		if err != nil {
			return fmt.Errorf("error adding policy: %v", err)
		}
	} else if effect == "deny" {
		// Handle deny case if required (in Casbin, the default behavior can be modeled with additional policies)
		// Alternatively, this can be managed by the "eft" in the model configuration if required
		_, err := pm.e.AddPolicy(user, resource, action, "deny")
		if err != nil {
			return fmt.Errorf("error adding deny policy: %v", err)
		}
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
