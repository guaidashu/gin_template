package helper_test

import (
	"gin_template/app/libs/helper"
	"testing"
)

func TestCompareFloat(t *testing.T) {
	ret, err := helper.CompareFloat(32.43421, 32.43422)
	if err != nil {
		t.Error(err)
	}

	t.Log(ret)
}
