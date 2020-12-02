// @gtg version=1.0.0
// @gtg backend=adder.todolist
// @gtg title=`Adder todolist API`
// @gtg servers=`http://adder.back-adder-gen.svc.k8s.test;test cluster`
//go:generate go get -u github.com/seniorGolang/tg/cmd/tg
//go:generate tg transport --services . --out ../transport --outSwagger ../swagger.yaml
package service

import (
	"context"

	"github.com/seniorGolang/tg/test/service/types"
)

// @gtg http-server log metrics
type GTGTodoList interface {
	// @gtg http-server-method=POST
	// @gtg http-server-uri-path=adder/todolist/create
	Create(ctx context.Context, items []types.RawItem) (response types.CreateGetResponse, err error)

	// @gtg http-server-method=GET
	// @gtg http-server-uri-path=adder/todolist/get
	Get(ctx context.Context) (response types.CreateGetResponse, err error)

	// @gtg http-server-method=POST
	// @gtg http-server-uri-path=adder/todolist/add
	Add(ctx context.Context, todo types.RawItem) (response types.AddDeleteResponse, err error)

	// @gtg http-server-method=PATCH
	// @gtg http-server-uri-path=adder/todolist/update
	Update(ctx context.Context, todo types.ListItem) (response types.AddDeleteResponse, err error)

	// @gtg http-server-method=DELETE
	// @gtg http-server-uri-path=adder/todolist/delete
	Delete(ctx context.Context, id int) (response types.AddDeleteResponse, err error)
}
