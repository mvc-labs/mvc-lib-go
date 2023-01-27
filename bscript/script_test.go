package bscript_test

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mvc-labs/mvc-lib-go/keys/bec"
	"github.com/mvc-labs/mvc-lib-go/keys/bip32"
	"github.com/mvc-labs/mvc-lib-go/keys/chaincfg"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/mvc-labs/mvc-lib-go/bscript"
)

func TestNewP2PKHFromPubKeyStr(t *testing.T) {
	t.Parallel()

	scriptP2PKH, err := bscript.NewP2PKHFromPubKeyStr(
		"023717efaec6761e457f55c8417815505b695209d0bbfed8c3265be425b373c2d6",
	)
	assert.NoError(t, err)
	assert.NotNil(t, scriptP2PKH)
	assert.Equal(t,
		"76a9144d5d1920331b71735a97a606d9734aed83cb3dfa88ac",
		hex.EncodeToString(*scriptP2PKH),
	)
}

func TestNewP2PKHFromPubKey(t *testing.T) {
	t.Parallel()

	pk, _ := hex.DecodeString("023717efaec6761e457f55c8417815505b695209d0bbfed8c3265be425b373c2d6")

	pubkey, err := bec.ParsePubKey(pk, bec.S256())
	assert.NoError(t, err)

	scriptP2PKH, err := bscript.NewP2PKHFromPubKeyEC(pubkey)
	assert.NoError(t, err)
	assert.NotNil(t, scriptP2PKH)
	assert.Equal(t,
		"76a9144d5d1920331b71735a97a606d9734aed83cb3dfa88ac",
		hex.EncodeToString(*scriptP2PKH),
	)
}

func TestNewP2PKHFromBip32ExtKey(t *testing.T) {
	t.Parallel()

	t.Run("output is added", func(t *testing.T) {
		var b [64]byte
		_, err := rand.Read(b[:])
		assert.NoError(t, err)

		key, err := bip32.NewMaster(b[:], &chaincfg.TestNet)
		assert.NoError(t, err)

		script, derivationPath, err := bscript.NewP2PKHFromBip32ExtKey(key)

		assert.NoError(t, err)
		assert.NotEmpty(t, derivationPath)
		assert.NotNil(t, script)
		assert.True(t, script.IsP2PKH())
	})

	t.Run("invalid key errors", func(t *testing.T) {
		var b [64]byte
		_, err := rand.Read(b[:])
		assert.NoError(t, err)

		script, derivationPath, err := bscript.NewP2PKHFromBip32ExtKey(&bip32.ExtendedKey{})

		assert.Error(t, err)
		assert.Empty(t, derivationPath)
		assert.Nil(t, script)
	})
}

func TestNewFromHexString(t *testing.T) {
	t.Parallel()

	s, err := bscript.NewFromHexString("76a914e2a623699e81b291c0327f408fea765d534baa2a88ac")
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t,
		"76a914e2a623699e81b291c0327f408fea765d534baa2a88ac",
		hex.EncodeToString(*s),
	)
}

func TestScript_ToASM(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		script string
		expASM string
	}{
		"valid script": {
			script: "76a914e2a623699e81b291c0327f408fea765d534baa2a88ac",
			expASM: "OP_DUP OP_HASH160 e2a623699e81b291c0327f408fea765d534baa2a OP_EQUALVERIFY OP_CHECKSIG",
		},
		"empty script:": {
			script: "",
			expASM: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s, err := bscript.NewFromHexString(test.script)
			assert.NoError(t, err)

			asm, err := s.ToASM()
			assert.NoError(t, err)

			assert.Equal(t, test.expASM, asm)
		})
	}
}

func TestNewFromASM(t *testing.T) {
	t.Parallel()

	s, err := bscript.NewFromASM("OP_DUP OP_HASH160 e2a623699e81b291c0327f408fea765d534baa2a OP_EQUALVERIFY OP_CHECKSIG")
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t,
		"76a914e2a623699e81b291c0327f408fea765d534baa2a88ac",
		hex.EncodeToString(*s),
	)
}

func TestScript_IsP2PKH(t *testing.T) {
	t.Parallel()

	b, err := hex.DecodeString("76a91403ececf2d12a7f614aef4c82ecf13c303bd9975d88ac")
	assert.NoError(t, err)

	scriptPub := bscript.NewFromBytes(b)
	assert.NotNil(t, scriptPub)
	assert.Equal(t, true, scriptPub.IsP2PKH())
}

