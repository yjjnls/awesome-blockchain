// Copyright 2017 ING Bank N.V.
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package byteconversion

import (
	"testing"
	"math/big"
	"reflect"
)


func TestFromBytes127(t *testing.T) {

	bytesIn := []byte{127}
	actualInt, err := FromByteArray(bytesIn)
	expectedInt := big.NewInt(127)

	if err != nil {
		t.Errorf("Unexpected error.")
	}

	if actualInt.Cmp(expectedInt) != 0 {
		t.Errorf("Assert failure: incorrect value: ", actualInt)
	}
}

func TestFromBytesNegLarge(t *testing.T) {

	bytesIn := []byte{228, 20, 131, 38, 208, 100, 246, 105, 110, 11, 247, 198, 
							54, 252, 188, 185, 163, 179, 13, 6, 144, 164, 44, 232, 184, 
							135, 147, 140, 88, 87, 191, 46, 22, 23, 252, 216, 72, 25, 
							5, 124, 29, 81, 56, 242, 199, 0, 68, 132, 102, 246, 34, 
							203, 122, 8, 7, 44, 237, 1, 181, 36}

	actualInt, err := FromByteArray(bytesIn)

	if err != nil {
		t.Errorf("Unexpected error.")
	}

	expectedInt := new(big.Int)
	s := "-34046416216309720507914088123716811131285777464224028073691261328351974060"
	s += "2415131562281745277658128038564338505868351008385659638777486107364060"
	expectedInt.SetString(s, 10)

	if (actualInt.Cmp(expectedInt) != 0) {
		t.Errorf("Assert failure: incorrect value: ", actualInt)
	}
}

func TestToBytesNeg1(t *testing.T) {

	expectedBytes := []byte{255}
	actualBytes := ToByteArray(big.NewInt(-1))

	if !reflect.DeepEqual(expectedBytes, actualBytes) {
		t.Errorf("Assert failure: incorrect byte-array: ", actualBytes)
	}
}

func TestToBytesZero(t *testing.T) {

	expectedBytes := []byte{0}
	actualBytes := ToByteArray(big.NewInt(0))

	if !reflect.DeepEqual(expectedBytes, actualBytes) {
		t.Errorf("Assert failure: incorrect byte-array: ", actualBytes)
	}
}

func TestToBytesPos127(t *testing.T) {

	expectedBytes := []byte{127}
	actualBytes := ToByteArray(big.NewInt(127))

	if !reflect.DeepEqual(expectedBytes, actualBytes) {
		t.Errorf("Assert failure: incorrect byte-array: ", actualBytes)
	}
}


func TestToBytesPos128(t *testing.T) {

	expectedBytes := []byte{0, 128}
	actualBytes := ToByteArray(big.NewInt(128))

	if !reflect.DeepEqual(expectedBytes, actualBytes) {
		t.Errorf("Assert failure: incorrect byte-array: ", actualBytes)
	}
}

func TestToBytesNegLarge(t *testing.T) {

	expectedBytes := []byte{228, 20, 131, 38, 208, 100, 246, 105, 110, 11, 247, 198, 
							54, 252, 188, 185, 163, 179, 13, 6, 144, 164, 44, 232, 184, 
							135, 147, 140, 88, 87, 191, 46, 22, 23, 252, 216, 72, 25, 
							5, 124, 29, 81, 56, 242, 199, 0, 68, 132, 102, 246, 34, 
							203, 122, 8, 7, 44, 237, 1, 181, 36}

	c := new(big.Int)
	s := "-34046416216309720507914088123716811131285777464224028073691261328351974060"
	s += "2415131562281745277658128038564338505868351008385659638777486107364060"

	c.SetString(s, 10)
	actualBytes := ToByteArray(c)

	if !reflect.DeepEqual(expectedBytes, actualBytes) {
		t.Errorf("Assert failure: incorrect byte-array:", actualBytes)
	}
}

