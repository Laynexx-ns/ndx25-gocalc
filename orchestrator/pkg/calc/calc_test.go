package calc

//
//import (
//	"testing"
//)
//
//type Test struct {
//	Expression string
//	Expected   float64
//}
//
//func TestCalc(t *testing.T) {
//	testCases := []Test{
//		{"1+1", 1 + 1},
//		{"3+3*6", 3 + 3*6},
//		{"1+8/2*4", 1 + 8/2*4},
//		{"(1+1)*2", (1 + 1) * 2},
//		{"((1+4)*(1+2)+10)*4", ((1+4)*(1+2) + 10) * 4},
//		{"(4+3+2)/(1+2)* 10/3", (4 + 3 + 2) / (1 + 2) * 10 / 3},
//		{"(70/7)*10/((3+2)*(3+7))-2", (70/7)*10/((3+2)*(3+7)) - 2},
//		{"5", 5},
//	}
//
//	calculate(t, testCases)
//
//}
//
//func calculate(t *testing.T, testCases []Test) {
//	for _, tc := range testCases {
//		num, err := Calc(tc.Expression)
//		if err != nil {
//			t.Error(err)
//		}
//		if tc.Expected != num {
//			t.Errorf("Calc(%q) = %f, want %f", tc.Expression, num, tc.Expected)
//		}
//	}
//}