func TestScript_IsP2PK(t *testing.T) {
	t.Parallel()

	b, err := hex.DecodeString("2102f0d97c290e79bf2a8660c406aa56b6f189ff79f2245cc5aff82808b58131b4d5ac")
	assert.NoError(t, err)

	scriptPub := bscript.NewFromBytes(b)
	assert.NotNil(t, scriptPub)
	assert.Equal(t, true, scriptPub.IsP2PK())
}

func TestScript_IsP2SH(t *testing.T) {
	t.Parallel()

	b, err := hex.DecodeString("a9149de5aeaff9c48431ba4dd6e8af73d51f38e451cb87")
	assert.NoError(t, err)

	scriptPub := bscript.NewFromBytes(b)
	assert.NotNil(t, scriptPub)
	assert.Equal(t, true, scriptPub.IsP2SH())
}

func TestScript_IsData(t *testing.T) {
	t.Parallel()

	b, err := hex.DecodeString("006a04ac1eed884d53027b2276657273696f6e223a22302e31222c22686569676874223a3634323436302c22707265764d696e65724964223a22303365393264336535633366376264393435646662663438653761393933393362316266623366313166333830616533306432383665376666326165633561323730222c22707265764d696e65724964536967223a2233303435303232313030643736333630653464323133333163613836663031386330343665353763393338663139373735303734373333333533363062653337303438636165316166333032323030626536363034353430323162663934363465393966356139353831613938633963663439353430373539386335396234373334623266646234383262663937222c226d696e65724964223a22303365393264336535633366376264393435646662663438653761393933393362316266623366313166333830616533306432383665376666326165633561323730222c2276637478223a7b2274784964223a2235373962343335393235613930656533396133376265336230306239303631653734633330633832343133663664306132303938653162656137613235313566222c22766f7574223a307d2c226d696e6572436f6e74616374223a7b22656d61696c223a22696e666f407461616c2e636f6d222c226e616d65223a225441414c20446973747269627574656420496e666f726d6174696f6e20546563686e6f6c6f67696573222c226d65726368616e74415049456e64506f696e74223a2268747470733a2f2f6d65726368616e746170692e7461616c2e636f6d2f227d7d46304402206fd1c6d6dd32cc85ddd2f30bc068445dd901c6bd85e394e45bb254716d2bb228022041f0f8b1b33c2e3702aee4ad47155548045ed945738b43dc0faed2e86faa12e4")
	assert.NoError(t, err)

	scriptPub := bscript.NewFromBytes(b)
	assert.NotNil(t, scriptPub)
	assert.Equal(t, true, scriptPub.IsData())
}

