package items

import (
	"basement/main/internal/logg"
	"fmt"
	"testing"
)

func init() {
	// logg.EnableDebugLogger()
}

func TestMoveToOther(t *testing.T) {
	b1 := &Box{Label: "box1"}
	b2 := &Box{Label: "box2"}

	assertInnerBoxesNil(t, b1)
	assertOuterBoxEqNil(t, b1)
	assertInnerBoxesNil(t, b2)
	assertOuterBoxEqNil(t, b2)

	// useBoxNodeInterface(b1)
	// useBoxNodeInterface(b2)

	// Old state:
	// B1
	// B2
	//
	// New state:
	// B2: {B1}
	assertMoveToOtherNoError(t, b1, b2)
	// useBoxNodeInterface(b1)
	// useBoxNodeInterface(b2)

	assertInnerBoxesNil(t, b1)
	assertOuterBoxEq(t, b1, b2)

	assertInnerBoxesLengthEq(t, b2, 1)
	assertFirstInnerBoxEq(t, b2, b1)
	// using another way to check the length
	assertInnerBoxesLengthEq(t, b2, 1)
	assertOuterBoxEqNil(t, b2)
	assertMoveInsideSelfErrorNotNil(t, b1)

	// Old state:
	// B2:{ B1 }
	// B3
	//
	// New state:
	// B3: {B2: {B1}}
	b3 := &Box{Label: "box3"}

	// useBoxNodeInterface(b1)
	// useBoxNodeInterface(b2)
	// useBoxNodeInterface(b3)

	assertMoveToOtherNoError(t, b2, b3)
	assertInnerBoxesNil(t, b1)
	assertOuterBoxEq(t, b1, b2)

	assertInnerBoxesLengthEq(t, b2, 1)
	assertInnerBoxesLengthEq(t, b1.OuterBox, 1)
	assertFirstInnerBoxEq(t, b2, b1)

	assertOuterBoxEq(t, b2, b3)
	assertInnerBoxesLengthEq(t, b3, 1)
	assertFirstInnerBoxEq(t, b3, b2)
	// using another way to check the length
	assertInnerBoxesLengthEq(t, b2.OuterBox, 1)

	// useBoxNodeInterface(b1)
	// useBoxNodeInterface(b2)
	// useBoxNodeInterface(b3)

	// Moving b1 out of b2
	// Old state:
	// B3: {B2: {B1}}
	//
	// New state:
	// B1
	// B3: {B2}
	err := b1.MoveOutOfOtherBox()
	if err != nil {
		t.Error(fmt.Errorf("something wrong while moving box '%v' out: %w", b1.Label, err))
	}

	assertInnerBoxesNil(t, b1)
	assertOuterBoxEqNil(t, b1)

	assertInnerBoxesNil(t, b2)
	assertInnerBoxesLengthEq(t, b2, 0)
	assertOuterBoxEq(t, b2, b3)

	assertInnerBoxesLengthEq(t, b3, 1)
	// assertFirstInnerBoxEq(t, b3, b2)
	assertOuterBoxEqNil(t, b3)
}

func assertMoveToOtherNoError(t *testing.T, b1 *Box, other *Box) {
	t.Helper()
	err := b1.MoveTo(other)
	if err != nil {
		t.Error("something wrong while moving box %w", err)
	}
}

func assertMoveInsideSelfErrorNotNil(t *testing.T, b *Box) {
	t.Helper()
	err := b.MoveTo(b)
	if err == nil {
		t.Errorf("Moving '%v' to itself should return an error that it can't move inside itself, but error is nil.", b.Label)
	}
}

func assertFirstInnerBoxEq(t *testing.T, b1 *Box, b2 *Box) {
	t.Helper()
	if b1.InnerBoxes == nil {
		t.Errorf("'%v' should have InnerBoxes[0] == '%v', but InnerBoxes is nil.", b1.Label, b2.Label)
		return
	}
	if b1.InnerBoxes[0] != b2 {
		t.Errorf("'%v' should have InnerBoxes[0] == '%v'", b1.Label, b2.Label)
	}
}

func assertOuterBoxEq(t *testing.T, b1 *Box, b2 *Box) {
	t.Helper()
	if b1.OuterBox != b2 {
		t.Errorf("'%v' should have '%v' as outer box but has '%v'", b1.Label, b2.Label, b1.OuterBox.Label)
	}
}

func assertOuterBoxNotEq(t *testing.T, b1 *Box, b2 *Box) {
	t.Helper()
	if b1.OuterBox == b2 {
		t.Errorf("'%v' should not have '%v' as outer box but has '%v'", b1.Label, b2.Label, b1.OuterBox.Label)
	}
}

func assertOuterBoxEqNil(t *testing.T, b *Box) {
	t.Helper()
	if b.OuterBox != nil {
		t.Errorf("'%v' outer box is not nil, '%v'", b.Label, b.OuterBox.Label)
	}

}

func assertInnerBoxesLengthEq(t *testing.T, b *Box, length int) {
	t.Helper()
	if len(b.InnerBoxes) != length {
		t.Errorf("'%v' should have InnerBoxes of length %v, has length %v", b.Label, length, len(b.InnerBoxes))
	}

}

func assertInnerBoxesNil(t *testing.T, b *Box) {
	t.Helper()
	if b.InnerBoxes != nil {
		t.Errorf("'%v' inner boxes is not nil. %v", b.Label, b.InnerBoxes)
	}

}

func useBoxNodeInterface(box BoxNode) {
	mbox := box.Self()
	boxes := fmt.Sprintf("\n-------------------------\nBox \"%v\": \n", mbox.Label)

	boxes += fmt.Sprintf("%v.OuterBox:\n", mbox.Label)
	if mbox.OuterBox != nil {
		boxes += fmt.Sprintf("\t- %v\n", mbox.OuterBox.Label)
	} else {
		boxes += fmt.Sprint("\t- no outer box\n")
	}

	boxes += fmt.Sprintf("%v.InnerBoxes:\n", mbox.Label)
	if len(mbox.InnerBoxes) != 0 {
		for _, b := range mbox.InnerBoxes {
			boxes += fmt.Sprintf("\t- %v\n", b.Label)
		}
	} else {
		boxes += fmt.Sprint("\t- no inner boxes\n")
	}
	// boxes += fmt.Sprint("This box is inside: ", mbox.OuterBox.Label)
	// boxes += fmt.Sprintf("Outer box has: %v", mbox.OuterBox.InnerBoxes[0].Label)
	logg.Debug(boxes, "-------------------------\n")
}
