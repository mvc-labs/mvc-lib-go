// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC

package interpreter

import (
	"testing"

	"github.com/mvc-labs/mvc-lib-go/bscript"
	"github.com/mvc-labs/mvc-lib-go/bscript/interpreter/errs"
)

// TestOpcodeDisabled tests the opcodeDisabled function manually because all
// disabled opcodes result in a script execution failure when executed normally,
// so the function is not called under normal circumstances.
func TestOpcodeDisabled(t *testing.T) {
	t.Parallel()

	tests := []byte{bscript.Op2MUL, bscript.Op2DIV}
	for _, opcodeVal := range tests {
		pop := ParsedOpcode{op: opcodeArray[opcodeVal], Data: nil}
		err := opcodeDisabled(&pop, nil)
		if !errs.IsErrorCode(err, errs.ErrDisabledOpcode) {
			t.Errorf("opcodeDisabled: unexpected error - got %v, "+
				"want %v", err, errs.ErrDisabledOpcode)
			continue
		}
	}
}
