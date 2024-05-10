package api

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// V2PipelineRequestBuilder builds and executes requests for operations under \api\v2\pipeline
type V2PipelineRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ByUuid gets an item from the go.artefactual.dev/ssclient/kiota.api.v2.pipeline.item collection
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
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/api/v2/pipeline", pathParameters),
    }
    return m
}
// NewV2PipelineRequestBuilder instantiates a new V2PipelineRequestBuilder and sets the default values.
func NewV2PipelineRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*V2PipelineRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewV2PipelineRequestBuilderInternal(urlParams, requestAdapter)
}