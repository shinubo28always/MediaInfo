// Unrated Coder t.me/Unrated_Coder

package bot

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"

	"bot/pkg/config"
	"bot/pkg/database"
	"bot/pkg/formatter"
	"bot/pkg/metadata"
	"bot/pkg/streamer"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/html"
	"github.com/gotd/td/telegram/updates"
	updhook "github.com/gotd/td/telegram/updates/hook"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"
)

type Bot struct {
	cfg      *config.Config
	db       *database.Database
	streamer *streamer.Streamer
}

func NewBot(cfg *config.Config, db *database.Database, streamer *streamer.Streamer) *Bot {
	return &Bot{
		cfg:      cfg,
		db:       db,
		streamer: streamer,
	}
}

func generateToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func getMessageID(upds tg.UpdatesClass) int {
	switch v := upds.(type) {
	case *tg.Updates:
		for _, u := range v.Updates {
			if newMessage, ok := u.(*tg.UpdateNewMessage); ok {
				return newMessage.Message.GetID()
			}
			if newChannelMessage, ok := u.(*tg.UpdateNewChannelMessage); ok {
				return newChannelMessage.Message.GetID()
			}
		}
	case *tg.UpdateShortSentMessage:
		return v.ID
	case *tg.UpdateShortMessage:
		return v.ID
	case *tg.UpdateShortChatMessage:
		return v.ID
	case *tg.UpdateShort:
		if newMessage, ok := v.Update.(*tg.UpdateNewMessage); ok {
			return newMessage.Message.GetID()
		}
		if newChannelMessage, ok := v.Update.(*tg.UpdateNewChannelMessage); ok {
			return newChannelMessage.Message.GetID()
		}
	case *tg.UpdatesCombined:
		for _, u := range v.Updates {
			if newMessage, ok := u.(*tg.UpdateNewMessage); ok {
				return newMessage.Message.GetID()
			}
			if newChannelMessage, ok := u.(*tg.UpdateNewChannelMessage); ok {
				return newChannelMessage.Message.GetID()
			}
		}
	}
	return 0
}

func getRepliedUser(ctx context.Context, api *tg.Client, replyToMsgID int) (int64, error) {
	var inputMessages []tg.InputMessageClass
	inputMessages = append(inputMessages, &tg.InputMessageID{ID: replyToMsgID})

	res, err := api.MessagesGetMessages(ctx, inputMessages)
	if err != nil {
		return 0, err
	}

	switch mClass := res.(type) {
	case *tg.MessagesMessages:
		for _, msg := range mClass.Messages {
			if m, ok := msg.(*tg.Message); ok {
				if peerUser, ok := m.FromID.(*tg.PeerUser); ok {
					return peerUser.UserID, nil
				}
			}
		}
	case *tg.MessagesMessagesSlice:
		for _, msg := range mClass.Messages {
			if m, ok := msg.(*tg.Message); ok {
				if peerUser, ok := m.FromID.(*tg.PeerUser); ok {
					return peerUser.UserID, nil
				}
			}
		}
	}
	return 0, fmt.Errorf("could not find replied user")
}

