// @tg version=1.0.0
// @tg backend=adder.todolist
// @tg title=`Adder todolist API`
// @tg servers=`http://adder.back-adder-gen.svc.k8s.test;test cluster`
//go:generate go get -u github.com/seniorGolang/tg/cmd/tg
//go:generate tg transport --services . --out ../transport --outSwagger ../swagger.yaml
package service

import (
	"context"

	"github.com/seniorGolang/tg/test/service/types"
)

// @tg http-server log trace metrics
type TodoList interface {
	// @tg http-method=POST
	// @tg http-path=adder/todolist/create
	// @tg http-response=git.wildberries.ru/suppliers-portal-eu/back-adder-gen/pkg/adder/transport/todolist:CreateResponseHandler
	Create(ctx context.Context, items []types.RawItem) (response types.CreateGetResponse, err error)

	// @tg http-method=GET
	// @tg http-path=adder/todolist/get
	// @tg http-response=git.wildberries.ru/suppliers-portal-eu/back-adder-gen/pkg/adder/transport/todolist:GetResponseHandler
	Get(ctx context.Context) (response types.CreateGetResponse, err error)

	// @tg http-method=POST
	// @tg http-path=adder/todolist/add
	// @tg http-response=git.wildberries.ru/suppliers-portal-eu/back-adder-gen/pkg/adder/transport/todolist:AddResponseHandler
	Add(ctx context.Context, todo types.RawItem) (response types.AddDeleteResponse, err error)

	// @tg http-method=PATCH
	// @tg http-path=adder/todolist/update
	// @tg http-response=git.wildberries.ru/suppliers-portal-eu/back-adder-gen/pkg/adder/transport/todolist:UpdateResponseHandler
	Update(ctx context.Context, todo types.ListItem) (response types.AddDeleteResponse, err error)

	// @tg http-method=DELETE
	// @tg http-path=adder/todolist/delete
	// @tg http-response=git.wildberries.ru/suppliers-portal-eu/back-adder-gen/pkg/adder/transport/todolist:DeleteResponseHandler
	Delete(ctx context.Context, id int) (response types.AddDeleteResponse, err error)
}
