package sec

import "testing"

var (
	testKey      = []byte("0102030405060708090a0b0c0d0e0f10")
	testData     = []byte("a data string to be encrypted")
	testPassword = []byte("mySup3Rp@s$w0rd")
)

func TestEncryptDecrypt(t *testing.T) {
	c, err := NewAES(testKey)
	if err != nil {
		t.Errorf("error creating crypter: %v", err)
	}

	data, err := c.Encrypt(testData)
	if err != nil {
		t.Errorf("error encrypting data: %v", err)
	}

	dec, err := c.Decrypt(data)
	if err != nil {
		t.Errorf("error decrypting data: %v", err)
	}

	if string(dec) != string(testData) {
		t.Errorf("decrypted data doesn't match the initial string: '%s'/'%s'",
			string(dec), string(testData))
	}
}

func TestPassKey(t *testing.T) {
	t1, err := CreatePassKey(testPassword)
	if err != nil {
		t.Error(err)
	}

	t2, err := CreatePassKey(testPassword)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < len(t1); i++ {
		if t1[i] != t2[i] {
			t.Errorf("t1 and t2 passkeys don't match")
		}
	}
}

func BenchmarkEncrypt(b *testing.B) {

	c, err := NewAES(testKey)
	if err != nil {
		b.Errorf("error creating crypter: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Encrypt(testData)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	c, err := NewAES(testKey)
	if err != nil {
		b.Errorf("error creating crypter: %v", err)
	}
	data, err := c.Encrypt(testData)
	if err != nil {
		b.Errorf("error encrypting data: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Decrypt(data)
	}
}