func TestToBytesPosLarge(t *testing.T) {

	expectedBytes := []byte{25, 165, 147, 207, 93, 105, 151, 6, 97, 172, 17, 104, 255, 
							15, 255, 106, 246, 80, 123, 85, 231, 140, 213, 98, 148, 126, 
							243, 177, 33, 73, 55, 165, 160, 240, 79, 198, 55, 253, 14, 
							55, 32, 177, 71, 135, 142, 229, 100, 102, 171, 66, 115, 74, 
							137, 135, 74, 20, 127, 132, 155, 6, 126, 202, 57, 173, 72, 
							172, 226, 124, 34, 57, 156, 232, 3, 188, 209, 157, 156, 145, 
							127, 253, 208, 219, 188, 171, 157, 155, 121, 59, 97, 242, 
							121, 187, 52, 222, 168, 83, 10, 89, 57, 83, 109, 245, 248, 
							143, 16, 106, 224, 68, 176, 175, 132, 118, 253, 177, 20, 77, 
							123, 204, 224, 204, 16, 203, 207, 33, 129, 105, 196, 100, 
							140, 179, 5, 167, 73, 146, 98, 210, 60, 247, 253, 95, 19, 
							210, 189, 122, 157, 89, 42, 65, 26, 4, 123, 86, 255, 118, 
							188, 109, 65, 90, 164, 231, 37, 144, 52, 20, 123, 16, 24, 
							18, 139, 147, 149, 145, 241, 82, 242, 163, 254, 236, 26, 
							205, 162, 208, 161, 145, 227, 15, 105, 61, 208, 29, 103, 4, 
							218, 177, 143, 148, 155, 160, 183, 116, 93, 232, 140, 47, 
							48, 61, 167, 130, 135, 160, 67, 69, 13, 156, 78, 212, 45, 
							205, 139, 232, 173, 241, 235, 67, 201, 117, 187, 231, 40, 
							246, 57, 235, 157, 45, 229, 218, 104, 4, 175, 202, 30, 9, 
							118, 237, 41, 227, 44, 60, 33, 29, 66, 125, 181, 117, 249, 
							209, 154, 13, 92, 216, 249, 18, 150, 214, 108, 211, 214, 
							59, 34, 52, 63, 150, 67, 215, 35, 217, 152, 26, 173, 129, 
							242, 101, 184, 80, 194, 17, 176, 100, 24, 211, 54, 2, 159, 
							254, 31, 137, 86, 185, 234, 126, 108, 147, 241, 239, 80, 
							34, 146, 2, 184, 200, 238, 38, 34, 65, 208, 117, 124, 76, 
							120, 41, 18, 66, 174, 207, 163, 163, 58, 70, 37, 1, 214, 
							52, 211, 66, 17, 226, 55, 205, 83, 147, 237, 76, 46, 224, 
							93, 80, 200, 233, 101, 137, 19, 229, 70, 0, 89, 160, 174, 
							33, 120, 9, 151, 214, 187, 219, 171, 155, 233, 110, 236, 
							165, 232, 19, 233, 129, 194, 89, 196, 79, 183, 137, 113, 
							31, 170, 107, 65, 45, 20, 29, 80, 177, 233, 186, 201, 208, 
							153, 137, 123, 18, 161, 116, 55, 163, 29, 81, 160, 9, 86, 
							70, 168, 226, 79, 145, 40, 242, 140, 145, 65, 129, 140, 
							161, 173, 95, 240, 240, 179, 43, 56, 108, 86, 222, 221, 
							121, 246, 154, 181, 96, 231, 173, 252, 87, 40, 46, 46, 226, 
							162, 158, 70, 148, 48, 147, 235, 118, 139, 25, 232, 67, 233, 
							118, 68, 181, 196, 195, 39, 193, 65, 141, 97, 163, 204, 186, 
							57, 124, 96, 89, 123, 30, 120, 210, 196, 89, 114, 43, 35, 39, 
							201, 210, 81, 93, 14, 53, 184, 55, 77, 239, 40, 15, 174, 50, 
							175, 84, 57, 222, 239, 120, 197, 20, 84, 23, 192, 248, 65, 32}

	c := new(big.Int)
	c.SetString("104629761028213599805480925925790425959567049134816086059352755190570739745518666205937735343498567571920655134925379760277594103917074168324328020614304601458655903777672976689775490077688338107235672544860350622774919412889869264967814894367526323662860621314829844532302923598163101691161809563689257119439517772642328497948311622477782095358615210758790515715457810324166800501005841394919183757519983518882638580018812501425864890227616096353048530301822817514776417452329353267096156269934135291712469590753885382455774960306786815780130886957965944086429085979982216274733117683220681240533762036024585652327095570558251184159582343561742487484870623990453783273002077229085268928493917134578545501339771123831252046790356292192853434440371050718893013786884617886032764594150932952606679129922556451438850853950795322470594816924482063623508100143990257521028260138831741855201480669520432545573066469798604777686564518950203058887546628141972672996503838284845932369286957158959203830688307302165544372812778767992596780236523180938001711578036548650984477439658901737951530184210064832275649697658679190072031726344089085533172632374570201218311185378464524281671726449877984821159571188587422006373681926621549415758315808", 10)
	actualBytes := ToByteArray(c)

	if !reflect.DeepEqual(expectedBytes, actualBytes) {
		t.Errorf("Assert failure: incorrect byte-array:", actualBytes)
	}
}

func TestToBytesPos39568(t *testing.T) {

	expectedBytes := []byte{0, 154, 144}
	actualBytes := ToByteArray(big.NewInt(39568))

	if !reflect.DeepEqual(expectedBytes, actualBytes) {
		t.Errorf("Assert failure: incorrect byte-array:", actualBytes)
	}
}
