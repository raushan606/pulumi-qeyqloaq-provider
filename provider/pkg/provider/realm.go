package provider

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// Realm represents a Keycloak realm resource with merge strategy
// This provider only updates fields that are explicitly managed,
// preserving manual changes to other realm attributes in Keycloak UI
type Realm struct{}

// RealmArgs represents the input arguments for creating/updating a realm
// Only these fields will be managed by Pulumi - all others are preserved
type RealmArgs struct {
	// The name of the realm (required)
	Name string `pulumi:"name"`

	// Whether the realm is enabled (managed field)
	Enabled *bool `pulumi:"enabled,optional"`

	// Display name for the realm (managed field)
	DisplayName *string `pulumi:"displayName,optional"`

	// HTML display name (managed field)
	DisplayNameHtml *string `pulumi:"displayNameHtml,optional"`

	// Login theme (managed field)
	LoginTheme *string `pulumi:"loginTheme,optional"`

	// Account theme (managed field)
	AccountTheme *string `pulumi:"accountTheme,optional"`

	// Admin theme (managed field)
	AdminTheme *string `pulumi:"adminTheme,optional"`

	// Email theme (managed field)
	EmailTheme *string `pulumi:"emailTheme,optional"`

	// SMTP server configuration (managed field)
	SmtpServer *SmtpServerConfig `pulumi:"smtpServer,optional"`
}

// SmtpServerConfig represents SMTP configuration for the realm
type SmtpServerConfig struct {
	Host     *string `pulumi:"host,optional"`
	Port     *int    `pulumi:"port,optional"`
	From     *string `pulumi:"from,optional"`
	FromName *string `pulumi:"fromName,optional"`
	StartTls *bool   `pulumi:"startTls,optional"`
	Auth     *bool   `pulumi:"auth,optional"`
	Username *string `pulumi:"username,optional"`
	Password *string `pulumi:"password,optional"`
}

// RealmState represents the output state of a realm resource
type RealmState struct {
	// The ID of the realm (same as name)
	ID string `pulumi:"realmId"`

	// The name of the realm
	Name string `pulumi:"name"`

	// Whether the realm is enabled
	Enabled *bool `pulumi:"enabled,optional"`

	// Display name for the realm
	DisplayName *string `pulumi:"displayName,optional"`

	// HTML display name
	DisplayNameHtml *string `pulumi:"displayNameHtml,optional"`

	// Login theme
	LoginTheme *string `pulumi:"loginTheme,optional"`

	// Account theme
	AccountTheme *string `pulumi:"accountTheme,optional"`

	// Admin theme
	AdminTheme *string `pulumi:"adminTheme,optional"`

	// Email theme
	EmailTheme *string `pulumi:"emailTheme,optional"`

	// SMTP server configuration
	SmtpServer *SmtpServerConfig `pulumi:"smtpServer,optional"`
}

// List of fields that this provider manages (merge strategy)
var managedFields = []string{
	"realm", "enabled", "displayName", "displayNameHtml",
	"loginTheme", "accountTheme", "adminTheme", "emailTheme", "smtpServer",
}

// Annotate provides schema documentation for the Realm resource
func (r *Realm) Annotate(a infer.Annotator) {
	a.Describe(&r, "A Keycloak realm resource that preserves manual changes to unmanaged attributes")
}

// Annotate provides schema documentation for RealmArgs
func (args *RealmArgs) Annotate(a infer.Annotator) {
	a.Describe(&args.Name, "The name of the realm")
	a.Describe(&args.Enabled, "Whether the realm is enabled")
	a.Describe(&args.DisplayName, "Display name shown in the admin console and login pages")
	a.Describe(&args.DisplayNameHtml, "HTML display name for the realm")
	a.Describe(&args.LoginTheme, "Theme used for login pages")
	a.Describe(&args.AccountTheme, "Theme used for account management pages")
	a.Describe(&args.AdminTheme, "Theme used for admin console")
	a.Describe(&args.EmailTheme, "Theme used for email templates")
	a.Describe(&args.SmtpServer, "SMTP server configuration for email sending")

	// Set default values
	a.SetDefault(&args.Enabled, true)
}

// Annotate provides schema documentation for SmtpServerConfig
func (smtp *SmtpServerConfig) Annotate(a infer.Annotator) {
	a.Describe(&smtp.Host, "SMTP server hostname")
	a.Describe(&smtp.Port, "SMTP server port")
	a.Describe(&smtp.From, "From email address")
	a.Describe(&smtp.FromName, "From display name")
	a.Describe(&smtp.StartTls, "Whether to use STARTTLS")
	a.Describe(&smtp.Auth, "Whether SMTP authentication is required")
	a.Describe(&smtp.Username, "SMTP username")
	a.Describe(&smtp.Password, "SMTP password")

	// Set default values
	a.SetDefault(&smtp.Port, 587)
	a.SetDefault(&smtp.StartTls, true)
	a.SetDefault(&smtp.Auth, false)
}

