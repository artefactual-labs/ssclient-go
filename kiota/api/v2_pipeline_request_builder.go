package api

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    if24bd427556b5f40ce1336ebc33d491ebd9ce71ce225ad2b47d523c1b0f25dee "go.artefactual.dev/ssclient/kiota/models"
)

// V2PipelineRequestBuilder builds and executes requests for operations under \api\v2\pipeline
type V2PipelineRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
type V2PipelineRequestBuilderGetQueryParameters struct {
    Description *string `uriparametername:"description"`
    Uuid *string `uriparametername:"uuid"`
}
// V2PipelineRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type V2PipelineRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *V2PipelineRequestBuilderGetQueryParameters
}
// ByUuid gets an item from the xgo.artefactual.dev/ssclient/kiota.api.v2.pipeline.item collection
// returns a *V2PipelineWithUuItemRequestBuilder when successful
func (m *V2PipelineRequestBuilder) ByUuid(uuid string)(*V2PipelineWithUuItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if uuid != "" {
        urlTplParams["uuid"] = uuid
    }
    return NewV2PipelineWithUuItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewV2PipelineRequestBuilderInternal instantiates a new V2PipelineRequestBuilder and sets the default values.
func NewV2PipelineRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*V2PipelineRequestBuilder) {
    m := &V2PipelineRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/api/v2/pipeline{?description*,uuid*}", pathParameters),
    }
    return m
}
// NewV2PipelineRequestBuilder instantiates a new V2PipelineRequestBuilder and sets the default values.
func NewV2PipelineRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*V2PipelineRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewV2PipelineRequestBuilderInternal(urlParams, requestAdapter)
}
// returns a PipelineListable when successful
func (m *V2PipelineRequestBuilder) Get(ctx context.Context, requestConfiguration *V2PipelineRequestBuilderGetRequestConfiguration)(if24bd427556b5f40ce1336ebc33d491ebd9ce71ce225ad2b47d523c1b0f25dee.PipelineListable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, if24bd427556b5f40ce1336ebc33d491ebd9ce71ce225ad2b47d523c1b0f25dee.CreatePipelineListFromDiscriminatorValue, nil)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(if24bd427556b5f40ce1336ebc33d491ebd9ce71ce225ad2b47d523c1b0f25dee.PipelineListable), nil
}
// returns a *RequestInformation when successful
func (m *V2PipelineRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *V2PipelineRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        if requestConfiguration.QueryParameters != nil {
            requestInfo.AddQueryParameters(*(requestConfiguration.QueryParameters))
        }
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *V2PipelineRequestBuilder when successful
func (m *V2PipelineRequestBuilder) WithUrl(rawUrl string)(*V2PipelineRequestBuilder) {
    return NewV2PipelineRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