func TestScript_IsMultisigOut(t *testing.T) { // TODO: check this
	t.Parallel()

	t.Run("is multisig", func(t *testing.T) {
		b, err := hex.DecodeString("5201110122013353ae")
		assert.NoError(t, err)

		scriptPub := bscript.NewFromBytes(b)
		assert.NotNil(t, scriptPub)
		assert.Equal(t, true, scriptPub.IsMultiSigOut())
	})

	t.Run("is not multisig and no error", func(t *testing.T) {
		//Test Txid:de22e20422dbba8e8eeab87d5532480499abb01d6619bb66fe374f4d4a7500ee, vout:1

		b, err := hex.DecodeString("5101400176018801a901ac615e7961007901687f7700005279517f75007f77007901fd8763615379537f75517f77007901007e81517a7561537a75527a527a5379535479937f75537f77527a75517a67007901fe8763615379557f75517f77007901007e81517a7561537a75527a527a5379555479937f75557f77527a75517a67007901ff8763615379597f75517f77007901007e81517a7561537a75527a527a5379595479937f75597f77527a75517a67615379517f75007f77007901007e81517a7561537a75527a527a5379515479937f75517f77527a75517a6868685179517a75517a75517a75517a7561517a7561007982770079011494527951797f77537952797f750001127900a063610113795a7959797e01147e51797e5a797e58797e517a7561610079011479007958806152790079827700517902fd009f63615179515179517951938000795179827751947f75007f77517a75517a75517a7561517a75675179030000019f6301fd615279525179517951938000795179827751947f75007f77517a75517a75517a75617e517a756751790500000000019f6301fe615279545179517951938000795179827751947f75007f77517a75517a75517a75617e517a75675179090000000000000000019f6301ff615279585179517951938000795179827751947f75007f77517a75517a75517a75617e517a7568686868007953797e517a75517a75517a75617e517a75517a7561527951797e537a75527a527a527975757568607900a06351790112797e610079011279007958806152790079827700517902fd009f63615179515179517951938000795179827751947f75007f77517a75517a75517a7561517a75675179030000019f6301fd615279525179517951938000795179827751947f75007f77517a75517a75517a75617e517a756751790500000000019f6301fe615279545179517951938000795179827751947f75007f77517a75517a75517a75617e517a75675179090000000000000000019f6301ff615279585179517951938000795179827751947f75007f77517a75517a75517a75617e517a7568686868007953797e517a75517a75517a75617e517a75517a7561527951797e537a75527a527a5279757575685e7900a063615f795a7959797e01147e51797e5a797e58797e517a75616100796079007958806152790079827700517902fd009f63615179515179517951938000795179827751947f75007f77517a75517a75517a7561517a75675179030000019f6301fd615279525179517951938000795179827751947f75007f77517a75517a75517a75617e517a756751790500000000019f6301fe615279545179517951938000795179827751947f75007f77517a75517a75517a75617e517a75675179090000000000000000019f6301ff615279585179517951938000795179827751947f75007f77517a75517a75517a75617e517a7568686868007953797e517a75517a75517a75617e517a75517a7561527951797e537a75527a527a5279757575685c7900a063615d795a7959797e01147e51797e5a797e58797e517a75616100795e79007958806152790079827700517902fd009f63615179515179517951938000795179827751947f75007f77517a75517a75517a7561517a75675179030000019f6301fd615279525179517951938000795179827751947f75007f77517a75517a75517a75617e517a756751790500000000019f6301fe615279545179517951938000795179827751947f75007f77517a75517a75517a75617e517a75675179090000000000000000019f6301ff615279585179517951938000795179827751947f75007f77517a75517a75517a75617e517a7568686868007953797e517a75517a75517a75617e517a75517a7561527951797e537a75527a527a5279757575680079aa007961011679007982775179517958947f7551790128947f77517a75517a75618769011679a954798769011779011779ac69610115796100792097dfd76851bf465e8f715593b217714858bbe9570ff3bd5e33840a34e20ff0262102ba79df5f8ae7604a9830f03c7933028186aede0675a16f025dc4f8be8eec0382210ac407f0e4bd44bfc207355a778b046225a7068fc59ee7eda43ad905aadbffc800206c266b30e6a1319c66dc401e5bd6b432ba49688eecd118297041da8074ce0810201008ce7480da41702918d1ec8e6849ba32b4d65b1e40dc669c31a1e6306b266c011379011379855679aa616100790079517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e01007e81517a756157795679567956795679537956795479577995939521414136d08c5ed2bf3ba048afe6dcaebafeffffffffffffffffffffffffffffff0061517951795179517997527a75517a5179009f635179517993527a75517a685179517a75517a7561527a75517a517951795296a0630079527994527a75517a68537982775279827754527993517993013051797e527e53797e57797e527e52797e5579517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7e56797e0079517a75517a75517a75517a75517a75517a75517a75517a75517a75517a75517a75517a756100795779ac517a75517a75517a75517a75517a75517a75517a75517a75517a7561517a75617777777777777777777777777777777777777777777777776ae0cfa0c0930b63270459fe368d5ed31da74c00de")
		assert.NoError(t, err)

		scriptPub := bscript.NewFromBytes(b)
		assert.NotNil(t, scriptPub)
		assert.Equal(t, false, scriptPub.IsMultiSigOut())
	})
}

func TestScript_PublicKeyHash(t *testing.T) {
	t.Parallel()

	t.Run("get as bytes", func(t *testing.T) {
		b, err := hex.DecodeString("76a91404d03f746652cfcb6cb55119ab473a045137d26588ac")
		assert.NoError(t, err)

		s := bscript.NewFromBytes(b)
		assert.NotNil(t, s)

		var pkh []byte
		pkh, err = s.PublicKeyHash()
		assert.NoError(t, err)
		assert.Equal(t, "04d03f746652cfcb6cb55119ab473a045137d265", hex.EncodeToString(pkh))
	})

	t.Run("get as string", func(t *testing.T) {
		s, err := bscript.NewFromHexString("76a91404d03f746652cfcb6cb55119ab473a045137d26588ac")
		assert.NoError(t, err)
		assert.NotNil(t, s)

		var pkh []byte
		pkh, err = s.PublicKeyHash()
		assert.NoError(t, err)
		assert.Equal(t, "04d03f746652cfcb6cb55119ab473a045137d265", hex.EncodeToString(pkh))
	})

	t.Run("empty script", func(t *testing.T) {
		s := &bscript.Script{}

		_, err := s.PublicKeyHash()
		assert.Error(t, err)
		assert.EqualError(t, err, "script is empty")
	})

	t.Run("nonstandard script", func(t *testing.T) {
		// example tx 37d4cc9f8a3b62e7f2e7c97c07a3282bfa924739c0e174733ff1b764ef8e3ebc
		s, err := bscript.NewFromHexString("76")
		assert.NoError(t, err)
		assert.NotNil(t, s)

		_, err = s.PublicKeyHash()
		assert.Error(t, err)
		assert.EqualError(t, err, "not a P2PKH")
	})
}

