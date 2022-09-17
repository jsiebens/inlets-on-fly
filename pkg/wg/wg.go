package wg

import (
	"bytes"
	"encoding/base64"
	"errors"
	"net"

	"golang.zx2c4.com/wireguard/device"
)

type Config struct {
	LocalPrivateKey PrivateKey
	LocalNetwork    *net.IPNet

	RemotePublicKey PublicKey
	RemoteNetwork   *net.IPNet

	Endpoint  string
	DNS       net.IP
	KeepAlive int
	MTU       int
	LogLevel  int
}

type PrivateKey device.NoisePrivateKey

type PublicKey device.NoisePublicKey

func (pk *PrivateKey) UnmarshalText(text []byte) error {
	buf, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return err
	}
	if len(buf) != device.NoisePrivateKeySize {
		return errors.New("invalid noise private key")
	}

	copy(pk[:], buf)
	return nil
}

func (pk *PublicKey) UnmarshalText(text []byte) error {
	buf, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return err
	}
	if len(buf) != device.NoisePublicKeySize {
		return errors.New("invalid noise private key")
	}

	copy(pk[:], buf)
	return nil
}

func (pk PrivateKey) ToHex() string {
	const hex = "0123456789abcdef"
	buf := new(bytes.Buffer)
	buf.Reset()
	for i := 0; i < len(pk); i++ {
		buf.WriteByte(hex[pk[i]>>4])
		buf.WriteByte(hex[pk[i]&0xf])
	}
	return buf.String()
}

func (pk PublicKey) ToHex() string {
	const hex = "0123456789abcdef"
	buf := new(bytes.Buffer)
	buf.Reset()
	for i := 0; i < len(pk); i++ {
		buf.WriteByte(hex[pk[i]>>4])
		buf.WriteByte(hex[pk[i]&0xf])
	}
	return buf.String()
}
