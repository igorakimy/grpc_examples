// Code generated by MockGen. DO NOT EDIT.
// Source: .\server\ecommerce\product_info_grpc.pb.go

// Package mock_ecommerce is a generated GoMock package.
package mock_ecommerce

import (
	context "context"
	ecommerce "productinfo/server/ecommerce"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// MockProductInfoClient is a mock of ProductInfoClient interface.
type MockProductInfoClient struct {
	ctrl     *gomock.Controller
	recorder *MockProductInfoClientMockRecorder
}

// MockProductInfoClientMockRecorder is the mock recorder for MockProductInfoClient.
type MockProductInfoClientMockRecorder struct {
	mock *MockProductInfoClient
}

// NewMockProductInfoClient creates a new mock instance.
func NewMockProductInfoClient(ctrl *gomock.Controller) *MockProductInfoClient {
	mock := &MockProductInfoClient{ctrl: ctrl}
	mock.recorder = &MockProductInfoClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductInfoClient) EXPECT() *MockProductInfoClientMockRecorder {
	return m.recorder
}

// AddProduct mocks base method.
func (m *MockProductInfoClient) AddProduct(ctx context.Context, in *ecommerce.Product, opts ...grpc.CallOption) (*wrapperspb.StringValue, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddProduct", varargs...)
	ret0, _ := ret[0].(*wrapperspb.StringValue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddProduct indicates an expected call of AddProduct.
func (mr *MockProductInfoClientMockRecorder) AddProduct(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProduct", reflect.TypeOf((*MockProductInfoClient)(nil).AddProduct), varargs...)
}

// GetProduct mocks base method.
func (m *MockProductInfoClient) GetProduct(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*ecommerce.Product, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetProduct", varargs...)
	ret0, _ := ret[0].(*ecommerce.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProduct indicates an expected call of GetProduct.
func (mr *MockProductInfoClientMockRecorder) GetProduct(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProduct", reflect.TypeOf((*MockProductInfoClient)(nil).GetProduct), varargs...)
}

// MockProductInfoServer is a mock of ProductInfoServer interface.
type MockProductInfoServer struct {
	ctrl     *gomock.Controller
	recorder *MockProductInfoServerMockRecorder
}

// MockProductInfoServerMockRecorder is the mock recorder for MockProductInfoServer.
type MockProductInfoServerMockRecorder struct {
	mock *MockProductInfoServer
}

// NewMockProductInfoServer creates a new mock instance.
func NewMockProductInfoServer(ctrl *gomock.Controller) *MockProductInfoServer {
	mock := &MockProductInfoServer{ctrl: ctrl}
	mock.recorder = &MockProductInfoServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductInfoServer) EXPECT() *MockProductInfoServerMockRecorder {
	return m.recorder
}

// AddProduct mocks base method.
func (m *MockProductInfoServer) AddProduct(arg0 context.Context, arg1 *ecommerce.Product) (*wrapperspb.StringValue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProduct", arg0, arg1)
	ret0, _ := ret[0].(*wrapperspb.StringValue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddProduct indicates an expected call of AddProduct.
func (mr *MockProductInfoServerMockRecorder) AddProduct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProduct", reflect.TypeOf((*MockProductInfoServer)(nil).AddProduct), arg0, arg1)
}

// GetProduct mocks base method.
func (m *MockProductInfoServer) GetProduct(arg0 context.Context, arg1 *wrapperspb.StringValue) (*ecommerce.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProduct", arg0, arg1)
	ret0, _ := ret[0].(*ecommerce.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProduct indicates an expected call of GetProduct.
func (mr *MockProductInfoServerMockRecorder) GetProduct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProduct", reflect.TypeOf((*MockProductInfoServer)(nil).GetProduct), arg0, arg1)
}

// MockUnsafeProductInfoServer is a mock of UnsafeProductInfoServer interface.
type MockUnsafeProductInfoServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeProductInfoServerMockRecorder
}

// MockUnsafeProductInfoServerMockRecorder is the mock recorder for MockUnsafeProductInfoServer.
type MockUnsafeProductInfoServerMockRecorder struct {
	mock *MockUnsafeProductInfoServer
}

// NewMockUnsafeProductInfoServer creates a new mock instance.
func NewMockUnsafeProductInfoServer(ctrl *gomock.Controller) *MockUnsafeProductInfoServer {
	mock := &MockUnsafeProductInfoServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeProductInfoServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeProductInfoServer) EXPECT() *MockUnsafeProductInfoServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedProductInfoServer mocks base method.
func (m *MockUnsafeProductInfoServer) mustEmbedUnimplementedProductInfoServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedProductInfoServer")
}

// mustEmbedUnimplementedProductInfoServer indicates an expected call of mustEmbedUnimplementedProductInfoServer.
func (mr *MockUnsafeProductInfoServerMockRecorder) mustEmbedUnimplementedProductInfoServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedProductInfoServer", reflect.TypeOf((*MockUnsafeProductInfoServer)(nil).mustEmbedUnimplementedProductInfoServer))
}
