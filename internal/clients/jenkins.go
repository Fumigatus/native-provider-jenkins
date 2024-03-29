/*
Copyright 2022 The ANKA SOFTWARE Authors.
*/

package clients

import (
	"context"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane/provider-nativeproviderjenkins/apis/v1alpha1"
)

// Config for vRA Client authentication struct
type Config struct {
	BaseURL      string
	RefreshToken string
}

// NewClient creates new vRA Client with provided vRA Configurations.
func NewClient(c Config) {//*vra.API {
	//transport := GetTransport(c)
	//apiClient := vra.New(transport, strfmt.Default)

	//return apiClient
}

// NewClientWithAuthentication creates new vRA Client with provided vRA Configurations and Credentials.
func NewClientWithAuthentication(c Config, bearerToken string) {//*vra.API {
	//transport := GetTransportWithAuthentication(c, bearerToken)
	//apiClient := vra.New(transport, strfmt.Default)

	//return apiClient
}

// GetTransport returns the REST config
func GetTransport(c Config) {//runtime.ClientTransport {
	//transport := httptransport.New(c.BaseURL, "", nil)
	//transport.SetDebug(true)

	//return transport
}
/*
// GetTransportWithAuthentication returns the REST config with authentication header
func GetTransportWithAuthentication(c Config, bearerToken string) runtime.ClientTransport {
	transport := httptransport.New(c.BaseURL, "", nil)
	transport.SetDebug(true)
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Bearer "+bearerToken)

	return transport
}
*/
// GetConfig constructs a Config that can be used to authenticate to vRA
func GetConfig(ctx context.Context, c client.Client, mg resource.Managed) (*Config, error) {
	switch {
	case mg.GetProviderConfigReference() != nil:
		return UseProviderConfig(ctx, c, mg)
	default:
		return nil, errors.New("providerConfigRef is not given")
	}
}

// UseProviderConfig to produce a config that can be used to authenticate to AWS.
func UseProviderConfig(ctx context.Context, c client.Client, mg resource.Managed) (*Config, error) {
	pc := &v1alpha1.ProviderConfig{}
	if err := c.Get(ctx, types.NamespacedName{Name: mg.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, "cannot get referenced Provider")
	}

	t := resource.NewProviderConfigUsageTracker(c, &v1alpha1.ProviderConfigUsage{})
	if err := t.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, "cannot track ProviderConfig usage")
	}

	switch s := pc.Spec.Credentials.Source; s { //nolint:exhaustive
	case xpv1.CredentialsSourceSecret:
		csr := pc.Spec.Credentials.SecretRef
		if csr == nil {
			return nil, errors.New("no credentials secret referenced")
		}
		s := &corev1.Secret{}
		if err := c.Get(ctx, types.NamespacedName{Namespace: csr.Namespace, Name: csr.Name}, s); err != nil {
			return nil, errors.Wrap(err, "cannot get credentials secret")
		}
		return &Config{BaseURL: "", RefreshToken: string(s.Data[csr.Key])}, nil
	default:
		return nil, errors.Errorf("credentials source %s is not currently supported", s)
	}
}

type FieldOption int

// Field options.
const (
	// FieldRequired causes zero values to be converted to a pointer to the zero
	// value, rather than a nil pointer. AWS Go SDK types use pointer fields,
	// with a nil pointer indicating an unset field. Our ToPtr functions return
	// a nil pointer for a zero values, unless FieldRequired is set.
	FieldRequired FieldOption = iota
)

// String converts the supplied string for use with the AWS Go SDK.
func String(v string, o ...FieldOption) *string {
	for _, fo := range o {
		if fo == FieldRequired && v == "" {
			return &v
		}
	}

	if v == "" {
		return nil
	}

	return &v
}

// StringValue converts the supplied string pointer to a string, returning the
// empty string if the pointer is nil.
// TODO(muvaf): is this really meaningful? why not implement it?
func StringValue(v *string) string {
	return *v
}

// StringSliceToPtr converts the supplied string array to an array of string pointers.
func StringSliceToPtr(slice []string) []*string {
	if slice == nil {
		return nil
	}

	res := make([]*string, len(slice))
	for i, s := range slice {
		res[i] = String(s)
	}
	return res
}

// StringPtrSliceToValue converts the supplied string pointer array to an array of strings.
func StringPtrSliceToValue(slice []*string) []string {
	if slice == nil {
		return nil
	}

	res := make([]string, len(slice))
	for i, s := range slice {
		res[i] = StringValue(s)
	}
	return res
}
