package keyboard

import "github.com/g3n/engine/window"

// decode a key event into string data
func decode(keyEv window.Key, shift bool) string {
	switch keyEv {
	case window.KeyGraveAccent:
		return map[bool]string{false: "`", true: "~"}[shift]
	case window.Key1:
		return map[bool]string{false: "1", true: "!"}[shift]
	case window.Key2:
		return map[bool]string{false: "2", true: "@"}[shift]
	case window.Key3:
		return map[bool]string{false: "3", true: "#"}[shift]
	case window.Key4:
		return map[bool]string{false: "4", true: "$"}[shift]
	case window.Key5:
		return map[bool]string{false: "5", true: "%"}[shift]
	case window.Key6:
		return map[bool]string{false: "6", true: "^"}[shift]
	case window.Key7:
		return map[bool]string{false: "7", true: "&"}[shift]
	case window.Key8:
		return map[bool]string{false: "8", true: "*"}[shift]
	case window.Key9:
		return map[bool]string{false: "9", true: "()"}[shift]
	case window.Key0:
		return map[bool]string{false: "0", true: ")"}[shift]
	case window.KeyMinus:
		return map[bool]string{false: "-", true: "_"}[shift]
	case window.KeyEqual:
		return map[bool]string{false: "=", true: "+"}[shift]
	case window.KeyQ:
		return map[bool]string{false: "q", true: "Q"}[shift]
	case window.KeyW:
		return map[bool]string{false: "w", true: "W"}[shift]
	case window.KeyE:
		return map[bool]string{false: "e", true: "E"}[shift]
	case window.KeyR:
		return map[bool]string{false: "r", true: "R"}[shift]
	case window.KeyT:
		return map[bool]string{false: "t", true: "T"}[shift]
	case window.KeyY:
		return map[bool]string{false: "y", true: "Y"}[shift]
	case window.KeyU:
		return map[bool]string{false: "u", true: "U"}[shift]
	case window.KeyI:
		return map[bool]string{false: "i", true: "I"}[shift]
	case window.KeyO:
		return map[bool]string{false: "o", true: "O"}[shift]
	case window.KeyP:
		return map[bool]string{false: "p", true: "P"}[shift]
	case window.KeyLeftBracket:
		return map[bool]string{false: "[", true: "{"}[shift]
	case window.KeyRightBracket:
		return map[bool]string{false: "]", true: "}"}[shift]
	case window.KeyBackslash:
		return map[bool]string{false: "\\", true: "|"}[shift]
	case window.KeyA:
		return map[bool]string{false: "a", true: "A"}[shift]
	case window.KeyS:
		return map[bool]string{false: "s", true: "S"}[shift]
	case window.KeyD:
		return map[bool]string{false: "d", true: "D"}[shift]
	case window.KeyF:
		return map[bool]string{false: "f", true: "F"}[shift]
	case window.KeyG:
		return map[bool]string{false: "g", true: "G"}[shift]
	case window.KeyH:
		return map[bool]string{false: "h", true: "H"}[shift]
	case window.KeyJ:
		return map[bool]string{false: "j", true: "J"}[shift]
	case window.KeyK:
		return map[bool]string{false: "k", true: "K"}[shift]
	case window.KeyL:
		return map[bool]string{false: "l", true: "L"}[shift]
	case window.KeySemicolon:
		return map[bool]string{false: ";", true: ":"}[shift]
	case window.KeyApostrophe:
		return map[bool]string{false: "'", true: "\""}[shift]
	case window.KeyZ:
		return map[bool]string{false: "z", true: "Z"}[shift]
	case window.KeyX:
		return map[bool]string{false: "x", true: "X"}[shift]
	case window.KeyC:
		return map[bool]string{false: "c", true: "C"}[shift]
	case window.KeyV:
		return map[bool]string{false: "v", true: "V"}[shift]
	case window.KeyB:
		return map[bool]string{false: "b", true: "B"}[shift]
	case window.KeyN:
		return map[bool]string{false: "n", true: "N"}[shift]
	case window.KeyM:
		return map[bool]string{false: "m", true: "M"}[shift]
	case window.KeyComma:
		return map[bool]string{false: ",", true: "<"}[shift]
	case window.KeyPeriod:
		return map[bool]string{false: ".", true: ">"}[shift]
	case window.KeySlash:
		return map[bool]string{false: "/", true: "?"}[shift]
	case window.KeySpace:
		return " "
	case window.KeyKPDivide:
		return "/"
	case window.KeyKPMultiply:
		return "*"
	case window.KeyKPSubtract:
		return "-"
	case window.KeyKPAdd:
		return "+"
	case window.KeyKPDecimal:
		return "."
	case window.KeyKP0:
		return "0"
	case window.KeyKP1:
		return "1"
	case window.KeyKP2:
		return "2"
	case window.KeyKP3:
		return "3"
	case window.KeyKP4:
		return "4"
	case window.KeyKP5:
		return "5"
	case window.KeyKP6:
		return "6"
	case window.KeyKP7:
		return "7"
	case window.KeyKP8:
		return "8"
	case window.KeyKP9:
		return "9"
	}
	return ""
}
