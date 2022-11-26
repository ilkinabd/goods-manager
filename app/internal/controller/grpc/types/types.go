package types

import (
	pbCommon "github.com/ilkinabd/goods-contracts/gen/go/common/v1"
	"github.com/ilkinabd/goods-manager/app/pkg/api/filter"
)

type FilterOperator = filter.Operator

func IntOperatorFromPB(e pbCommon.IntFilterField_Operator) FilterOperator {
	switch e {
	case pbCommon.IntFilterField_OPERATOR_EQ:
		return filter.OperatorEq
	case pbCommon.IntFilterField_OPERATOR_NEQ:
		return filter.OperatorNotEq
	case pbCommon.IntFilterField_OPERATOR_LT:
		return filter.OperatorLowerThan
	case pbCommon.IntFilterField_OPERATOR_LTE:
		return filter.OperatorLowerThanEq
	case pbCommon.IntFilterField_OPERATOR_GT:
		return filter.OperatorGreaterThan
	case pbCommon.IntFilterField_OPERATOR_GTE:
		return filter.OperatorGreaterThanEq
	default:
		return ""
	}
}

func StringOperatorFromPB(e pbCommon.StringFilterField_Operator) FilterOperator {
	switch e {
	case pbCommon.StringFilterField_OPERATOR_EQ:
		return filter.OperatorEq
	case pbCommon.StringFilterField_OPERATOR_NEQ:
		return filter.OperatorNotEq
	case pbCommon.StringFilterField_OPERATOR_LIKE:
		return filter.OperatorLike
	default:
		return ""
	}
}
