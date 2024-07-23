package api

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// V2LocationDefaultRequestBuilder builds and executes requests for operations under \api\v2\location\default
type V2LocationDefaultRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ByPurpose gets an item from the xgo.artefactual.dev/ssclient/kiota.api.v2.location.default.item collection
// returns a *V2LocationDefaultWithPurposeItemRequestBuilder when successful
func (m *V2LocationDefaultRequestBuilder) ByPurpose(purpose string)(*V2LocationDefaultWithPurposeItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if purpose != "" {
        urlTplParams["purpose"] = purpose
    }
    return NewV2LocationDefaultWithPurposeItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewV2LocationDefaultRequestBuilderInternal instantiates a new V2LocationDefaultRequestBuilder and sets the default values.
func NewV2LocationDefaultRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*V2LocationDefaultRequestBuilder) {
    m := &V2LocationDefaultRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/api/v2/location/default", pathParameters),
    }
    return m
}
// NewV2LocationDefaultRequestBuilder instantiates a new V2LocationDefaultRequestBuilder and sets the default values.
func NewV2LocationDefaultRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*V2LocationDefaultRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewV2LocationDefaultRequestBuilderInternal(urlParams, requestAdapter)
}
