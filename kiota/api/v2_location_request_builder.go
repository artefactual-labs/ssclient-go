package api

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i2b7a3625368152c59661ed1a63c26960f7e9cda05d0fbc8e5d79ac57ca250e0a "go.artefactual.dev/ssclient/kiota/models"
)

// V2LocationRequestBuilder builds and executes requests for operations under \api\v2\location
type V2LocationRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// V2LocationRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type V2LocationRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ByUuid gets an item from the go.artefactual.dev/ssclient/kiota.api.v2.location.item collection
// returns a *V2LocationWithUuItemRequestBuilder when successful
func (m *V2LocationRequestBuilder) ByUuid(uuid string)(*V2LocationWithUuItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if uuid != "" {
        urlTplParams["uuid"] = uuid
    }
    return NewV2LocationWithUuItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewV2LocationRequestBuilderInternal instantiates a new V2LocationRequestBuilder and sets the default values.
func NewV2LocationRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*V2LocationRequestBuilder) {
    m := &V2LocationRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/api/v2/location", pathParameters),
    }
    return m
}
// NewV2LocationRequestBuilder instantiates a new V2LocationRequestBuilder and sets the default values.
func NewV2LocationRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*V2LocationRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewV2LocationRequestBuilderInternal(urlParams, requestAdapter)
}
// DefaultEscaped the default property
// returns a *V2LocationDefaultRequestBuilder when successful
func (m *V2LocationRequestBuilder) DefaultEscaped()(*V2LocationDefaultRequestBuilder) {
    return NewV2LocationDefaultRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// returns a LocationListable when successful
func (m *V2LocationRequestBuilder) Get(ctx context.Context, requestConfiguration *V2LocationRequestBuilderGetRequestConfiguration)(i2b7a3625368152c59661ed1a63c26960f7e9cda05d0fbc8e5d79ac57ca250e0a.LocationListable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, i2b7a3625368152c59661ed1a63c26960f7e9cda05d0fbc8e5d79ac57ca250e0a.CreateLocationListFromDiscriminatorValue, nil)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(i2b7a3625368152c59661ed1a63c26960f7e9cda05d0fbc8e5d79ac57ca250e0a.LocationListable), nil
}
// returns a *RequestInformation when successful
func (m *V2LocationRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *V2LocationRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *V2LocationRequestBuilder when successful
func (m *V2LocationRequestBuilder) WithUrl(rawUrl string)(*V2LocationRequestBuilder) {
    return NewV2LocationRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
