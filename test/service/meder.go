// @tg version=1.0.0
// @tg backend=adder.api
// @tg title=`Adder API`
// @tg servers=`http://adder.back-adder-gen.svc.k8s.test;test cluster`
package service

import (
	"context"
	"github.com/seniorGolang/tg/test/service/types"

	"github.com/seniorGolang/gokit/types/uuid"
)

// @tg jsonRPC-server log trace metrics
type Meder interface {
	// @tg summary=`Получение суммы двух чисел`
	Add(ctx context.Context, firstNumber, secondNumber float64) (sum float64, division []float64, err error)

	// @tg summary=`Получение uuid по id`
	HAHA(ctx context.Context, id []int) (genUUID uuid.UUID, cools []types.CoolThing, err error)

	// @tg summary=`Возвращает исходные данные`
	DoNothing(ctx context.Context, thing types.CoolThing) (out types.CoolThing, testMap map[string]int, err error)
}