func TestErrorIsAppended(t *testing.T) {
	script, _ := hex.DecodeString("6a0548656c6c6f0548656c6c")
	s := bscript.Script(script)

	asm, err := s.ToASM()
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(asm, "[error]"), "toASM() should end with [error]")
}

func TestScript_AppendOpcodes(t *testing.T) {
	tests := map[string]struct {
		script    string
		appends   []byte
		expScript string
		expErr    error
	}{
		"successful single append": {
			script:    "OP_2 OP_2 OP_ADD",
			appends:   []byte{bscript.OpEQUALVERIFY},
			expScript: "OP_2 OP_2 OP_ADD OP_EQUALVERIFY",
		},
		"successful multiple append": {
			script:    "OP_2 OP_2 OP_ADD",
			appends:   []byte{bscript.OpEQUAL, bscript.OpVERIFY},
			expScript: "OP_2 OP_2 OP_ADD OP_EQUAL OP_VERIFY",
		},
		"unsuccessful push adata append": {
			script:  "OP_2 OP_2 OP_ADD",
			appends: []byte{bscript.OpEQUAL, bscript.OpPUSHDATA1, 0x44},
			expErr:  bscript.ErrInvalidOpcodeType,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			script, err := bscript.NewFromASM(test.script)
			assert.NoError(t, err)

			err = script.AppendOpcodes(test.appends...)
			if test.expErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, test.expErr, errors.Unwrap(err).Error())
			} else {
				assert.NoError(t, err)
				asm, err := script.ToASM()
				assert.NoError(t, err)
				assert.Equal(t, test.expScript, asm)
			}
		})
	}
}

func TestScript_Equals(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		script1 *bscript.Script
		script2 *bscript.Script
		exp     bool
	}{
		"P2PKH scripts that equal should return true": {
			script1: func() *bscript.Script {
				s, err := bscript.NewP2PKHFromAddress("n2wmGVP89x3DsLNqk3NvctfQy9m9pvt7mk")
				assert.NoError(t, err)
				return s
			}(),
			script2: func() *bscript.Script {
				s, err := bscript.NewP2PKHFromAddress("n2wmGVP89x3DsLNqk3NvctfQy9m9pvt7mk")
				assert.NoError(t, err)
				return s
			}(),
			exp: true,
		}, "scripts from bytes equal should return true": {
			script1: func() *bscript.Script {
				b, err := hex.DecodeString("5201110122013353ae")
				assert.NoError(t, err)

				return bscript.NewFromBytes(b)
			}(),
			script2: func() *bscript.Script {
				b, err := hex.DecodeString("5201110122013353ae")
				assert.NoError(t, err)

				return bscript.NewFromBytes(b)
			}(),
			exp: true,
		}, "scripts from hex, equal should return true": {
			script1: func() *bscript.Script {
				s, err := bscript.NewFromHexString("76a91404d03f746652cfcb6cb55119ab473a045137d26588ac")
				assert.NoError(t, err)
				assert.NotNil(t, s)
				return s
			}(),
			script2: func() *bscript.Script {
				s, err := bscript.NewFromHexString("76a91404d03f746652cfcb6cb55119ab473a045137d26588ac")
				assert.NoError(t, err)
				assert.NotNil(t, s)
				return s
			}(),
			exp: true,
		}, "scripts from hex, not equal should return false": {
			script1: func() *bscript.Script {
				s, err := bscript.NewFromHexString("76a91404d03f746652cfcb6cb55119ab473a045137d26566ac")
				assert.NoError(t, err)
				assert.NotNil(t, s)
				return s
			}(),
			script2: func() *bscript.Script {
				s, err := bscript.NewFromHexString("76a91404d03f746652cfcb6cb55119ab473a045137d26588ac")
				assert.NoError(t, err)
				assert.NotNil(t, s)
				return s
			}(),
			exp: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.exp, test.script1.Equals(test.script2))
			assert.Equal(t, test.exp, test.script1.EqualsBytes(*test.script2))
			assert.Equal(t, test.exp, test.script1.EqualsHex(test.script2.String()))
		})
	}
}

