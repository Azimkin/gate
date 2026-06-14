package chat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/edition/java/proto/nbtconv"
	"go.minekube.com/gate/pkg/edition/java/proto/version"
)

func TestComponentHolderAsComponentAcceptsNBTStyleByteBooleans(t *testing.T) {
	tag, err := nbtconv.SnbtToBinaryTag(`{text:"hi",italic:0B,bold:1B}`)
	require.NoError(t, err)

	holder := &ComponentHolder{
		Protocol:  version.Minecraft_1_21_5.Protocol,
		BinaryTag: tag,
	}
	got, err := holder.AsComponent()
	require.NoError(t, err)

	text, ok := got.(*component.Text)
	require.Truef(t, ok, "got %T", got)
	require.Equal(t, "hi", text.Content)
	require.Equal(t, component.False, text.S.Italic)
	require.Equal(t, component.True, text.S.Bold)
}

func TestComponentHolderAsComponentNormalizesJSONStyleByteBooleans(t *testing.T) {
	holder := &ComponentHolder{
		Protocol: version.Minecraft_1_20_2.Protocol,
		JSON: json.RawMessage(`{
			"text": "hi",
			"obfuscated": "0B",
			"extra": [{
				"text": " child",
				"strikethrough": "1B"
			}]
		}`),
	}

	got, err := holder.AsComponent()
	require.NoError(t, err)

	text, ok := got.(*component.Text)
	require.Truef(t, ok, "got %T", got)
	require.Equal(t, "hi", text.Content)
	require.Equal(t, component.False, text.S.Obfuscated)

	require.Len(t, text.Extra, 1)
	child, ok := text.Extra[0].(*component.Text)
	require.Truef(t, ok, "got %T", text.Extra[0])
	require.Equal(t, " child", child.Content)
	require.Equal(t, component.True, child.S.Strikethrough)
}
