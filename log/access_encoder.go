package log

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/fatih/color"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

const _hex = "0123456789abcdef"

var bufferPool = buffer.NewPool()

var _sliceEncoderPool = sync.Pool{
	New: func() interface{} {
		return &sliceArrayEncoder{elems: make([]interface{}, 0, 2)}
	},
}

func getSliceEncoder() *sliceArrayEncoder {
	return _sliceEncoderPool.Get().(*sliceArrayEncoder)
}

func putSliceEncoder(e *sliceArrayEncoder) {
	e.elems = e.elems[:0]
	_sliceEncoderPool.Put(e)
}

type accessEncoder struct {
	*zapcore.EncoderConfig
	buf *buffer.Buffer
}

func newAccessEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return &accessEncoder{
		EncoderConfig: &cfg,
		buf:           bufferPool.Get(),
	}
}

func (enc *accessEncoder) AddArray(_ string, marshaler zapcore.ArrayMarshaler) error {
	return enc.AppendArray(marshaler)
}

func (enc *accessEncoder) AddObject(_ string, marshaler zapcore.ObjectMarshaler) error {
	return enc.AppendObject(marshaler)
}

func (enc *accessEncoder) AddBinary(key string, value []byte) {
	enc.AddString(key, base64.StdEncoding.EncodeToString(value))
}

func (enc *accessEncoder) AddByteString(_ string, value []byte) {
	enc.AppendByteString(value)
}

func (enc *accessEncoder) AddBool(_ string, value bool) {
	enc.AppendBool(value)
}

func (enc *accessEncoder) AddComplex128(_ string, value complex128) {
	enc.AppendComplex128(value)
}

func (enc *accessEncoder) AddDuration(_ string, value time.Duration) {
	enc.AppendDuration(value)
}

func (enc *accessEncoder) AddFloat64(_ string, value float64) {
	enc.AppendFloat64(value)
}

func (enc *accessEncoder) AddInt64(_ string, value int64) {
	enc.AppendInt64(value)
}

func (enc *accessEncoder) AddString(_, value string) {
	enc.AppendString(value)
}

func (enc *accessEncoder) AddTime(_ string, value time.Time) {
	enc.AppendTime(value)
}

func (enc *accessEncoder) AddUint64(_ string, value uint64) {
	enc.AppendUint64(value)
}

func (enc *accessEncoder) AddReflected(_ string, value interface{}) error {
	marshaled, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = enc.buf.Write(marshaled)
	return err
}

func (enc *accessEncoder) AppendBool(b bool) {
	enc.buf.AppendBool(b)
}

func (enc *accessEncoder) AppendByteString(bytes []byte) {
	enc.safeAddByteString(bytes)
}

func (enc *accessEncoder) AppendComplex128(c complex128) {
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(c)), float64(imag(c))
	enc.buf.AppendByte('"')
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendFloat(r, 64)
	enc.buf.AppendByte('+')
	enc.buf.AppendFloat(i, 64)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte('"')
}

func (enc *accessEncoder) AppendInt64(i int64) {
	enc.buf.AppendInt(i)
}

func (enc *accessEncoder) AppendString(s string) {
	enc.safeAddString(s)
}

func (enc *accessEncoder) AppendUint64(u uint64) {
	enc.buf.AppendUint(u)
}

func (enc *accessEncoder) AppendDuration(duration time.Duration) {
	cur := enc.buf.Len()
	enc.EncodeDuration(duration, enc)
	if cur == enc.buf.Len() {
		// User-supplied EncodeDuration is a no-op. Fall back to nanoseconds to keep
		// JSON valid.
		enc.AppendInt64(int64(duration))
	}
}

func (enc *accessEncoder) AppendTime(t time.Time) {
	cur := enc.buf.Len()
	enc.EncodeTime(t, enc)
	if cur == enc.buf.Len() {
		// User-supplied EncodeTime is a no-op. Fall back to nanos since epoch to keep
		// output JSON valid.
		enc.AppendInt64(t.UnixNano())
	}
}

func (enc *accessEncoder) AppendArray(marshaler zapcore.ArrayMarshaler) error {
	return marshaler.MarshalLogArray(enc)
}

func (enc *accessEncoder) AppendObject(marshaler zapcore.ObjectMarshaler) error {
	return marshaler.MarshalLogObject(enc)
}

func (enc *accessEncoder) AppendReflected(value interface{}) error {
	marshaled, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = enc.buf.Write(marshaled)
	return err
}

