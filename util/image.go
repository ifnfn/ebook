package util

import (
	"bytes"
	"io"
	"io/ioutil"
	"sync"

	"github.com/ifnfn/util/stores"
	"github.com/ifnfn/util/system"
)

// CryptoSave 将文件加密上传
func CryptoSave(key string, file io.Reader, store stores.Store) error {
	var err error
	w := bytes.NewBuffer(nil)

	if err = system.AesEncrpyt(file, w, "zhuzhg_yaodaidi", 1); err == nil {
		err = store.Save(key, w)
	}

	return err
}

// CryptoGet 解密
func CryptoGet(key string, store stores.Store) (io.ReadCloser, error) {
	var err error
	var r io.ReadCloser

	if r, err = store.Get(key); err == nil {
		w := bytes.NewBuffer(nil)
		if err = system.AesDecrypt(r, w, "zhuzhg_yaodaidi", 1); err == nil {
			return ioutil.NopCloser(w), nil
		}
	}

	return nil, err
}

// UploadImage 上传图片，生成大小图
func UploadImage(src []byte, store stores.Store, wg *sync.WaitGroup) string {
	fileName := system.GetMD5(src)[:20]

	go func() {
		if _, err := store.Stat(fileName); err != nil {
			r := bytes.NewBuffer(src)
			if err := CryptoSave(fileName, r, store); err == nil { // 保存大图
				println("uploadImage ok!")
			} else {
				println("big", err.Error())
			}
		}
		if wg != nil {
			wg.Done()
		}

	}()

	return fileName
}
