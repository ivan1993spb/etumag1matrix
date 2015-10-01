package main

import "testing"

func TestHandle(t *testing.T) {
	var respfld *respfield
	reqfld := &reqfield{}

	respfld = handle(reqfld)
	if respfld.Status != STATUS_INVALID_INPUT {
		t.Error("TestHandle error status")
	}

	reqfld.Col = []float64{}
	reqfld.Row = []float64{}
	respfld = handle(reqfld)
	if respfld.Status != STATUS_INVALID_INPUT {
		t.Error("TestHandle error status")
	}

	reqfld.Col = []float64{1, 2, 3}
	reqfld.Row = []float64{}
	respfld = handle(reqfld)
	if respfld.Status != STATUS_INVALID_INPUT {
		t.Error("TestHandle error status")
	}

	reqfld.Col = []float64{1, 2, 3}
	reqfld.Row = []float64{3, 2, 1}
	respfld = handle(reqfld)
	if respfld.Status == STATUS_INVALID_INPUT {
		t.Error("TestHandle error status")
	}
	if respfld.Value != 3*1+2*2+3*1 {
		t.Error("TestHandle error status")
	}
}
