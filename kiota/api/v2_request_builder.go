package api

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// V2RequestBuilder builds and executes requests for operations under \api\v2
type V2RequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewV2RequestBuilderInternal instantiates a new V2RequestBuilder and sets the default values.
func NewV2RequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*V2RequestBuilder) {
    m := &V2RequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/api/v2", pathParameters),
    }
    return m
}
// NewV2RequestBuilder instantiates a new V2RequestBuilder and sets the default values.
func NewV2RequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*V2RequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewV2RequestBuilderInternal(urlParams, requestAdapter)
}
// File the file property
// returns a *V2FileRequestBuilder when successful
func (m *V2RequestBuilder) File()(*V2FileRequestBuilder) {
    return NewV2FileRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Location the location property
// returns a *V2LocationRequestBuilder when successful
func (m *V2RequestBuilder) Location()(*V2LocationRequestBuilder) {
    return NewV2LocationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Pipeline the pipeline property
// returns a *V2PipelineRequestBuilder when successful
func (m *V2RequestBuilder) Pipeline()(*V2PipelineRequestBuilder) {
    return NewV2PipelineRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