func TestScript_MinPushSize(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		data   [][]byte
		expLen int
	}{
		"OpX / OpNeg returns 1": {
			data: [][]byte{
				{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9},
				{10}, {11}, {12}, {13}, {14}, {15}, {16}, {0x81},
			},
			expLen: 1,
		},
		"OP_DATA_1 + data returns 2": {
			data: [][]byte{
				{0x17}, {0x18}, {0x19}, {0x20}, {0x21}, {0x22}, {0x23}, {0x24}, {0x25}, {0x26},
				{0x27}, {0x28}, {0x29}, {0x30}, {0x31}, {0x32}, {0x33}, {0x34}, {0x35}, {0x36},
				{0x37}, {0x38}, {0x39}, {0x40}, {0x41}, {0x42}, {0x43}, {0x44}, {0x45}, {0x46},
				{0x47}, {0x48}, {0x49}, {0x50}, {0x51}, {0x52}, {0x53}, {0x54}, {0x55}, {0x56},
				{0x57}, {0x58}, {0x59}, {0x60}, {0x61}, {0x62}, {0x63}, {0x64}, {0x65}, {0x66},
				{0x67}, {0x68}, {0x69}, {0x70}, {0x71}, {0x72}, {0x73}, {0x74}, {0x75}, {0x76},
				{0x78}, {0x79}, {0x7a}, {0x7b}, {0x7c}, {0x7d}, {0x7e}, {0x7f}, {0x80},
				{0x82}, {0x83}, {0x84}, {0x85}, {0x86}, {0x87}, {0x88}, {0x89}, {0x8a}, {0x8b},
				{0x8c}, {0x8d}, {0x8e}, {0x8f}, {0x90}, {0x91}, {0x92}, {0x93}, {0x94}, {0x95},
				{0x96}, {0x97}, {0x98}, {0x99}, {0x9a}, {0x9b}, {0x9c}, {0x9d}, {0x9e}, {0x9f},
				{0xa0}, {0xa1}, {0xa2}, {0xa3}, {0xa4}, {0xa5}, {0xa6}, {0xa7}, {0xa8}, {0xa9},
				{0xaa}, {0xab}, {0xac}, {0xad}, {0xae}, {0xaf}, {0xb0}, {0xb1}, {0xb2}, {0xb3},
				{0xb4}, {0xb5}, {0xb6}, {0xb7}, {0xb8}, {0xb9}, {0xba}, {0xbb}, {0xbc}, {0xbd},
				{0xbe}, {0xbf}, {0xc0}, {0xc1}, {0xc2}, {0xc3}, {0xc4}, {0xc5}, {0xc6}, {0xc7},
				{0xc8}, {0xc9}, {0xca}, {0xcb}, {0xcc}, {0xcd}, {0xce}, {0xcf}, {0xd0}, {0xd1},
				{0xd2}, {0xd3}, {0xd4}, {0xd5}, {0xd6}, {0xd7}, {0xd8}, {0xd9}, {0xda}, {0xdb},
				{0xdc}, {0xdd}, {0xde}, {0xdf}, {0xe0}, {0xe1}, {0xe2}, {0xe3}, {0xe4}, {0xe5},
				{0xe6}, {0xe7}, {0xe8}, {0xe9}, {0xea}, {0xeb}, {0xec}, {0xed}, {0xee}, {0xef},
				{0xf0}, {0xf1}, {0xf2}, {0xf3}, {0xf4}, {0xf5}, {0xf6}, {0xf7}, {0xf8}, {0xf9},
				{0xfa}, {0xfb}, {0xfc}, {0xfd}, {0xfe}, {0xff},
			},
			expLen: 2,
		},
		"OP_DATA_2 onward returns len(data)+1": {
			data: [][]byte{func() []byte {
				return bytes.Repeat([]byte{0x00}, 23)
			}()},
			expLen: 23 + 1,
		},
		"OP_DATA_75 returns len(data)+1 (max)": {
			data: [][]byte{func() []byte {
				return bytes.Repeat([]byte{0x00}, 75)
			}()},
			expLen: 75 + 1,
		},
		"OP_PUSHDATA1 + length byte + data returns len(data)+2": {
			data: [][]byte{func() []byte {
				return bytes.Repeat([]byte{0x00}, 86)
			}()},
			expLen: 86 + 2,
		},
		"OP_PUSHDATA1 + length byte + data returns len(data)+2 (max)": {
			data: [][]byte{func() []byte {
				return bytes.Repeat([]byte{0x00}, 255)
			}()},
			expLen: 255 + 2,
		},
		"OP_PUSHDATA2 + length byte + data returns len(data)+3": {
			data: [][]byte{func() []byte {
				return bytes.Repeat([]byte{0x00}, 256)
			}()},
			expLen: 256 + 3,
		},
		"OP_PUSHDATA2 + length byte + data returns len(data)+3 (max)": {
			data: [][]byte{func() []byte {
				return bytes.Repeat([]byte{0x00}, 65535)
			}()},
			expLen: 65535 + 3,
		},
		"OP_PUSHDATA4 + length byte + data returns len(data)+5": {
			data: [][]byte{func() []byte {
				return bytes.Repeat([]byte{0x00}, 65536)
			}()},
			expLen: 65536 + 5,
		},
		// These tests cause the CI to OOM due to the massive slices being created
		//"OP_PUSHDATA4 + length byte + data returns len(data)+5 (max)": {
		//	data: [][]byte{func() []byte {
		//		return bytes.Repeat([]byte{0x00}, 0xffffffff)
		//	}()},
		//	expLen: 0xffffffff + 5,
		//},
		//"data too large returns 0": {
		//	data: [][]byte{func() []byte {
		//		return bytes.Repeat([]byte{0x00}, 0xffffffff+1)
		//	}()},
		//	expLen: 0,
		//},
		"Op0 returns 1": {
			data:   [][]byte{},
			expLen: 1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			for _, data := range test.data {
				assert.Equal(t, test.expLen, bscript.MinPushSize(data), "data: %x", data)
			}
		})
	}
}

