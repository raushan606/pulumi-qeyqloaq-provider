package provider

import (
	"context"
	"fmt"

	gocloak "github.com/Nerzal/gocloak/v13"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// Realm represents a Keycloak realm resource with merge strategy
// This provider only updates fields that are explicitly managed,
// preserving manual changes to other realm attributes in Keycloak UI
type Realm struct{}

type RealmArgs struct {
	Name            string            `pulumi:"name"`
	Enabled         *bool             `pulumi:"enabled,optional"`
	DisplayName     *string           `pulumi:"displayName,optional"`
	DisplayNameHtml *string           `pulumi:"displayNameHtml,optional"`
	LoginTheme      *string           `pulumi:"loginTheme,optional"`
	AccountTheme    *string           `pulumi:"accountTheme,optional"`
	AdminTheme      *string           `pulumi:"adminTheme,optional"`
	EmailTheme      *string           `pulumi:"emailTheme,optional"`
	SmtpServer      *SmtpServerConfig `pulumi:"smtpServer,optional"`
}

func (args RealmArgs) toKeycloakRealm() gocloak.RealmRepresentation {
	keycloakRealmRepresentation := gocloak.RealmRepresentation{
		Realm: &args.Name,
	}

	if args.Enabled != nil {
		keycloakRealmRepresentation.Enabled = args.Enabled
	} else {
		enabled := true
		keycloakRealmRepresentation.Enabled = &enabled
	}

	if args.DisplayName != nil {
		keycloakRealmRepresentation.DisplayName = args.DisplayName
	}
	if args.DisplayNameHtml != nil {
		keycloakRealmRepresentation.DisplayNameHTML = args.DisplayNameHtml
	}
	if args.LoginTheme != nil {
		keycloakRealmRepresentation.LoginTheme = args.LoginTheme
	}
	if args.AccountTheme != nil {
		keycloakRealmRepresentation.AccountTheme = args.AccountTheme
	}
	if args.AdminTheme != nil {
		keycloakRealmRepresentation.AdminTheme = args.AdminTheme
	}
	if args.EmailTheme != nil {
		keycloakRealmRepresentation.EmailTheme = args.EmailTheme
	}
	if args.SmtpServer != nil {
		smtpConfig := convertSmtpConfig(args.SmtpServer)
		keycloakRealmRepresentation.SMTPServer = &smtpConfig
	}
	return keycloakRealmRepresentation
}

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

type RealmState struct {
	ID              string            `pulumi:"realmId"` // The ID of the realm (same as name)
	Name            string            `pulumi:"name"`
	Enabled         *bool             `pulumi:"enabled,optional"`
	DisplayName     *string           `pulumi:"displayName,optional"`
	DisplayNameHtml *string           `pulumi:"displayNameHtml,optional"`
	LoginTheme      *string           `pulumi:"loginTheme,optional"`
	AccountTheme    *string           `pulumi:"accountTheme,optional"`
	AdminTheme      *string           `pulumi:"adminTheme,optional"`
	EmailTheme      *string           `pulumi:"emailTheme,optional"`
	SmtpServer      *SmtpServerConfig `pulumi:"smtpServer,optional"`
}

// Annotate provides schema documentation for the Realm resource
func (r *Realm) Annotate(a infer.Annotator) {
	a.Describe(&r, "A Keycloak realm resource that preserves manual changes to unmanaged attributes")
}

// WireDependencies controls how outputs and secrets flow through values
func (Realm) WireDependencies(f infer.FieldSelector, args *RealmArgs, state *RealmState) {
	f.OutputField(&state.Name).DependsOn(f.InputField(&args.Name))
	f.OutputField(&state.DisplayName).DependsOn(f.InputField(&args.DisplayName))
	f.OutputField(&state.LoginTheme).DependsOn(f.InputField(&args.LoginTheme))
	f.OutputField(&state.AccountTheme).DependsOn(f.InputField(&args.AccountTheme))
	f.OutputField(&state.AdminTheme).DependsOn(f.InputField(&args.AdminTheme))
	f.OutputField(&state.EmailTheme).DependsOn(f.InputField(&args.EmailTheme))
	f.OutputField(&state.SmtpServer).DependsOn(f.InputField(&args.SmtpServer))
}

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

	a.SetDefault(&args.Enabled, true)
}

func (smtp *SmtpServerConfig) Annotate(a infer.Annotator) {
	a.Describe(&smtp.Host, "SMTP server hostname")
	a.Describe(&smtp.Port, "SMTP server port")
	a.Describe(&smtp.From, "From email address")
	a.Describe(&smtp.FromName, "From display name")
	a.Describe(&smtp.StartTls, "Whether to use STARTTLS")
	a.Describe(&smtp.Auth, "Whether SMTP authentication is required")
	a.Describe(&smtp.Username, "SMTP username")
	a.Describe(&smtp.Password, "SMTP password")

	a.SetDefault(&smtp.Port, 587)
	a.SetDefault(&smtp.StartTls, true)
	a.SetDefault(&smtp.Auth, false)
}

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

func (r *Realm) Create(ctx context.Context, req infer.CreateRequest[RealmArgs]) (infer.CreateResponse[RealmState], error) {
	config := infer.GetConfig[ProviderConfig](ctx)

	client := gocloak.NewClient(config.URL)
	token, err := client.LoginAdmin(ctx, config.Username, config.Password, *config.Realm)
	if err != nil {
		return infer.CreateResponse[RealmState]{}, fmt.Errorf("failed to authenticate: %w", err)
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

	_, err = client.CreateRealm(ctx, token.AccessToken, req.Inputs.toKeycloakRealm())
	if err != nil {
		return infer.CreateResponse[RealmState]{}, fmt.Errorf("failed to create realm: %w", err)
	}

	// Read the current state
	state, err := readRealmState(ctx, client, token.AccessToken, req.Inputs.Name)
	if err != nil {
		return infer.CreateResponse[RealmState]{}, fmt.Errorf("failed to read realm state: %w", err)
	}

	return infer.CreateResponse[RealmState]{
		ID:     req.Inputs.Name,
		Output: state,
	}, nil
}

func (*Realm) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[RealmArgs], error) {
	args, f, err := infer.DefaultCheck[RealmArgs](ctx, req.NewInputs)
	return infer.CheckResponse[RealmArgs]{
		Inputs:   args,
		Failures: f,
	}, err
}

