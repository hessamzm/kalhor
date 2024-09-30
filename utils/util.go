package utils

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/goravel/framework/facades"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var KlDebug bool

// DONT ETDIT THISS
type ecb struct {
	b         cipher.Block
	blockSize int
}

// پیاده‌سازی Encryptor برای ECB
type ecbEncrypter ecb

// تابع رمزنگاری با الگوریتم DES/ECB/NoPadding
// DONT ETDIT THISS
func encryptPinBlock(key, pinBlock []byte) (string, error) {

	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Pad pinBlock to make sure its length is a multiple of the block size
	pinBlock = PKCS5Padding(pinBlock, block.BlockSize())

	if len(pinBlock)%block.BlockSize() != 0 {
		return "", fmt.Errorf("invalid block size")
	}

	ecb := NewECBEncrypter(block)
	encrypted := make([]byte, len(pinBlock))
	ecb.CryptBlocks(encrypted, pinBlock)

	return BytesToHex(encrypted), nil
}

// ایجاد Encryptor برای ECB
// DONT ETDIT THISS
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(&ecb{
		b:         b,
		blockSize: b.BlockSize(),
	})
}

// پیاده‌سازی روش CryptBlocks برای رمزنگاری با ECB
// DONT ETDIT THISS
func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("ECBEncrypt: input not full blocks")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

// تابع برای تبدیل رشته هگز به بایت‌ها
func HexToBytes(hexStr string) ([]byte, error) {
	return hex.DecodeString(hexStr)
}

// تابع برای تبدیل بایت‌ها به رشته هگز
func BytesToHex(b []byte) string {
	return hex.EncodeToString(b)
}

// تابع برای رمزنگاری Mellat
// DONT EDITE THIS
func MellatEncrypt(s string) string {
	hx := hex.EncodeToString([]byte(s))

	// this key for mellat beh partakht and plz dont change it
	key, err := HexToBytes("2C7D202B960A96AA")
	if err != nil {
		log.Fatal(err)
	}

	pinBlock, err := HexToBytes(hx)
	if err != nil {
		log.Fatal(err)
	}

	// رمزنگاری
	encryptedPan, err := encryptPinBlock(key, pinBlock)
	if err != nil {
		log.Fatal(err)
	}

	if KlDebug {
		fmt.Println("encryptedPan =", encryptedPan)
	}

	return encryptedPan
}

// PKCS5Padding pads the input to a multiple of the block size
// DONT EDITE THIS
func PKCS5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func GetDatamethodget(url string, ErrorArgomans string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		facades.Log().Error(ErrorArgomans, err)
		return nil, err
	}
	defer resp.Body.Close()

	// پردازش پاسخ به عنوان map[string]interface{}
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		facades.Log().Error("Failed to decode response", err)
		return nil, err
	}

	if KlDebug {
		// لاگ گرفتن داده دریافتی
		fmt.Printf("Response data: %+v\n", data)
	}

	return data, nil

	////how to useg
	//USEDATA, ok := data["KEY FROM DATA"].(YOUR TYPE)
	//if !ok {
	//	println("Error converting USEDATA")
	//
	//}
}
func SafeParseFloat(value string) (float64, error) {
	// حذف فاصله‌ها و کاراکترهای غیر ضروری
	value = strings.TrimSpace(value)
	// تبدیل به float64
	return strconv.ParseFloat(value, 64)

	//how to useg
	//USEDATA , err := utils.SafeParseFloat(DATA BY TYPE STRING)
	//USEDATA NOW TYPE FLOAT64
}

func InitDebug() {
	// تعریف فلگ
	flag.BoolVar(&KlDebug, "kldebug", false, "Enable debug mode")
	// پردازش فلگ‌ها
	flag.Parse()
}

// تابع برای لاگ کردن در صورت فعال بودن فلگ debug
func LogDebug(message string) {
	if KlDebug {
		fmt.Println("DEBUG:", message)
	}
}