// Annotate provides schema documentation for RealmState
func (state *RealmState) Annotate(a infer.Annotator) {
	a.Describe(&state.ID, "The unique identifier of the realm")
	a.Describe(&state.Name, "The name of the realm")
	a.Describe(&state.Enabled, "Whether the realm is enabled")
	a.Describe(&state.DisplayName, "Display name shown in the admin console and login pages")
	a.Describe(&state.DisplayNameHtml, "HTML display name for the realm")
	a.Describe(&state.LoginTheme, "Theme used for login pages")
	a.Describe(&state.AccountTheme, "Theme used for account management pages")
	a.Describe(&state.AdminTheme, "Theme used for admin console")
	a.Describe(&state.EmailTheme, "Theme used for email templates")
	a.Describe(&state.SmtpServer, "SMTP server configuration for email sending")
}

// Create implementation for the realm resource
func (r *Realm) Create(ctx context.Context, req infer.CreateRequest[RealmArgs]) (infer.CreateResponse[RealmState], error) {
	config := getProviderConfig(ctx)
	client := getKeycloakClient(ctx)

	// Authenticate to get fresh token
	token, err := client.LoginAdmin(ctx, config.Username, config.Password, *config.Realm)
	if err != nil {
		return infer.CreateResponse[RealmState]{}, fmt.Errorf("failed to authenticate: %w", err)
	}

	// Check if realm already exists
	exists, err := realmExists(ctx, client, req.Inputs.Name)
	if err != nil {
		return infer.CreateResponse[RealmState]{}, fmt.Errorf("failed to check if realm exists: %w", err)
	}

	if req.DryRun {
		return infer.CreateResponse[RealmState]{
			ID: req.Inputs.Name,
			Output: RealmState{
				ID:              req.Inputs.Name,
				Name:            req.Inputs.Name,
				Enabled:         req.Inputs.Enabled,
				DisplayName:     req.Inputs.DisplayName,
				DisplayNameHtml: req.Inputs.DisplayNameHtml,
				LoginTheme:      req.Inputs.LoginTheme,
				AccountTheme:    req.Inputs.AccountTheme,
				AdminTheme:      req.Inputs.AdminTheme,
				EmailTheme:      req.Inputs.EmailTheme,
				SmtpServer:      req.Inputs.SmtpServer,
			},
		}, nil
	}

	// Create realm if it doesn't exist with minimal config
	if !exists {
		enabled := true
		if req.Inputs.Enabled != nil {
			enabled = *req.Inputs.Enabled
		}

		minimalRealm := gocloak.RealmRepresentation{
			Realm:   &req.Inputs.Name,
			Enabled: &enabled,
		}

		if req.Inputs.DisplayName != nil {
			minimalRealm.DisplayName = req.Inputs.DisplayName
		}

		_, err = client.CreateRealm(ctx, token.AccessToken, minimalRealm)
		if err != nil {
			return infer.CreateResponse[RealmState]{}, fmt.Errorf("failed to create realm: %w", err)
		}
	}

	// Update with managed fields only (merge strategy)
	err = r.updateManagedFields(ctx, client, token.AccessToken, req.Inputs)
	if err != nil {
		return infer.CreateResponse[RealmState]{}, fmt.Errorf("failed to update managed fields: %w", err)
	}

	// Read the current state
	state, err := r.readRealmState(ctx, client, token.AccessToken, req.Inputs.Name)
	if err != nil {
		return infer.CreateResponse[RealmState]{}, fmt.Errorf("failed to read realm state: %w", err)
	}

	return infer.CreateResponse[RealmState]{
		ID:     req.Inputs.Name,
		Output: state,
	}, nil
}