func TestScript_MarshalJSON(t *testing.T) {
	script, err := bscript.NewFromASM("OP_2 OP_2 OP_ADD OP_4 OP_EQUALVERIFY")
	assert.NoError(t, err)

	bb, err := json.Marshal(script)
	assert.NoError(t, err)

	assert.Equal(t, `"5252935488"`, string(bb))
}

func TestScript_UnmarshalJSON(t *testing.T) {
	tests := map[string]struct {
		jsonString string
		exp        string
	}{
		"script with content": {
			jsonString: `"5252935488"`,
			exp:        "5252935488",
		},
		"empty script": {
			jsonString: `""`,
			exp:        "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var out *bscript.Script
			assert.NoError(t, json.Unmarshal([]byte(test.jsonString), &out))
			assert.Equal(t, test.exp, out.String())
		})
	}
}

func TestScriptToAsm(t *testing.T) {
	script, _ := hex.DecodeString("006a2231394878696756345179427633744870515663554551797131707a5a56646f4175744d7301e4b8bbe381aa54574954544552e381a8425356e381ae547765746368e381a7e381aee98195e381840a547765746368e381a7e381afe887aae58886e381aee69bb8e38184e3819fe38384e382a8e383bce38388e381afe4b880e795aae69c80e5889de381bee381a70ae38195e3818be381aee381bce381a3e381a6e38184e381a4e381a7e38282e7a2bae8aa8de58fafe883bde381a7e8aaade381bfe8bebce381bfe381a7e995b7e69982e996930ae5819ce6ada2e381afe38182e3828ae381bee3819be38293e380825954e383aae383b3e382afe381aee58b95e794bbe38292e8a696e881b4e38197e3819fe5a0b4e590880ae99fb3e6a5bde381afe382b9e382afe383ade383bce383abe38197e381a6e38282e98094e58887e3828ce3819ae881b4e38193e38188e3819fe381bee381be0ae38384e382a4e38383e382bfe383bce381afe69c80e5889de381aee383ace382b9e381bee381a7e8a18ce38191e381aae38184e381a7e38197e38287e380820a746578742f706c61696e04746578741f7477657463685f7477746578745f313634343834393439353138332e747874017c223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540b7477646174615f6a736f6e046e756c6c0375726c046e756c6c07636f6d6d656e74046e756c6c076d625f757365720439373038057265706c794035366462363536376363306230663539316265363561396135313731663533396635316334333165643837356464326136373431643733353061353539363762047479706504706f73740974696d657374616d70046e756c6c036170700674776574636807696e766f6963652461366637336133312d336334342d346164612d393937352d386537386261666661623765017c22313550636948473232534e4c514a584d6f53556157566937575371633768436676610d424954434f494e5f454344534122314c6970354b335671677743415662674d7842536547434d344355364e344e6b75744c58494b4b554a35765a7753336b4c456e353749356a36485a2b43325733393834314e543532334a4c374534387655706d6f57306b4677613767392b51703246434f4d42776a556a7a76454150624252784d496a746c6b476b3d")
	s := bscript.Script(script)

	asm, err := s.ToASM()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	expected := "0 OP_RETURN 31394878696756345179427633744870515663554551797131707a5a56646f417574 e4b8bbe381aa54574954544552e381a8425356e381ae547765746368e381a7e381aee98195e381840a547765746368e381a7e381afe887aae58886e381aee69bb8e38184e3819fe38384e382a8e383bce38388e381afe4b880e795aae69c80e5889de381bee381a70ae38195e3818be381aee381bce381a3e381a6e38184e381a4e381a7e38282e7a2bae8aa8de58fafe883bde381a7e8aaade381bfe8bebce381bfe381a7e995b7e69982e996930ae5819ce6ada2e381afe38182e3828ae381bee3819be38293e380825954e383aae383b3e382afe381aee58b95e794bbe38292e8a696e881b4e38197e3819fe5a0b4e590880ae99fb3e6a5bde381afe382b9e382afe383ade383bce383abe38197e381a6e38282e98094e58887e3828ce3819ae881b4e38193e38188e3819fe381bee381be0ae38384e382a4e38383e382bfe383bce381afe69c80e5889de381aee383ace382b9e381bee381a7e8a18ce38191e381aae38184e381a7e38197e38287e38082 746578742f706c61696e 1954047348 7477657463685f7477746578745f313634343834393439353138332e747874 124 3150755161374b36324d694b43747373534c4b79316b683536575755374d74555235 5522771 7477646174615f6a736f6e 1819047278 7107189 1819047278 636f6d6d656e74 1819047278 6d625f75736572 942683961 7265706c79 35366462363536376363306230663539316265363561396135313731663533396635316334333165643837356464326136373431643733353061353539363762 1701869940 1953722224 74696d657374616d70 1819047278 7368801 747765746368 696e766f696365 61366637336133312d336334342d346164612d393937352d386537386261666661623765 124 313550636948473232534e4c514a584d6f5355615756693757537163376843667661 424954434f494e5f4543445341 314c6970354b335671677743415662674d7842536547434d344355364e344e6b7574 494b4b554a35765a7753336b4c456e353749356a36485a2b43325733393834314e543532334a4c374534387655706d6f57306b4677613767392b51703246434f4d42776a556a7a76454150624252784d496a746c6b476b3d"

	if asm != expected {
		t.Errorf("\nExpected %q\ngot      %q", expected, asm)
	}

}