func (b *Bot) Run(ctx context.Context) error {
	dispatcher := tg.NewUpdateDispatcher()
	gaps := updates.New(updates.Config{
		Handler: dispatcher,
	})

	var api *tg.Client
	var sender *message.Sender

	handler := func(ctx context.Context, entities tg.Entities, u message.AnswerableMessageUpdate) error {
		if api == nil || sender == nil {
			return nil
		}
			m, ok := u.GetMessage().(*tg.Message)
			if !ok || m.Out {
				return nil
			}

			var userID int64
			if m.FromID != nil {
				if peerUser, ok := m.FromID.(*tg.PeerUser); ok {
					userID = peerUser.UserID
				}
			}

			if userID == 0 {
				return nil
			}

			banned, err := b.db.IsBanned(ctx, userID)
			if err == nil && banned {
				_, _ = sender.Reply(entities, u).Text(ctx, "❌ **ʏᴏᴜ ᴀʀᴇ ʙᴀɴɴᴇᴅ ꜰʀᴏᴍ ᴜꜱɪɴɢ ᴛʜɪꜱ ʙᴏᴛ.**\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ")
				return nil
			}

			if strings.HasPrefix(m.Message, "/") {
				parts := strings.Fields(m.Message)
				cmd := parts[0]

				if cmd == "/ping" {
					_, _ = sender.Reply(entities, u).Text(ctx, "🏓 **ᴘᴏɴɢ!**\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ")
					return nil
				}

				if cmd == "/start" {
					var username string
					if user, ok := entities.Users[userID]; ok {
						username = user.Username
					}
					_ = b.db.AddUser(ctx, userID, username)

					startText := "👋 **ᴡᴇʟᴄᴏᴍᴇ ᴛᴏ ᴍᴇᴅɪᴀɪɴꜰᴏ ʙᴏᴛ!**\n\n" +
						"ꜱᴇɴᴅ ᴍᴇ ᴀɴʏ ᴠɪᴅᴇᴏ, ᴀᴜᴅɪᴏ, ᴏʀ ᴅᴏᴄᴜᴍᴇɴᴛ (ᴍᴇᴅɪᴀ) ᴛᴏ ɢᴇᴛ ɪᴛꜱ ᴛᴇᴄʜɴɪᴄᴀʟ ᴅᴇᴛᴀɪʟꜱ ᴡɪᴛʜᴏᴜᴛ ᴅᴏᴡɴʟᴏᴀᴅɪɴɢ ᴛʜᴇ ꜰᴜʟʟ ꜰɪʟᴇ.\n\n" +
						"ᴘᴏᴡᴇʀᴇᴅ ʙʏ `ᴍᴇᴅɪᴀɪɴꜰᴏ` ꜱᴛʀᴇᴀᴍɪɴɢ ᴛᴇᴄʜɴᴏʟᴏɢʏ.\n\n" +
						"ᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ"

					startMarkup := &tg.ReplyInlineMarkup{
						Rows: []tg.KeyboardButtonRow{
							{
								Buttons: []tg.KeyboardButtonClass{
									&tg.KeyboardButtonURL{
										Text: "🟢 ᴊᴏɪɴ ᴜᴘᴅᴀᴛᴇ",
										URL:  "https://t.me/Unrated_Coder",
									},
								},
							},
							{
								Buttons: []tg.KeyboardButtonClass{
									&tg.KeyboardButtonCallback{
										Text: "🔵 ʜᴇʟᴘ",
										Data: []byte("help"),
									},
									&tg.KeyboardButtonCallback{
										Text: "🔵 ᴀʙᴏᴜᴛ",
										Data: []byte("about"),
									},
								},
							},
						},
					}

					_, _ = sender.Reply(entities, u).Markup(startMarkup).Text(ctx, startText)
					return nil
				}

				if cmd == "/help" {
					helpText := "ʜᴇʟᴘ ɪɴꜱᴛʀᴜᴄᴛɪᴏɴꜱ:\n" +
						"ꜱᴇɴᴅ ᴀɴʏ ᴠɪᴅᴇᴏ, ᴀᴜᴅɪᴏ, ᴏʀ ᴅᴏᴄᴜᴍᴇɴᴛ (ᴍᴇᴅɪᴀ) ᴅɪʀᴇᴄᴛʟʏ ᴛᴏ ᴛʜᴇ ʙᴏᴛ.\n" +
						"ᴛʜᴇ ʙᴏᴛ ᴡɪʟʟ ɪɴꜱᴛᴀɴᴛʟʏ ᴘʀᴏʙᴇ ᴛʜᴇ ᴍᴇᴅɪᴀ ʜᴇᴀᴅᴇʀ ᴀɴᴅ ᴘʀᴏᴠɪᴅᴇ ᴀ ᴅᴇᴛᴀɪʟᴇᴅ ᴍᴇᴅɪᴀɪɴꜰᴏ ʀᴇᴘᴏʀᴛ.\n\n" +
						"ᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ"

					backMarkup := &tg.ReplyInlineMarkup{
						Rows: []tg.KeyboardButtonRow{
							{
								Buttons: []tg.KeyboardButtonClass{
									&tg.KeyboardButtonCallback{
										Text: "🔙 ʙᴀᴄᴋ",
										Data: []byte("back"),
									},
								},
							},
						},
					}

					_, _ = sender.Reply(entities, u).Markup(backMarkup).Text(ctx, helpText)
					return nil
				}

				if cmd == "/about" {
					aboutText := "ᴀʙᴏᴜᴛ ᴛʜɪꜱ ʙᴏᴛ:\n" +
						"ᴛʜɪꜱ ɪꜱ ᴀ ʜɪɢʜ-ᴘᴇʀꜰᴏʀᴍᴀɴᴄᴇ ᴛᴇʟᴇɢʀᴀᴍ ᴍᴇᴅɪᴀɪɴꜰᴏ ʙᴏᴛ ᴡʀɪᴛᴛᴇɴ ɪɴ ɢᴏ.\n" +
						"ɪᴛ ᴜꜱᴇꜱ ᴍᴛᴘʀᴏᴛᴏ ᴀɴᴅ ʀᴀɴɢᴇ-ʙᴀꜱᴇᴅ ꜱᴛʀᴇᴀᴍɪɴɢ ᴛᴏ ᴇxᴛʀᴀᴄᴛ ᴍᴇᴅɪᴀɪɴꜰᴏ ɪɴ ꜱᴜʙ-ꜱᴇᴄᴏɴᴅꜱ.\n\n" +
						"ᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ"

					backMarkup := &tg.ReplyInlineMarkup{
						Rows: []tg.KeyboardButtonRow{
							{
								Buttons: []tg.KeyboardButtonClass{
									&tg.KeyboardButtonCallback{
										Text: "🔙 ʙᴀᴄᴋ",
										Data: []byte("back"),
									},
								},
							},
						},
					}

					_, _ = sender.Reply(entities, u).Markup(backMarkup).Text(ctx, aboutText)
					return nil
				}

				if cmd == "/log" || cmd == "/logs" {
					isAdmin, err := b.db.IsAdmin(ctx, userID, b.cfg.OwnerID)
					if err == nil && isAdmin {
						logs := GlobalLogCapturer.GetLogs()

						up := uploader.NewUploader(api)
						file, err := up.FromBytes(ctx, "Log.txt", logs)
						if err == nil {
							_, _ = sender.Reply(entities, u).Media(ctx, message.UploadedDocument(file))
						} else {
							_, _ = sender.Reply(entities, u).Text(ctx, fmt.Sprintf("❌ **ꜰᴀɪʟᴇᴅ ᴛᴏ ᴜᴘʟᴏᴀᴅ ʟᴏɢꜱ**: `%v`\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ", err))
						}
					} else {
						_, _ = sender.Reply(entities, u).Text(ctx, "❌ **ᴏɴʟʏ ᴀᴅᴍɪɴꜱ ᴄᴀɴ ᴀᴄᴄᴇꜱꜱ ʟᴏɢꜱ.**\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ")
					}
					return nil
				}

				if cmd == "/users" {
					isAdmin, err := b.db.IsAdmin(ctx, userID, b.cfg.OwnerID)
					if err == nil && isAdmin {
						count, _ := b.db.GetAllUsersCount(ctx)
						_, _ = sender.Reply(entities, u).Text(ctx, fmt.Sprintf("📊 **ᴛᴏᴛᴀʟ ᴜꜱᴇʀꜱ**: `%d`\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ", count))
					}
					return nil
				}

				if cmd == "/remadmin" {
					if userID == b.cfg.OwnerID {
						var userToRemove int64
						header, hasReply := m.GetReplyTo()
						if hasReply {
							if rHeader, ok := header.(*tg.MessageReplyHeader); ok {
								repID, _ := rHeader.GetReplyToMsgID()
								uToRemove, err := getRepliedUser(ctx, api, repID)
								if err == nil {
									userToRemove = uToRemove
								}
							}
						} else if len(parts) > 1 {
							idVal, err := strconv.ParseInt(parts[1], 10, 64)
							if err == nil {
								userToRemove = idVal
							}
						}

						if userToRemove == 0 {
							_, _ = sender.Reply(entities, u).Text(ctx, "❌ ʀᴇᴘʟʏ ᴛᴏ ᴀ ᴜꜱᴇʀ ᴏʀ ᴘʀᴏᴠɪᴅᴇ ᴀ ᴜꜱᴇʀ ɪᴅ ᴛᴏ ʀᴇᴍᴏᴠᴇ.\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ")
							return nil
						}

						_ = b.db.RemoveAdmin(ctx, userToRemove)
						_, _ = sender.Reply(entities, u).Text(ctx, fmt.Sprintf("✅ **ᴀᴅᴍɪɴ ʀᴇᴍᴏᴠᴇᴅ**: `%d`\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ", userToRemove))
					} else {
						_, _ = sender.Reply(entities, u).Text(ctx, "❌ ᴏɴʟʏ ᴛʜᴇ ᴏᴡɴᴇʀ ᴄᴀɴ ʀᴇᴍᴏᴠᴇ ᴀᴅᴍɪɴꜱ.\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ")
					}
					return nil
				}

				if cmd == "/admins" {
					isAdmin, err := b.db.IsAdmin(ctx, userID, b.cfg.OwnerID)
					if err == nil && isAdmin {
						admins, err := b.db.GetAllAdmins(ctx)
						if err != nil {
							_, _ = sender.Reply(entities, u).Text(ctx, "❌ ꜰᴀɪʟᴇᴅ ᴛᴏ ꜰᴇᴛᴄʜ ᴀᴅᴍɪɴꜱ.\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ")
							return nil
						}

						var rows []tg.KeyboardButtonRow
						rows = append(rows, tg.KeyboardButtonRow{
							Buttons: []tg.KeyboardButtonClass{
								&tg.KeyboardButtonCallback{
									Text: fmt.Sprintf("🟢 ᴛᴏᴛᴀʟ ᴀᴅᴍɪɴꜱ: %d", len(admins)),
									Data: []byte("admin_stats"),
								},
							},
						})

						for _, adminID := range admins {
							rows = append(rows, tg.KeyboardButtonRow{
								Buttons: []tg.KeyboardButtonClass{
									&tg.KeyboardButtonCallback{
										Text: fmt.Sprintf("👤 %d", adminID),
										Data: []byte(fmt.Sprintf("admin_id_%d", adminID)),
									},
									&tg.KeyboardButtonCallback{
										Text: "🔴 ʀᴇᴍᴏᴠᴇ",
										Data: []byte(fmt.Sprintf("remadmin_%d", adminID)),
									},
								},
							})
						}

						rows = append(rows, tg.KeyboardButtonRow{
							Buttons: []tg.KeyboardButtonClass{
								&tg.KeyboardButtonCallback{
									Text: "🔄 ʀᴇꜰʀᴇꜱʜ",
									Data: []byte("refresh_admins"),
								},
							},
						})

						markup := &tg.ReplyInlineMarkup{Rows: rows}
						_, _ = sender.Reply(entities, u).Markup(markup).Text(ctx, "👑 **ᴀᴅᴍɪɴɪꜱᴛʀᴀᴛᴏʀꜱ ᴅᴀꜱʜʙᴏᴀʀᴅ**\n\nᴍᴀɴᴀɢᴇ ʙᴏᴛ ᴀᴅᴍɪɴɪꜱᴛʀᴀᴛᴏʀꜱ ᴅɪʀᴇᴄᴛʟʏ ꜰʀᴏᴍ ᴛʜᴇ ɪɴʟɪɴᴇ ᴍᴇɴᴜ ʙᴇʟᴏᴡ.\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ")
					} else {
						_, _ = sender.Reply(entities, u).Text(ctx, "❌ ᴏɴʟʏ ᴀᴅᴍɪɴꜱ ᴄᴀɴ ᴠɪᴇᴡ ᴛʜɪꜱ ᴅᴀꜱʜʙᴏᴀʀᴅ.\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ")
					}
					return nil
				}

				if cmd == "/ban" {
					isAdmin, err := b.db.IsAdmin(ctx, userID, b.cfg.OwnerID)
					if err == nil && isAdmin {
						var userToBan int64
						header, hasReply := m.GetReplyTo()
						if hasReply {
							if rHeader, ok := header.(*tg.MessageReplyHeader); ok {
								repID, _ := rHeader.GetReplyToMsgID()
								uToBan, err := getRepliedUser(ctx, api, repID)
								if err == nil {
									userToBan = uToBan
								}
							}
						} else if len(parts) > 1 {
							idVal, err := strconv.ParseInt(parts[1], 10, 64)
							if err == nil {
								userToBan = idVal
							}
						}

						if userToBan == 0 {
							_, _ = sender.Reply(entities, u).Text(ctx, "❌ ʀᴇᴘʟʏ ᴛᴏ ᴀ ᴜꜱᴇʀ ᴏʀ ᴘʀᴏᴠɪᴅᴇ ᴀ ᴜꜱᴇʀ ɪᴅ ᴛᴏ ʙᴀɴ.\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ")
							return nil
						}

						_ = b.db.BanUser(ctx, userToBan)
						_, _ = sender.Reply(entities, u).Text(ctx, fmt.Sprintf("🚫 **ᴜꜱᴇʀ ʙᴀɴɴᴇᴅ**: `%d`\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ", userToBan))
					}
					return nil
				}

				if cmd == "/unban" {
					isAdmin, err := b.db.IsAdmin(ctx, userID, b.cfg.OwnerID)
					if err == nil && isAdmin {
						var userToUnban int64
						if len(parts) > 1 {
							idVal, err := strconv.ParseInt(parts[1], 10, 64)
							if err == nil {
								userToUnban = idVal
							}
						}

						if userToUnban == 0 {
							_, _ = sender.Reply(entities, u).Text(ctx, "❌ ᴘʟᴇᴀꜱᴇ ᴘʀᴏᴠɪᴅᴇ ᴀ ᴜꜱᴇʀ ɪᴅ ᴛᴏ ᴜɴʙᴀɴ.\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ")
							return nil
						}

						_ = b.db.UnbanUser(ctx, userToUnban)
						_, _ = sender.Reply(entities, u).Text(ctx, fmt.Sprintf("✅ **ᴜꜱᴇʀ ᴜɴʙᴀɴɴᴇᴅ**: `%d`\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ", userToUnban))
					}
					return nil
				}

				if cmd == "/add_admin" {
					if userID == b.cfg.OwnerID {
						var userToAdd int64
						header, hasReply := m.GetReplyTo()
						if hasReply {
							if rHeader, ok := header.(*tg.MessageReplyHeader); ok {
								repID, _ := rHeader.GetReplyToMsgID()
								uToAdd, err := getRepliedUser(ctx, api, repID)
								if err == nil {
									userToAdd = uToAdd
								}
							}
						} else if len(parts) > 1 {
							idVal, err := strconv.ParseInt(parts[1], 10, 64)
							if err == nil {
								userToAdd = idVal
							}
						}

						if userToAdd == 0 {
							_, _ = sender.Reply(entities, u).Text(ctx, "❌ ʀᴇᴘʟʏ ᴛᴏ ᴀ ᴜꜱᴇʀ ᴏʀ ᴘʀᴏᴠɪᴅᴇ ᴀ ᴜꜱᴇʀ ɪᴅ.\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ")
							return nil
						}

						_ = b.db.AddAdmin(ctx, userToAdd)
						_, _ = sender.Reply(entities, u).Text(ctx, fmt.Sprintf("👑 **ᴀᴅᴍɪɴ ᴀᴅᴅᴇᴅ**: `%d`\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ", userToAdd))
					} else {
						_, _ = sender.Reply(entities, u).Text(ctx, "❌ ᴏɴʟʏ ᴛʜᴇ ᴏᴡɴᴇʀ ᴄᴀɴ ᴀᴅᴅ ᴀᴅᴍɪɴꜱ.\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ")
					}
					return nil
				}
			}

			if m.Media != nil {
				if mediaDoc, ok := m.Media.(*tg.MessageMediaDocument); ok {
					doc, ok := mediaDoc.Document.(*tg.Document)
					if ok {
						upds, err := sender.Reply(entities, u).Text(ctx, "🔍 **ᴇxᴛʀᴀᴄᴛɪɴɢ ᴍᴇᴅɪᴀɪɴꜰᴏ...**\n*ᴛʜɪꜱ ᴍᴀʏ ᴛᴀᴋᴇ ᴀ ꜰᴇᴡ ꜱᴇᴄᴏɴᴅꜱ ᴀꜱ ɪ ᴘʀᴏʙᴇ ᴛʜᴇ ʜᴇᴀᴅᴇʀ.*")
						if err != nil {
							log.Printf("Error sending status message: %v", err)
							return nil
						}
						msgID := getMessageID(upds)

						token := generateToken()
						b.streamer.RegisterDocument(token, doc)

						url := fmt.Sprintf("http://127.0.0.1:8080/stream/%s", token)
						data, err := metadata.ExtractMediaInfo(ctx, url)
						if err != nil {
							if msgID != 0 {
								_, _ = sender.Answer(entities, u).Edit(msgID).Text(ctx, fmt.Sprintf("❌ **ᴇʀʀᴏʀ**: `%v`", err))
							}
							return nil
						}

						text := formatter.FormatOutput(data)

						var fileName string = "mediainfo"
						for _, attr := range doc.Attributes {
							if fnAttr, ok := attr.(*tg.DocumentAttributeFilename); ok {
								fileName = fnAttr.FileName
								break
							}
						}

						fullText := fmt.Sprintf("📊 **ᴍᴇᴅɪᴀɪɴꜰᴏ ꜰᴏʀ**: `%s`\n\n%s\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ", fileName, text)

						if len(fullText) > 4096 {
							if msgID != 0 {
								_, _ = sender.Answer(entities, u).Edit(msgID).Text(ctx, fmt.Sprintf("📊 **ᴍᴇᴅɪᴀɪɴꜰᴏ ꜰᴏʀ**: `%s` ɪꜱ ᴛᴏᴏ ʟᴏɴɢ ᴛᴏ ᴅɪꜱᴘʟᴀʏ ɪɴʟɪɴᴇ. ꜱᴇɴᴅɪɴɢ ᴅᴇᴛᴀɪʟᴇᴅ ʀᴇᴘᴏʀᴛ ᴀꜱ ᴀ ᴛᴇxᴛ ꜰɪʟᴇ...", fileName))
							}

							cleanText := strings.ReplaceAll(text, "```", "")
							cleanText = strings.TrimSpace(cleanText)

							up := uploader.NewUploader(api)
							file, err := up.FromBytes(ctx, fileName+".txt", []byte(cleanText))
							if err == nil {
								_, _ = sender.Reply(entities, u).Media(ctx, message.UploadedDocument(file))
							}
						} else {
							if msgID != 0 {
								_, _ = sender.Answer(entities, u).Edit(msgID).Text(ctx, fullText)
							}
						}
					}
				}
			}

			return nil
		}

	dispatcher.OnNewMessage(func(msgCtx context.Context, entities tg.Entities, u *tg.UpdateNewMessage) error {
		go func() {
			_ = handler(ctx, entities, u)
		}()
		return nil
	})

	dispatcher.OnNewChannelMessage(func(chanCtx context.Context, entities tg.Entities, u *tg.UpdateNewChannelMessage) error {
		go func() {
			_ = handler(ctx, entities, u)
		}()
		return nil
	})

	dispatcher.OnBotCallbackQuery(func(cbCtx context.Context, entities tg.Entities, u *tg.UpdateBotCallbackQuery) error {
		go func() {
			_ = func() error {
				if api == nil {
					return nil
				}
				data, ok := u.GetData()
		if !ok {
			return nil
		}

		cmd := string(data)
		var inputPeer tg.InputPeerClass
		switch p := u.Peer.(type) {
		case *tg.PeerUser:
			if user, ok := entities.Users[p.UserID]; ok {
				inputPeer = &tg.InputPeerUser{
					UserID:     user.ID,
					AccessHash: user.AccessHash,
				}
			} else {
				inputPeer = &tg.InputPeerUser{
					UserID: p.UserID,
				}
			}
		case *tg.PeerChat:
			inputPeer = &tg.InputPeerChat{
				ChatID: p.ChatID,
			}
		case *tg.PeerChannel:
			if channel, ok := entities.Channels[p.ChannelID]; ok {
				inputPeer = &tg.InputPeerChannel{
					ChannelID:  channel.ID,
					AccessHash: channel.AccessHash,
				}
			} else {
				inputPeer = &tg.InputPeerChannel{
					ChannelID: p.ChannelID,
				}
			}
		}

		if inputPeer == nil {
			return nil
		}

		startText := "👋 **ᴡᴇʟᴄᴏᴍᴇ ᴛᴏ ᴍᴇᴅɪᴀɪɴꜰᴏ ʙᴏᴛ!**\n\n" +
			"ꜱᴇɴᴅ ᴍᴇ ᴀɴʏ ᴠɪᴅᴇᴏ, ᴀᴜᴅɪᴏ, ᴏʀ ᴅᴏᴄᴜᴍᴇɴᴛ (ᴍᴇᴅɪᴀ) ᴛᴏ ɢᴇᴛ ɪᴛꜱ ᴛᴇᴄʜɴɪᴄᴀʟ ᴅᴇᴛᴀɪʟꜱ ᴡɪᴛʜᴏᴜᴛ ᴅᴏᴡɴʟᴏᴀᴅɪɴɢ ᴛʜᴇ ꜰᴜʟʟ ꜰɪʟᴇ.\n\n" +
			"ᴘᴏᴡᴇʀᴇᴅ ʙʏ `ᴍᴇᴅɪᴀɪɴꜰᴏ` ꜱᴛʀᴇᴀᴍɪɴɢ ᴛᴇᴄʜɴᴏʟᴏɢʏ.\n\n" +
			"ᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ"

		startMarkup := &tg.ReplyInlineMarkup{
			Rows: []tg.KeyboardButtonRow{
				{
					Buttons: []tg.KeyboardButtonClass{
						&tg.KeyboardButtonURL{
							Text: "🟢 ᴊᴏɪɴ ᴜᴘᴅᴀᴛᴇ",
							URL:  "https://t.me/Unrated_Coder",
						},
					},
				},
				{
					Buttons: []tg.KeyboardButtonClass{
						&tg.KeyboardButtonCallback{
							Text: "🔵 ʜᴇʟᴘ",
							Data: []byte("help"),
						},
						&tg.KeyboardButtonCallback{
							Text: "🔵 ᴀʙᴏᴜᴛ",
							Data: []byte("about"),
						},
					},
				},
			},
		}

		backMarkup := &tg.ReplyInlineMarkup{
			Rows: []tg.KeyboardButtonRow{
				{
					Buttons: []tg.KeyboardButtonClass{
						&tg.KeyboardButtonCallback{
							Text: "🔙 ʙᴀᴄᴋ",
							Data: []byte("back"),
						},
					},
				},
			},
		}

		if cmd == "help" {
			helpText := "ʜᴇʟᴘ ɪɴꜱᴛʀᴜᴄᴛɪᴏɴꜱ:\n" +
				"ꜱᴇɴᴅ ᴀɴʏ ᴠɪᴅᴇᴏ, ᴀᴜᴅɪᴏ, ᴏʀ ᴅᴏᴄᴜᴍᴇɴᴛ (ᴍᴇᴅɪᴀ) ᴅɪʀᴇᴄᴛʟʏ ᴛᴏ ᴛʜᴇ ʙᴏᴛ.\n" +
				"ᴛʜᴇ ʙᴏᴛ ᴡɪʟʟ ɪɴꜱᴛᴀɴᴛʟʏ ᴘʀᴏʙᴇ ᴛʜᴇ ᴍᴇᴅɪᴀ ʜᴇᴀᴅᴇʀ ᴀɴᴅ ᴘʀᴏᴠɪᴅᴇ ᴀ ᴅᴇᴛᴀɪʟᴇᴅ ᴍᴇᴅɪᴀɪɴꜰᴏ ʀᴇᴘᴏʀᴛ.\n\n" +
				"ᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ"

			req := &tg.MessagesEditMessageRequest{
				Peer:    inputPeer,
				ID:      u.MsgID,
				Message: helpText,
			}
			req.SetReplyMarkup(backMarkup)
			_, _ = api.MessagesEditMessage(ctx, req)

			_, _ = api.MessagesSetBotCallbackAnswer(ctx, &tg.MessagesSetBotCallbackAnswerRequest{
				QueryID: u.QueryID,
			})
			return nil
		}

		if cmd == "about" {
			aboutText := "ᴀʙᴏᴜᴛ ᴛʜɪꜱ ʙᴏᴛ:\n" +
				"ᴛʜɪꜱ ɪꜱ ᴀ ʜɪɢʜ-ᴘᴇʀꜰᴏʀᴍᴀɴᴄᴇ ᴛᴇʟᴇɢʀᴀᴍ ᴍᴇᴅɪᴀɪɴꜰᴏ ʙᴏᴛ ᴡʀɪᴛᴛᴇɴ ɪɴ ɢᴏ.\n" +
				"ɪᴛ ᴜꜱᴇꜱ ᴍᴛᴘʀᴏᴛᴏ ᴀɴᴅ ʀᴀɴɢᴇ-ʙᴀꜱᴇᴅ ꜱᴛʀᴇᴀᴍɪɴɢ ᴛᴏ ᴇxᴛʀᴀᴄᴛ ᴍᴇᴅɪᴀɪɴꜰᴏ ɪɴ ꜱᴜʙ-ꜱᴇᴄᴏɴᴅꜱ.\n\n" +
				"ᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ"

			req := &tg.MessagesEditMessageRequest{
				Peer:    inputPeer,
				ID:      u.MsgID,
				Message: aboutText,
			}
			req.SetReplyMarkup(backMarkup)
			_, _ = api.MessagesEditMessage(ctx, req)

			_, _ = api.MessagesSetBotCallbackAnswer(ctx, &tg.MessagesSetBotCallbackAnswerRequest{
				QueryID: u.QueryID,
			})
			return nil
		}

		if cmd == "back" {
			req := &tg.MessagesEditMessageRequest{
				Peer:    inputPeer,
				ID:      u.MsgID,
				Message: startText,
			}
			req.SetReplyMarkup(startMarkup)
			_, _ = api.MessagesEditMessage(ctx, req)

			_, _ = api.MessagesSetBotCallbackAnswer(ctx, &tg.MessagesSetBotCallbackAnswerRequest{
				QueryID: u.QueryID,
			})
			return nil
		}

		if cmd == "admin_stats" {
			_, _ = api.MessagesSetBotCallbackAnswer(ctx, &tg.MessagesSetBotCallbackAnswerRequest{
				QueryID: u.QueryID,
				Message: "📊 ᴀᴅᴍɪɴ ꜱᴛᴀᴛɪꜱᴛɪᴄꜱ",
			})
			return nil
		}

		if strings.HasPrefix(cmd, "admin_id_") {
			_, _ = api.MessagesSetBotCallbackAnswer(ctx, &tg.MessagesSetBotCallbackAnswerRequest{
				QueryID: u.QueryID,
				Message: "👤 ᴀᴅᴍɪɴ ᴜꜱᴇʀ ɪᴅ",
			})
			return nil
		}

		if strings.HasPrefix(cmd, "remadmin_") {
			if u.UserID != b.cfg.OwnerID {
				_, _ = api.MessagesSetBotCallbackAnswer(ctx, &tg.MessagesSetBotCallbackAnswerRequest{
					QueryID: u.QueryID,
					Message: "❌ ᴏɴʟʏ ᴛʜᴇ ᴏᴡɴᴇʀ ᴄᴀɴ ʀᴇᴍᴏᴠᴇ ᴀᴅᴍɪɴꜱ.",
					Alert:   true,
				})
				return nil
			}

			targetStr := strings.TrimPrefix(cmd, "remadmin_")
			targetID, err := strconv.ParseInt(targetStr, 10, 64)
			if err == nil {
				_ = b.db.RemoveAdmin(ctx, targetID)
				_, _ = api.MessagesSetBotCallbackAnswer(ctx, &tg.MessagesSetBotCallbackAnswerRequest{
					QueryID: u.QueryID,
					Message: fmt.Sprintf("✅ ᴀᴅᴍɪɴ %d ʀᴇᴍᴏᴠᴇᴅ!", targetID),
					Alert:   true,
				})

				admins, err := b.db.GetAllAdmins(ctx)
				if err == nil {
					var rows []tg.KeyboardButtonRow
					rows = append(rows, tg.KeyboardButtonRow{
						Buttons: []tg.KeyboardButtonClass{
							&tg.KeyboardButtonCallback{
								Text: fmt.Sprintf("🟢 ᴛᴏᴛᴀʟ ᴀᴅᴍɪɴꜱ: %d", len(admins)),
								Data: []byte("admin_stats"),
							},
						},
					})

					for _, adminID := range admins {
						rows = append(rows, tg.KeyboardButtonRow{
							Buttons: []tg.KeyboardButtonClass{
								&tg.KeyboardButtonCallback{
									Text: fmt.Sprintf("👤 %d", adminID),
									Data: []byte(fmt.Sprintf("admin_id_%d", adminID)),
								},
								&tg.KeyboardButtonCallback{
									Text: "🔴 ʀᴇᴍᴏᴠᴇ",
									Data: []byte(fmt.Sprintf("remadmin_%d", adminID)),
								},
							},
						})
					}

					rows = append(rows, tg.KeyboardButtonRow{
						Buttons: []tg.KeyboardButtonClass{
							&tg.KeyboardButtonCallback{
								Text: "🔄 ʀᴇꜰʀᴇꜱʜ",
								Data: []byte("refresh_admins"),
							},
						},
					})

					markup := &tg.ReplyInlineMarkup{Rows: rows}
					req := &tg.MessagesEditMessageRequest{
						Peer:    inputPeer,
						ID:      u.MsgID,
						Message: "👑 **ᴀᴅᴍɪɴɪꜱᴛʀᴀᴛᴏʀꜱ ᴅᴀꜱʜʙᴏᴀʀᴅ**\n\nᴍᴀɴᴀɢᴇ ʙᴏᴛ ᴀᴅᴍɪɴɪꜱᴛʀᴀᴛᴏʀꜱ ᴅɪʀᴇᴄᴛʟʏ ꜰʀᴏᴍ ᴛʜᴇ ɪɴʟɪɴᴇ ᴍᴇɴᴜ ʙᴇʟᴏᴡ.\n\n%s",
					}
					req.SetReplyMarkup(markup)
					// Format with trailing credits as per rule
					req.Message = fmt.Sprintf("👑 **ᴀᴅᴍɪɴɪꜱᴛʀᴀᴛᴏʀꜱ ᴅᴀꜱʜʙᴏᴀʀᴅ**\n\nᴍᴀɴᴀɢᴇ ʙᴏᴛ ᴀᴅᴍɪɴɪꜱᴛʀᴀᴛᴏʀꜱ ᴅɪʀᴇᴄᴛʟʏ ꜰʀᴏᴍ ᴛʜᴇ ɪɴʟɪɴᴇ ᴍᴇɴᴜ ʙᴇʟᴏᴡ.\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ")
					_, _ = api.MessagesEditMessage(ctx, req)
				}
			}
			return nil
		}

		if cmd == "refresh_admins" {
			isAdmin, _ := b.db.IsAdmin(ctx, u.UserID, b.cfg.OwnerID)
			if !isAdmin {
				_, _ = api.MessagesSetBotCallbackAnswer(ctx, &tg.MessagesSetBotCallbackAnswerRequest{
					QueryID: u.QueryID,
					Message: "❌ ᴀᴄᴄᴇꜱꜱ ᴅᴇɴɪᴇᴅ.",
					Alert:   true,
				})
				return nil
			}

			admins, err := b.db.GetAllAdmins(ctx)
			if err == nil {
				var rows []tg.KeyboardButtonRow
				rows = append(rows, tg.KeyboardButtonRow{
					Buttons: []tg.KeyboardButtonClass{
						&tg.KeyboardButtonCallback{
							Text: fmt.Sprintf("🟢 ᴛᴏᴛᴀʟ ᴀᴅᴍɪɴꜱ: %d", len(admins)),
							Data: []byte("admin_stats"),
						},
					},
				})

				for _, adminID := range admins {
					rows = append(rows, tg.KeyboardButtonRow{
						Buttons: []tg.KeyboardButtonClass{
							&tg.KeyboardButtonCallback{
								Text: fmt.Sprintf("👤 %d", adminID),
								Data: []byte(fmt.Sprintf("admin_id_%d", adminID)),
							},
							&tg.KeyboardButtonCallback{
								Text: "🔴 ʀᴇᴍᴏᴠᴇ",
								Data: []byte(fmt.Sprintf("remadmin_%d", adminID)),
							},
						},
					})
				}

				rows = append(rows, tg.KeyboardButtonRow{
					Buttons: []tg.KeyboardButtonClass{
						&tg.KeyboardButtonCallback{
							Text: "🔄 ʀᴇꜰʀᴇꜱʜ",
							Data: []byte("refresh_admins"),
						},
					},
				})

				markup := &tg.ReplyInlineMarkup{Rows: rows}
				req := &tg.MessagesEditMessageRequest{
					Peer:    inputPeer,
					ID:      u.MsgID,
					Message: "👑 **ᴀᴅᴍɪɴɪꜱᴛʀᴀᴛᴏʀꜱ ᴅᴀꜱʜʙᴏᴀʀᴅ**\n\nᴍᴀɴᴀɢᴇ ʙᴏᴛ ᴀᴅᴍɪɴɪꜱᴛʀᴀᴛᴏʀꜱ ᴅɪʀᴇᴄᴛʟʏ ꜰʀᴏᴍ ᴛʜᴇ ɪɴʟɪɴᴇ ᴍᴇɴᴜ ʙᴇʟᴏᴡ.\n\nᴄʀᴇᴅɪᴛꜱ: ᴜɴʀᴀᴛᴇᴅ ᴄᴏᴅᴇʀ ᴛ.ᴍᴇ/ᴜɴʀᴀᴛᴇᴅ_ᴄᴏᴅᴇʀ",
				}
				req.SetReplyMarkup(markup)
				_, _ = api.MessagesEditMessage(ctx, req)

				_, _ = api.MessagesSetBotCallbackAnswer(ctx, &tg.MessagesSetBotCallbackAnswerRequest{
					QueryID: u.QueryID,
					Message: "🔄 ᴅᴀꜱʜʙᴏᴀʀᴅ ʀᴇꜰʀᴇꜱʜᴇᴅ!",
				})
			}
			return nil
		}

		return nil
			}()
		}()
		return nil
	})

	opts := telegram.Options{
		UpdateHandler: gaps,
		Logger:        NewZapLogger(),
		Middlewares: []telegram.Middleware{
			updhook.UpdateHook(gaps.Handle),
		},
	}

	return telegram.BotFromEnvironment(ctx, opts, func(setupCtx context.Context, client *telegram.Client) error {
		api = tg.NewClient(client)
		b.streamer.SetAPI(api)
		sender = message.NewSender(api)
		return nil
	}, func(connCtx context.Context, client *telegram.Client) error {
		_, err := api.BotsSetBotCommands(connCtx, &tg.BotsSetBotCommandsRequest{
			Scope:    &tg.BotCommandScopeDefault{},
			LangCode: "",
			Commands: []tg.BotCommand{
				{Command: "start", Description: "ꜱᴛᴀʀᴛ ᴛʜᴇ ʙᴏᴛ"},
				{Command: "help", Description: "ɢᴇᴛ ᴜꜱᴀɢᴇ ɪɴꜱᴛʀᴜᴄᴛɪᴏɴꜱ"},
				{Command: "ping", Description: "ᴘɪɴɢ ᴛʜᴇ ʙᴏᴛ"},
				{Command: "log", Description: "ɢᴇᴛ ʙᴏᴛ ꜱʏꜱᴛᴇᴍ ʟᴏɢꜱ (ᴀᴅᴍɪɴ)"},
				{Command: "logs", Description: "ɢᴇᴛ ʙᴏᴛ ꜱʏꜱᴛᴇᴍ ʟᴏɢꜱ (ᴀᴅᴍɪɴ)"},
				{Command: "users", Description: "ꜱʜᴏᴡ ᴛᴏᴛᴀʟ ᴜꜱᴇʀꜱ ꜱᴛᴀᴛɪꜱᴛɪᴄꜱ (ᴀᴅᴍɪɴ)"},
				{Command: "ban", Description: "ʙᴀɴ ᴀ ᴜꜱᴇʀ ꜰʀᴏᴍ ᴜꜱɪɴɢ ᴛʜᴇ ʙᴏᴛ (ᴀᴅᴍɪɴ)"},
				{Command: "unban", Description: "ᴜɴʙᴀɴ ᴀ ʙᴀɴɴᴇᴅ ᴜꜱᴇʀ (ᴀᴅᴍɪɴ)"},
				{Command: "add_admin", Description: "ɢʀᴀɴᴛ ᴀᴅᴍɪɴ ᴘʀɪᴠɪʟᴇɢᴇꜱ (ᴏᴡɴᴇʀ)"},
				{Command: "remadmin", Description: "ʀᴇᴍᴏᴠᴇ ᴀᴅᴍɪɴ ᴘʀɪᴠɪʟᴇɢᴇꜱ (ᴏᴡɴᴇʀ)"},
				{Command: "admins", Description: "ꜱʜᴏᴡ ᴀᴅᴍɪɴꜱ ᴅᴀꜱʜʙᴏᴀʀᴅ (ᴀᴅᴍɪɴ)"},
			},
		})
		if err != nil {
			log.Printf("Failed to set bot commands: %v", err)
		} else {
			log.Println("Bot commands registered successfully!")
		}

		// Notify owner that bot started
		ownerPeer := &tg.InputPeerUser{
			UserID: b.cfg.OwnerID,
		}
		_, errNotify := sender.To(ownerPeer).StyledText(connCtx, html.String(nil, "🚀 <b>ʙᴏᴛ ꜱᴛᴀʀᴛᴇᴅ ꜱᴜᴄᴄᴇꜱꜰᴜʟʟʏ!</b>"))
		if errNotify != nil {
			log.Printf("Failed to notify owner: %v", errNotify)
		} else {
			log.Println("Owner notified successfully!")
		}

		self, errSelf := client.Self(connCtx)
		if errSelf != nil {
			return errSelf
		}

		return gaps.Run(connCtx, api, self.ID, updates.AuthOptions{
			IsBot: true,
			OnStart: func(ctx context.Context) {
				log.Println("Updates manager (gaps) started successfully!")
			},
		})
	})
}