// Update implementation - only updates managed fields
func (r *Realm) Update(ctx context.Context, req infer.UpdateRequest[RealmArgs, RealmState]) (infer.UpdateResponse[RealmState], error) {
	config := getProviderConfig(ctx)
	client := getKeycloakClient(ctx)

	// Authenticate to get fresh token
	token, err := client.LoginAdmin(ctx, config.Username, config.Password, *config.Realm)
	if err != nil {
		return infer.UpdateResponse[RealmState]{}, fmt.Errorf("failed to authenticate: %w", err)
	}

	if req.DryRun {
		return infer.UpdateResponse[RealmState]{
			Output: RealmState{
				ID:              req.Inputs.Name,
				Name:            req.Inputs.Name,
				Enabled:         req.Inputs.Enabled,
				DisplayName:     req.Inputs.DisplayName,
				DisplayNameHtml: req.Inputs.DisplayNameHtml,
				LoginTheme:      req.Inputs.LoginTheme,
				AccountTheme:    req.Inputs.AccountTheme,
				AdminTheme:      req.Inputs.AdminTheme,
				EmailTheme:      req.Inputs.EmailTheme,
				SmtpServer:      req.Inputs.SmtpServer,
			},
		}, nil
	}

	// Update only managed fields (merge strategy)
	err = r.updateManagedFields(ctx, client, token.AccessToken, req.Inputs)
	if err != nil {
		return infer.UpdateResponse[RealmState]{}, fmt.Errorf("failed to update managed fields: %w", err)
	}

	// Read the current state
	state, err := r.readRealmState(ctx, client, token.AccessToken, req.Inputs.Name)
	if err != nil {
		return infer.UpdateResponse[RealmState]{}, fmt.Errorf("failed to read realm state: %w", err)
	}

	return infer.UpdateResponse[RealmState]{
		Output: state,
	}, nil
}

