package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrencyRate(t *testing.T) {
	rate, err := getCurrencyRate("USD", "2022-10-08")

	assert.NoError(t, err)
	assert.Equal(t, 61.2475, rate)
}

func TestGetCurrencyRate_InvalidCode(t *testing.T) {
	_, err := getCurrencyRate("ABC", "2022-10-08")

	assert.Error(t, err)
}

func TestGetCurrencyRate_InvalidDate(t *testing.T) {
	_, err := getCurrencyRate("USD", "2022-13-01")

	assert.Error(t, err)
}

func TestPrintCurrencyRate(t *testing.T) {
	output := captureOutput(func() {
		getCurrencyRate("USD", "2022-10-08")
	})

	assert.Contains(t, output, "USD (Доллар США): 61,2475")
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}
