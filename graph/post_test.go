package graph

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractRawBody(t *testing.T) {
	tests := []struct {
		name string
		body ItemBody
		want string
	}{
		{
			"string type",
			ItemBody{
				"string",
				"test string",
			},
			"test string",
		},
		{
			"unknown type",
			ItemBody{
				"json",
				"test string",
			},
			"test string",
		},
		{
			"html type",
			ItemBody{
				"html",
				`<html><body><div><div><div>test string<tableid="x_jSanity_hideInPlanner"><tr><tdstyle="padding:40px000"><tablewidth="224"border="0"cellspacing="0"cellpadding="0"><tr><tdstyle="background-color:#6264a7;padding:10px30px12px30px;"><divalign="center"style="text-align:center"><ahref="https://planner.cloud.microsoft/test-domain/Home/Task/test-hash-0123456789abcdef?Type=Comment&amp;Channel=GroupMailbox&amp;CreatedTime=639074801954330000"data-auth="NotApplicable"style="text-decoration:none"><spanstyle="color:white;font-size:14px;font-family:SegoeUI,sans-serif,serif,EmojiFont;text-decoration:none;">ReplyinMicrosoftPlanner</span></a></div></td></tr></table></td></tr><tr><tdstyle="padding:15px000"><div><spanstyle="color:#666666;font-size:10px;font-family:SegoeUI,sans-serif,serif,EmojiFont;font-weight:lighter;">Youcanalsoreplytothisemailtoaddataskcomment.</span></div></td></tr><tr><tdstyle="padding:10px000"><div><spanstyle="color:#666666;font-size:10px;font-family:SegoeUI,sans-serif,serif,EmojiFont;font-weight:lighter;">Thistaskisinthe<ahref="https://planner.cloud.microsoft/test-domain/Home/Task/test-hash-0123456789abcdef?Type=Comment&amp;Channel=GroupMailbox&amp;CreatedTime=639074801954330000"data-auth="NotApplicable">mytestplan</a>plan.</span></div></td></tr></table></div></div></div></body></html>`,
			},
			"test string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.body.rawBody())
		})
	}
}
