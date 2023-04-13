package solanawallet_test

import (
	"fmt"
	"testing"

	"github.com/dmitrymomot/random"
	"github.com/dmitrymomot/solana-wallets/internal/solanawallet"
	"github.com/stretchr/testify/require"
)

func TestEncryptAndDecryptMnemonic(t *testing.T) {
	pin := "123456"
	mnemonic, err := solanawallet.NewMnemonic(solanawallet.MnemonicLength12)
	require.NoError(t, err)
	require.NotEmpty(t, mnemonic)

	fmt.Println("mnemonic:", mnemonic)

	salt := random.String(20)
	client := solanawallet.NewClient(salt)

	encrypted, err := client.EnctyptMnemonic(mnemonic, pin)
	require.NoError(t, err)
	require.NotEmpty(t, encrypted)

	fmt.Println("encrypted mnemonic:", encrypted)

	t.Run("decrypt with correct pin", func(t *testing.T) {
		decrypted, err := client.DecryptMnemonic(encrypted, pin)
		require.NoError(t, err)
		require.NotEmpty(t, decrypted)
		require.Equal(t, mnemonic, decrypted)
	})

	t.Run("decrypt with incorrect pin", func(t *testing.T) {
		_, err := client.DecryptMnemonic(encrypted, "654321")
		require.Error(t, err)
	})

	t.Run("decrypt with incorrect salt", func(t *testing.T) {
		_, err := solanawallet.NewClient(random.String(32)).DecryptMnemonic(encrypted, pin)
		require.Error(t, err)
	})
}