func TestRunScriptExample2(t *testing.T) {
	script, _ := hex.DecodeString("006a0372756e01050c63727970746f6669676874734d16057b22696e223a312c22726566223a5b22343561666530303862396634393663333130356663396132636234373234316565643566646531333531303532616339353938323531636666623939376136385f6f31222c22643335343933633964313266656538363134313663333366653336346662336566373531363234373532313833316264623232303933333731303330383663325f6f31222c22306338623636326339363862316537376164626535666161653566666436633033653537353965373833376132353534653438643561356535326335346634385f6f31222c22336136376365633363313662646238343762393732626565326663316330373137633539656463616537626635663438633931666563636661363335616633335f6f31222c22313465323738633638666635323165303931366164376337313361653461303135366537363336316462643362326233353764666236303238653064636137615f6f31222c22613738663561366437326637383731316536366336323131666262643061306266643135616439316264643030343034393238613966616363363364613664395f6f31222c22373535633932326336366363656533353766356265656437383164323631336634313739346230323839333963333435316466653438393032303238343263355f6f31222c22316661333532383030333363343534663465313263323134383333343436643335313734663031666565373064346639653633366664393462363237316436325f6f31222c22636233356534656361336635616334303561636261636464383632346366303835636333626535336639323633663531616565373037393234616265316237385f6f31222c22373166626133383633343162393332333830656335626665646333613430626365343364343937346465636463393463343139613934613863653564666332335f6f31222c22363161653132323165646438626431646438336332326461326232616237643131346139313239363439366365336664306562613737333236623638613238335f6f31222c22386462643166643638373934353131636364616338333938333136393638306662616338356233613961626439636166366361326666343839633862313633385f6f31222c22386564316564633665656439386135326635373234396333353032663266333764623561336666356233643135613930363732353063383465363035366531315f6f31222c22343062396534373865333766383733636532386364383162666635323532346631383063623538353837376331656139383636343933383039363363646237385f6f31225d2c226f7574223a5b2262633038643265323932623036313031323463313337656361356566393632383464363534363139616162366630313461333932353532613066336339666162225d2c2264656c223a5b5d2c22637265223a5b5d2c2265786563223a5b7b226f70223a2243414c4c222c2264617461223a5b7b22246a6967223a307d2c227265736f6c7665222c5b2263633031363466343332613563383635306331393731376430373161656561646261393065643761303939386234396464656535663537316430323139373032222c313634373630333734373833342c302c2230336339386161663266623237393930613130356364336362616462656461383536613064363238343262666564633430353730343966636232343163333563222c5b302c302c305d5d5d7d5d7d")
	s := bscript.Script(script)

	asm, err := s.ToASM()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	expected := "0 OP_RETURN 7239026 5 63727970746f666967687473 7b22696e223a312c22726566223a5b22343561666530303862396634393663333130356663396132636234373234316565643566646531333531303532616339353938323531636666623939376136385f6f31222c22643335343933633964313266656538363134313663333366653336346662336566373531363234373532313833316264623232303933333731303330383663325f6f31222c22306338623636326339363862316537376164626535666161653566666436633033653537353965373833376132353534653438643561356535326335346634385f6f31222c22336136376365633363313662646238343762393732626565326663316330373137633539656463616537626635663438633931666563636661363335616633335f6f31222c22313465323738633638666635323165303931366164376337313361653461303135366537363336316462643362326233353764666236303238653064636137615f6f31222c22613738663561366437326637383731316536366336323131666262643061306266643135616439316264643030343034393238613966616363363364613664395f6f31222c22373535633932326336366363656533353766356265656437383164323631336634313739346230323839333963333435316466653438393032303238343263355f6f31222c22316661333532383030333363343534663465313263323134383333343436643335313734663031666565373064346639653633366664393462363237316436325f6f31222c22636233356534656361336635616334303561636261636464383632346366303835636333626535336639323633663531616565373037393234616265316237385f6f31222c22373166626133383633343162393332333830656335626665646333613430626365343364343937346465636463393463343139613934613863653564666332335f6f31222c22363161653132323165646438626431646438336332326461326232616237643131346139313239363439366365336664306562613737333236623638613238335f6f31222c22386462643166643638373934353131636364616338333938333136393638306662616338356233613961626439636166366361326666343839633862313633385f6f31222c22386564316564633665656439386135326635373234396333353032663266333764623561336666356233643135613930363732353063383465363035366531315f6f31222c22343062396534373865333766383733636532386364383162666635323532346631383063623538353837376331656139383636343933383039363363646237385f6f31225d2c226f7574223a5b2262633038643265323932623036313031323463313337656361356566393632383464363534363139616162366630313461333932353532613066336339666162225d2c2264656c223a5b5d2c22637265223a5b5d2c2265786563223a5b7b226f70223a2243414c4c222c2264617461223a5b7b22246a6967223a307d2c227265736f6c7665222c5b2263633031363466343332613563383635306331393731376430373161656561646261393065643761303939386234396464656535663537316430323139373032222c313634373630333734373833342c302c2230336339386161663266623237393930613130356364336362616462656461383536613064363238343262666564633430353730343966636232343163333563222c5b302c302c305d5d5d7d5d7d"
	if asm != expected {
		t.Errorf("\nExpected %q\ngot      %q", expected, asm)
	}
}

func TestRunScriptExample3(t *testing.T) {
	script, _ := hex.DecodeString("006a223139694733575459537362796f7333754a373333794b347a45696f69314665734e55010042666166383166326364346433663239383061623162363564616166656231656631333561626339643534386461633466366134656361623230653033656365362d300274780134")
	s := bscript.Script(script)

	asm, err := s.ToASM()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	expected := "0 OP_RETURN 3139694733575459537362796f7333754a373333794b347a45696f69314665734e55 0 666166383166326364346433663239383061623162363564616166656231656631333561626339643534386461633466366134656361623230653033656365362d30 30836 52"
	if asm != expected {
		t.Errorf("\nExpected %q\ngot      %q", expected, asm)
	}
}
