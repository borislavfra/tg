// @tg version=1.0.0
// @tg backend=adder.api
// @tg title=`Adder API`
// @tg servers=`http://adder.back-adder-gen.svc.k8s.test;test cluster`
package service

import (
	"context"

	"github.com/seniorGolang/gokit/types/uuid"
	"github.com/seniorGolang/tg/test/service/types"
)

// @tg jsonRPC-server log trace metrics
type Adder interface {
	// @tg summary=`Получение суммы двух чисел`
	Add(ctx context.Context, firstNumber, secondNumber float64) (sum float64, err error)

	// @tg summary=`Получение uuid по id`
	GetUUID(ctx context.Context, id []int) (genUUID uuid.UUID, err error)

	// @tg summary=`Возвращает исходные данные`
	DoNothing(ctx context.Context, thing types.CoolThing) (out types.CoolThing, err error)
}