// // Update implementation - only updates managed fields
// func (r *Realm) Update(ctx context.Context, req infer.UpdateRequest[RealmArgs, RealmState]) (infer.UpdateResponse[RealmState], error) {
// 	config := getProviderConfig(ctx)
// 	client := getKeycloakClient(ctx)

// 	// Authenticate to get fresh token
// 	token, err := client.LoginAdmin(ctx, config.Username, config.Password, *config.Realm)
// 	if err != nil {
// 		return infer.UpdateResponse[RealmState]{}, fmt.Errorf("failed to authenticate: %w", err)
// 	}

// 	if req.DryRun {
// 		return infer.UpdateResponse[RealmState]{
// 			Output: RealmState{
// 				ID:              req.Inputs.Name,
// 				Name:            req.Inputs.Name,
// 				Enabled:         req.Inputs.Enabled,
// 				DisplayName:     req.Inputs.DisplayName,
// 				DisplayNameHtml: req.Inputs.DisplayNameHtml,
// 				LoginTheme:      req.Inputs.LoginTheme,
// 				AccountTheme:    req.Inputs.AccountTheme,
// 				AdminTheme:      req.Inputs.AdminTheme,
// 				EmailTheme:      req.Inputs.EmailTheme,
// 				SmtpServer:      req.Inputs.SmtpServer,
// 			},
// 		}, nil
// 	}

// 	// Update only managed fields (merge strategy)
// 	err = r.updateManagedFields(ctx, client, token.AccessToken, req.Inputs)
// 	if err != nil {
// 		return infer.UpdateResponse[RealmState]{}, fmt.Errorf("failed to update managed fields: %w", err)
// 	}

// 	// Read the current state
// 	state, err := r.readRealmState(ctx, client, token.AccessToken, req.Inputs.Name)
// 	if err != nil {
// 		return infer.UpdateResponse[RealmState]{}, fmt.Errorf("failed to read realm state: %w", err)
// 	}

// 	return infer.UpdateResponse[RealmState]{
// 		Output: state,
// 	}, nil
// }

// // Delete implementation
// func (r *Realm) Delete(ctx context.Context, req infer.DeleteRequest[RealmState]) (infer.DeleteResponse, error) {
// 	config := getProviderConfig(ctx)
// 	client := getKeycloakClient(ctx)

// 	// Authenticate to get fresh token
// 	token, err := client.LoginAdmin(ctx, config.Username, config.Password, *config.Realm)
// 	if err != nil {
// 		return infer.DeleteResponse{}, fmt.Errorf("failed to authenticate: %w", err)
// 	}

// 	err = client.DeleteRealm(ctx, token.AccessToken, req.State.Name)
// 	if err != nil {
// 		// Check if realm was already deleted
// 		exists, checkErr := realmExists(ctx, client, req.State.Name)
// 		if checkErr == nil && !exists {
// 			// Realm already deleted, that's okay
// 			return infer.DeleteResponse{}, nil
// 		}
// 		return infer.DeleteResponse{}, fmt.Errorf("failed to delete realm: %w", err)
// 	}

// 	return infer.DeleteResponse{}, nil
// }

// // Read implementation
// func (r *Realm) Read(ctx context.Context, req infer.ReadRequest[RealmArgs, RealmState]) (infer.ReadResponse[RealmArgs, RealmState], error) {
// 	config := getProviderConfig(ctx)
// 	client := getKeycloakClient(ctx)

// 	// Authenticate to get fresh token
// 	token, err := client.LoginAdmin(ctx, config.Username, config.Password, *config.Realm)
// 	if err != nil {
// 		return infer.ReadResponse[RealmArgs, RealmState]{}, fmt.Errorf("failed to authenticate: %w", err)
// 	}

// 	// Check if realm exists
// 	exists, err := realmExists(ctx, client, req.ID)
// 	if err != nil {
// 		return infer.ReadResponse[RealmArgs, RealmState]{}, fmt.Errorf("failed to check if realm exists: %w", err)
// 	}

// 	if !exists {
// 		// Realm doesn't exist, return empty response
// 		return infer.ReadResponse[RealmArgs, RealmState]{}, nil
// 	}

// 	// Read the current state
// 	state, err := r.readRealmState(ctx, client, token.AccessToken, req.ID)
// 	if err != nil {
// 		return infer.ReadResponse[RealmArgs, RealmState]{}, fmt.Errorf("failed to read realm state: %w", err)
// 	}

// 	return infer.ReadResponse[RealmArgs, RealmState]{
// 		ID:     req.ID,
// 		Inputs: req.Inputs,
// 		State:  state,
// 	}, nil
// }

// updateManagedFields updates only the fields managed by this provider
func updateManagedFields(ctx context.Context, client *gocloak.GoCloak, token string, args RealmArgs) error {
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

// reads the current state, focusing on managed fields
func readRealmState(ctx context.Context, client *gocloak.GoCloak, token, realmName string) (RealmState, error) {
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

func realmExistsWithClient(ctx context.Context, client *gocloak.GoCloak, token, realmName string) (bool, error) {
	_, err := client.GetRealm(ctx, token, realmName)
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