// Delete implementation
func (r *Realm) Delete(ctx context.Context, req infer.DeleteRequest[RealmState]) (infer.DeleteResponse, error) {
	config := getProviderConfig(ctx)
	client := getKeycloakClient(ctx)

	// Authenticate to get fresh token
	token, err := client.LoginAdmin(ctx, config.Username, config.Password, *config.Realm)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("failed to authenticate: %w", err)
	}

	err = client.DeleteRealm(ctx, token.AccessToken, req.State.Name)
	if err != nil {
		// Check if realm was already deleted
		exists, checkErr := realmExists(ctx, client, req.State.Name)
		if checkErr == nil && !exists {
			// Realm already deleted, that's okay
			return infer.DeleteResponse{}, nil
		}
		return infer.DeleteResponse{}, fmt.Errorf("failed to delete realm: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Read implementation
func (r *Realm) Read(ctx context.Context, req infer.ReadRequest[RealmArgs, RealmState]) (infer.ReadResponse[RealmArgs, RealmState], error) {
	config := getProviderConfig(ctx)
	client := getKeycloakClient(ctx)

	// Authenticate to get fresh token
	token, err := client.LoginAdmin(ctx, config.Username, config.Password, *config.Realm)
	if err != nil {
		return infer.ReadResponse[RealmArgs, RealmState]{}, fmt.Errorf("failed to authenticate: %w", err)
	}

	// Check if realm exists
	exists, err := realmExists(ctx, client, req.ID)
	if err != nil {
		return infer.ReadResponse[RealmArgs, RealmState]{}, fmt.Errorf("failed to check if realm exists: %w", err)
	}

	if !exists {
		// Realm doesn't exist, return empty response
		return infer.ReadResponse[RealmArgs, RealmState]{}, nil
	}

	// Read the current state
	state, err := r.readRealmState(ctx, client, token.AccessToken, req.ID)
	if err != nil {
		return infer.ReadResponse[RealmArgs, RealmState]{}, fmt.Errorf("failed to read realm state: %w", err)
	}

	return infer.ReadResponse[RealmArgs, RealmState]{
		ID:     req.ID,
		Inputs: req.Inputs,
		State:  state,
	}, nil
}

// updateManagedFields updates only the fields managed by this provider
func (r *Realm) updateManagedFields(ctx context.Context, client *gocloak.GoCloak, token string, args RealmArgs) error {
	// Get current realm to preserve unmanaged fields
	currentRealm, err := client.GetRealm(ctx, token, args.Name)
	if err != nil {
		return fmt.Errorf("failed to get current realm: %w", err)
	}

	// Create update payload with only managed fields
	updateRealm := *currentRealm

	// Update managed fields
	if args.Enabled != nil {
		updateRealm.Enabled = args.Enabled
	}

	if args.DisplayName != nil {
		updateRealm.DisplayName = args.DisplayName
	}

	if args.DisplayNameHtml != nil {
		updateRealm.DisplayNameHTML = args.DisplayNameHtml
	}

	if args.LoginTheme != nil {
		updateRealm.LoginTheme = args.LoginTheme
	}

	if args.AccountTheme != nil {
		updateRealm.AccountTheme = args.AccountTheme
	}

	if args.AdminTheme != nil {
		updateRealm.AdminTheme = args.AdminTheme
	}

	if args.EmailTheme != nil {
		updateRealm.EmailTheme = args.EmailTheme
	}

	if args.SmtpServer != nil {
		smtpConfig := convertSmtpConfig(args.SmtpServer)
		updateRealm.SMTPServer = &smtpConfig
	}

	err = client.UpdateRealm(ctx, token, updateRealm)
	if err != nil {
		return fmt.Errorf("failed to update realm: %w", err)
	}

	return nil
}

// readRealmState reads the current state, focusing on managed fields
func (r *Realm) readRealmState(ctx context.Context, client *gocloak.GoCloak, token, realmName string) (RealmState, error) {
	realm, err := client.GetRealm(ctx, token, realmName)
	if err != nil {
		return RealmState{}, fmt.Errorf("failed to get realm: %w", err)
	}

	state := RealmState{
		ID:   *realm.Realm,
		Name: *realm.Realm,
	}

	// Only populate managed fields
	if realm.Enabled != nil {
		state.Enabled = realm.Enabled
	}

	if realm.DisplayName != nil {
		state.DisplayName = realm.DisplayName
	}

	if realm.DisplayNameHTML != nil {
		state.DisplayNameHtml = realm.DisplayNameHTML
	}

	if realm.LoginTheme != nil {
		state.LoginTheme = realm.LoginTheme
	}

	if realm.AccountTheme != nil {
		state.AccountTheme = realm.AccountTheme
	}

	if realm.AdminTheme != nil {
		state.AdminTheme = realm.AdminTheme
	}

	if realm.EmailTheme != nil {
		state.EmailTheme = realm.EmailTheme
	}

	if realm.SMTPServer != nil {
		state.SmtpServer = convertFromKeycloakSmtp(*realm.SMTPServer)
	}

	return state, nil
}

// Helper functions
func realmExists(ctx context.Context, client *gocloak.GoCloak, realmName string) (bool, error) {
	config := getProviderConfig(ctx)
	token, err := client.LoginAdmin(ctx, config.Username, config.Password, *config.Realm)
	if err != nil {
		return false, fmt.Errorf("failed to authenticate: %w", err)
	}

	_, err = client.GetRealm(ctx, token.AccessToken, realmName)
	if err != nil {
		// If it's a 404-like error, realm doesn't exist
		if err.Error() == "404" || err.Error() == "realm not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func convertSmtpConfig(smtp *SmtpServerConfig) map[string]string {
	if smtp == nil {
		return nil
	}

	result := make(map[string]string)

	if smtp.Host != nil {
		result["host"] = *smtp.Host
	}

	if smtp.Port != nil {
		result["port"] = fmt.Sprintf("%d", *smtp.Port)
	}

	if smtp.From != nil {
		result["from"] = *smtp.From
	}

	if smtp.FromName != nil {
		result["fromDisplayName"] = *smtp.FromName
	}

	if smtp.StartTls != nil {
		if *smtp.StartTls {
			result["starttls"] = "true"
		} else {
			result["starttls"] = "false"
		}
	}

	if smtp.Auth != nil && *smtp.Auth {
		result["auth"] = "true"
		if smtp.Username != nil {
			result["user"] = *smtp.Username
		}
		if smtp.Password != nil {
			result["password"] = *smtp.Password
		}
	} else {
		result["auth"] = "false"
	}

	return result
}

func convertFromKeycloakSmtp(keycloakSmtp map[string]string) *SmtpServerConfig {
	if len(keycloakSmtp) == 0 {
		return nil
	}

	smtp := &SmtpServerConfig{}

	if host, ok := keycloakSmtp["host"]; ok {
		smtp.Host = &host
	}

	if port, ok := keycloakSmtp["port"]; ok {
		if portInt := parseInt(port); portInt != nil {
			smtp.Port = portInt
		}
	}

	if from, ok := keycloakSmtp["from"]; ok {
		smtp.From = &from
	}

	if fromName, ok := keycloakSmtp["fromDisplayName"]; ok {
		smtp.FromName = &fromName
	}

	if starttls, ok := keycloakSmtp["starttls"]; ok {
		starttlsBool := starttls == "true"
		smtp.StartTls = &starttlsBool
	}

	if auth, ok := keycloakSmtp["auth"]; ok {
		authBool := auth == "true"
		smtp.Auth = &authBool

		if authBool {
			if user, ok := keycloakSmtp["user"]; ok {
				smtp.Username = &user
			}
			if password, ok := keycloakSmtp["password"]; ok {
				smtp.Password = &password
			}
		}
	}

	return smtp
}

func parseInt(s string) *int {
	if s == "" {
		return nil
	}

	// Simple integer parsing
	result := 0
	for _, char := range s {
		if char < '0' || char > '9' {
			return nil
		}
		result = result*10 + int(char-'0')
	}
	return &result
}
