package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

// AzureAuthenticator handles Azure authentication
type AzureAuthenticator struct {
	credential azcore.TokenCredential
	tenantID   string
}

// NewAzureAuthenticator creates a new authenticator
// It tries Azure CLI credential first, then falls back to environment-based auth
func NewAzureAuthenticator() (*AzureAuthenticator, error) {
	var cred azcore.TokenCredential
	var err error

	// Try Service Principal from environment variables first
	clientID := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	tenantID := os.Getenv("AZURE_TENANT_ID")

	if clientID != "" && clientSecret != "" && tenantID != "" {
		// Use Service Principal authentication
		cred, err = azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create service principal credential: %w", err)
		}
	} else {
		// Fall back to Azure CLI credential
		cred, err = azidentity.NewAzureCLICredential(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create Azure CLI credential: %w (hint: run 'az login' or set AZURE_CLIENT_ID, AZURE_CLIENT_SECRET, AZURE_TENANT_ID)", err)
		}
	}

	return &AzureAuthenticator{
		credential: cred,
		tenantID:   tenantID,
	}, nil
}

// GetCredential returns the Azure credential
func (a *AzureAuthenticator) GetCredential() azcore.TokenCredential {
	return a.credential
}

// GetTenantID returns the tenant ID if available
func (a *AzureAuthenticator) GetTenantID() string {
	return a.tenantID
}

// ValidateSubscriptionAccess validates that we can access the subscription
func (a *AzureAuthenticator) ValidateSubscriptionAccess(subscriptionID string) error {
	ctx := context.Background()

	clientFactory, err := armsubscriptions.NewClientFactory(a.credential, nil)
	if err != nil {
		return fmt.Errorf("failed to create subscription client factory: %w", err)
	}

	client := clientFactory.NewClient()

	// Try to get the subscription details
	_, err = client.Get(ctx, subscriptionID, nil)
	if err != nil {
		return fmt.Errorf("failed to access subscription %s: %w (check permissions)", subscriptionID, err)
	}

	return nil
}

// ListSubscriptions lists all accessible subscriptions
func (a *AzureAuthenticator) ListSubscriptions() ([]Subscription, error) {
	ctx := context.Background()

	clientFactory, err := armsubscriptions.NewClientFactory(a.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription client factory: %w", err)
	}

	client := clientFactory.NewClient()
	pager := client.NewListPager(nil)

	var subscriptions []Subscription
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list subscriptions: %w", err)
		}

		for _, sub := range page.Value {
			if sub.SubscriptionID == nil || sub.DisplayName == nil {
				continue
			}

			subscription := Subscription{
				ID:   *sub.SubscriptionID,
				Name: *sub.DisplayName,
			}

			if sub.State != nil {
				subscription.State = string(*sub.State)
			}

			if sub.TenantID != nil {
				subscription.TenantID = *sub.TenantID
			}

			subscriptions = append(subscriptions, subscription)
		}
	}

	return subscriptions, nil
}

// Subscription represents a basic subscription
type Subscription struct {
	ID       string
	Name     string
	State    string
	TenantID string
}