func (enc *accessEncoder) AddComplex64(key string, value complex64) {
	enc.AddComplex128(key, complex128(value))
}
func (enc *accessEncoder) AddFloat32(key string, value float32) { enc.AddFloat64(key, float64(value)) }
func (enc *accessEncoder) AddInt(key string, value int)         { enc.AddInt64(key, int64(value)) }
func (enc *accessEncoder) AddInt32(key string, value int32)     { enc.AddInt64(key, int64(value)) }
func (enc *accessEncoder) AddInt16(key string, value int16)     { enc.AddInt64(key, int64(value)) }
func (enc *accessEncoder) AddInt8(key string, value int8)       { enc.AddInt64(key, int64(value)) }
func (enc *accessEncoder) AddUint(key string, value uint)       { enc.AddUint64(key, uint64(value)) }
func (enc *accessEncoder) AddUint32(key string, value uint32)   { enc.AddUint64(key, uint64(value)) }
func (enc *accessEncoder) AddUint16(key string, value uint16)   { enc.AddUint64(key, uint64(value)) }
func (enc *accessEncoder) AddUint8(key string, value uint8)     { enc.AddUint64(key, uint64(value)) }
func (enc *accessEncoder) AddUintptr(key string, value uintptr) { enc.AddUint64(key, uint64(value)) }
func (enc *accessEncoder) AppendComplex64(c complex64)          { enc.AppendComplex128(complex128(c)) }
func (enc *accessEncoder) AppendFloat64(f float64)              { enc.appendFloat(f, 64) }
func (enc *accessEncoder) AppendFloat32(f float32)              { enc.appendFloat(float64(f), 32) }
func (enc *accessEncoder) AppendInt(i int)                      { enc.AppendInt64(int64(i)) }
func (enc *accessEncoder) AppendInt32(i int32)                  { enc.AppendInt64(int64(i)) }
func (enc *accessEncoder) AppendInt16(i int16)                  { enc.AppendInt64(int64(i)) }
func (enc *accessEncoder) AppendInt8(i int8)                    { enc.AppendInt64(int64(i)) }
func (enc *accessEncoder) AppendUint(u uint)                    { enc.AppendUint64(uint64(u)) }
func (enc *accessEncoder) AppendUint32(u uint32)                { enc.AppendUint64(uint64(u)) }
func (enc *accessEncoder) AppendUint16(u uint16)                { enc.AppendUint64(uint64(u)) }
func (enc *accessEncoder) AppendUint8(u uint8)                  { enc.AppendUint64(uint64(u)) }
func (enc *accessEncoder) AppendUintptr(u uintptr)              { enc.AppendUint64(uint64(u)) }
func (enc *accessEncoder) OpenNamespace(_ string)               {}

func (enc *accessEncoder) Clone() zapcore.Encoder {
	clone := enc.clone()
	_, _ = clone.buf.Write(enc.buf.Bytes())
	return clone
}

func (enc *accessEncoder) clone() *accessEncoder {
	clone := &accessEncoder{}
	clone.EncoderConfig = enc.EncoderConfig
	clone.buf = bufferPool.Get()
	return clone
}

func (enc *accessEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	line := bufferPool.Get()

	arr := getSliceEncoder()

	if enc.TimeKey != "" && enc.EncodeTime != nil {
		enc.EncodeTime(entry.Time, arr)
	}
	if enc.LevelKey != "" && enc.EncodeLevel != nil {
		enc.EncodeLevel(entry.Level, arr)
	}
	if entry.LoggerName != "" && enc.NameKey != "" {
		nameEncoder := enc.EncodeName

		if nameEncoder == nil {
			// Fall back to FullNameEncoder for backward compatibility.
			nameEncoder = zapcore.FullNameEncoder
		}

		nameEncoder(entry.LoggerName, arr)
	}

	toRender := make(map[string]interface{})

	for _, field := range fields {
		switch field.Type {
		case zapcore.StringType:
			toRender[field.Key] = field.String
		case zapcore.DurationType:
			fallthrough
		case zapcore.Int64Type:
			toRender[field.Key] = field.Integer
		default:
			toRender[field.Key] = field.Interface
		}
	}

	if status, ok := toRender["status_code"]; ok {
		if code, ok := status.(int64); ok {
			switch category := code / 100; category {
			case 1:
				fallthrough
			case 2:
				arr.AppendString(color.GreenString("%d", code))
			case 3:
				arr.AppendString(color.BlueString("%d", code))
			case 4:
				arr.AppendString(color.YellowString("%d", code))
			default:
				arr.AppendString(color.RedString("%d", code))
			}
		}
	}

	if latency, ok := toRender["latency"]; ok {
		duration := time.Duration(latency.(int64))
		arr.AppendDuration(duration)
	}

	if enc.MessageKey != "" {
		arr.AppendString(entry.Message)
	}

	for i := range arr.elems {
		if i > 0 {
			line.AppendByte('\t')
		}
		_, _ = fmt.Fprint(line, arr.elems[i])
	}

	putSliceEncoder(arr)

	line.AppendString(zapcore.DefaultLineEnding)

	return line, nil
}

func (enc *accessEncoder) addTabIfNecessary(buf *buffer.Buffer) {
	if buf.Len() > 0 {
		buf.AppendByte('\t')
	}
}

// safeAddString JSON-escapes a string and appends it to the internal buffer.
// Unlike the standard library's encoder, it doesn't attempt to protect the
// user from browser vulnerabilities or JSONP-related problems.
func (enc *accessEncoder) safeAddString(s string) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.AppendString(s[i : i+size])
		i += size
	}
}

// safeAddByteString is no-alloc equivalent of safeAddString(string(s)) for s []byte.
func (enc *accessEncoder) safeAddByteString(s []byte) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRune(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		_, _ = enc.buf.Write(s[i : i+size])
		i += size
	}
}

func (enc *accessEncoder) tryAddRuneSelf(b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}
	if 0x20 <= b && b != '\\' && b != '"' {
		enc.buf.AppendByte(b)
		return true
	}
	switch b {
	case '\\', '"':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte(b)
	case '\n':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('n')
	case '\r':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('r')
	case '\t':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('t')
	default:
		// Encode bytes < 0x20, except for the escape sequences above.
		enc.buf.AppendString(`\u00`)
		enc.buf.AppendByte(_hex[b>>4])
		enc.buf.AppendByte(_hex[b&0xF])
	}
	return true
}

func (enc *accessEncoder) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		enc.buf.AppendString(`\ufffd`)
		return true
	}
	return false
}

func (enc *accessEncoder) appendFloat(val float64, bitSize int) {
	switch {
	case math.IsNaN(val):
		enc.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		enc.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		enc.buf.AppendString(`"-Inf"`)
	default:
		enc.buf.AppendFloat(val, bitSize)
	}
}
