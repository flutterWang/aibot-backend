package assistants

import (
	"aibot-backend/io"
	"aibot-backend/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	gorilla "github.com/gorilla/websocket"
)

const (
	Chat = "/chat"
)

func GetAssistantList(c *gin.Context, io *bot_io.Io) ([]*model.AssistantInfo, error) {
	result, err := io.GetAssistantList(c)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func CreateAssistant(c *gin.Context, io *bot_io.Io, req *model.CreateAssistantReq) (string, error) {
	result, err := io.CreateAssistant(c, req)
	if err != nil {
		return "", err
	}

	return result, nil
}

func CheckAssistantDocs(c *gin.Context, io *bot_io.Io, config *model.Config, req *model.CheckAssistantDocsReq) (*model.CheckAssistantDocsResp, error) {
	assistant, err := io.GetAssistant(c, req.AssistantID)
	if err != nil {
		log.Println("get assistant err: ", err)
		return nil, err
	}

	var docInfo *model.CacleDocInfo
	if req.DocID != "" {
		for _, doc := range assistant.Docs {
			if doc.Id == req.DocID {
				docInfo = doc
			}
		}
	} else {
		log.Println("docID is empty")
		return nil, errors.New("docID is empty")
	}

	result, err := io.CheckAssistantDocs(c, config, docInfo)
	if err != nil {
		log.Println("check assistant docs err: ", err)
		return nil, err
	}
	log.Println("result: ", result)
	if result == "OK" {
		docInfo.Status = model.CacleDocInfoStatusOk
	} else {
		docInfo.Status = model.CacleDocInfoStatusNo
	}
	err = io.UpdateAssistant(c, assistant)
	if err != nil {
		log.Println("update  UpdateAssistant assistant err: ", err)
		return nil, err
	}

	return &model.CheckAssistantDocsResp{
		Status: result,
		DocID:  docInfo.Id,
	}, nil
}

func GetAssistantDocsList(c *gin.Context, io *bot_io.Io, assistantID string) ([]*model.CacleDocInfo, error) {
	assistant, err := io.GetAssistant(c, assistantID)
	if err != nil {
		return nil, err
	}

	return assistant.Docs, nil
}

func UpdateAssistantDocStatus(c *gin.Context, io *bot_io.Io, aiServer string, req model.AsyncCheckAssistantDocsReq, status string) error {
	assistant, err := io.GetAssistant(c, req.AssistantID)
	if err != nil {
		return err
	}

	for _, doc := range assistant.Docs {
		if doc.Id == req.DocID {
			doc.Status = status

			log.Println("status: ", status)
			if status == model.CacleDocInfoStatusOk {
				result, err := io.GetCheckDocResult(c, aiServer, req.DocID)
				if err != nil {
					return err
				}

				if result != nil || len(result.CheckResult) != 0 {
					doc.Result = "1"
				}

				for _, resultItem := range result.CheckResult {
					log.Println("---\n resultItem: ", resultItem, " \n Result: ", resultItem.Result, "\n---")
					if resultItem.Result == 0 {
						doc.Result = "0"
						break
					}
				}
			}
		}
	}

	err = io.UpdateAssistant(c, assistant)
	if err != nil {
		return err
	}

	return nil
}

func UploadAssistantDocs(c *gin.Context, io *bot_io.Io, req *model.UploadAssistantDocsReq) (*model.UploadAssistantDocsResp, error) {
	assistant, err := io.GetAssistant(c, req.AssistantID)
	if err != nil {
		return nil, err
	}

	docInfo := &model.CacleDocInfo{
		Name:     req.DocName,
		FilePath: req.DocFilePath,
		Id:       uuid.New().String(),
		Status:   model.CacleDocInfoStatusUn,
	}

	assistant.Docs = append(assistant.Docs, docInfo)
	err = io.UpdateAssistant(c, assistant)
	if err != nil {
		return nil, err
	}

	return &model.UploadAssistantDocsResp{
		DocID: docInfo.Id,
	}, nil
}

func AssistantUpdateDocsStatus(c *gin.Context, io *bot_io.Io, req *model.AssistantUpdateDocsStatusReq, status string) error {
	assistant, err := io.GetAssistant(c, req.AssistantID)
	if err != nil {
		return err
	}

	for i := range assistant.Docs {
		if assistant.Docs[i].Id == req.DocID {
			assistant.Docs[i].Status = status
		}
	}

	err = io.UpdateAssistant(c, assistant)
	if err != nil {
		return err
	}

	return nil
}

func GetAssistantChatList(c *gin.Context, io *bot_io.Io, req *model.AssistantChatListReq) (*model.AssistantChatListResp, error) {
	assistant, err := io.GetAssistant(c, req.AssistantID)
	if err != nil {
		return nil, err
	}

	// 只要 DeleteAt 为 0 的

	chats := make([]*model.ChatInfo, 0)
	for i := 0; i < len(assistant.Chats); i++ {
		if assistant.Chats[i].DeleteAt == 0 {
			chats = append(chats, assistant.Chats[i])
		}
	}

	return &model.AssistantChatListResp{
		Chats: chats,
	}, nil
}

func AssistantChatSend(c *gin.Context, botIO *bot_io.Io, config *model.Config) ([]*model.Conversation, error) {
	// 定义WebSocket升级器
	var (
		chatDone  = false
		assistant *model.AssistantInfo
	)
	upGrader := gorilla.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade err: ", err)
		return nil, err
	}

	defer func() {
		log.Println("close websocket")
		_ = ws.Close()

		if chatDone {
			err = botIO.UpdateAssistant(c, assistant)
			if err != nil {
				log.Println("update  UpdateAssistant assistant err: ", err)
				return
			}
		}
	}()

	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		fmt.Printf("Received message of type %d: %s\n", messageType, string(message))
		var req model.AssistantChatSendReq
		if err := json.Unmarshal(message, &req); err != nil {
			log.Println("read:", err)
			break
		}

		assistant, err = botIO.GetAssistant(c, req.AssistantID)
		if err != nil {
			return nil, err
		}

		var chatInfo *model.ChatInfo
		for _, chat := range assistant.Chats {
			if chat.Id == req.ChatID {
				chatInfo = chat
			}
		}

		if chatInfo == nil {
			return nil, errors.New("chat not found")
		}

		chatInfo.Conversation = append(chatInfo.Conversation, &model.Conversation{
			Role:    "user",
			Content: req.Content,
			Time:    time.Now().Format("2006-01-02 15:04:05"),
		})

		sendReq := &model.ApiSendContentToAIReq{
			BotName:      "motorbot",
			Conversation: chatInfo.Conversation,
			Details:      "xxx",
		}
		reqData, err := json.Marshal(sendReq)
		if err != nil {
			return nil, err
		}

		pythonReq, err := http.NewRequest("POST", config.AiServer+Chat, bytes.NewReader(reqData))
		client := &http.Client{
			Timeout: 2 * time.Minute,
		}
		rsp, err := client.Do(pythonReq)
		if err != nil {
			log.Println("read:", err)
			return nil, err
		}
		defer rsp.Body.Close()

		buffer := &bytes.Buffer{}
		for {
			chunk := make([]byte, 65535)
			n, readErr := rsp.Body.Read(chunk)
			if readErr != nil && readErr != io.EOF {
				log.Println("read:", readErr)
				break
			}

			if readErr == io.EOF {
				log.Println("Read error:", readErr)
				break
			}
			if n > 0 {
				buffer.Write(chunk[:n])
				data := buffer.Bytes()

				// 如果可以解析，就直接解析，并且输出
				if json.Valid(data) {
					var chatRsp model.StreamChatResp
					if unmarshalErr := json.Unmarshal(data, &chatRsp); unmarshalErr == nil {
						dataStr, err := json.Marshal(chatRsp)
						if err != nil {
							log.Printf("Error marshalling JSON: %v", err)
							continue // 跳过错误的数据，继续处理下一行
						}

						log.Println("get data")

						if chatRsp.StatusCode == 200 {
							log.Println("in here")
							var chatDoneData model.ChatDone
							if unmarshalErr := json.Unmarshal([]byte(chatRsp.Content), &chatDoneData); unmarshalErr == nil {
								log.Println("chatOK")
								chatInfo.Conversation = append(chatInfo.Conversation, chatDoneData.Conversation)
								chatDone = true
							} else {
								log.Printf("Error unmarshalling JSON: %v, Data: %s", unmarshalErr, string(data))
								continue
							}
						}

						writeErr := ws.WriteMessage(gorilla.TextMessage, dataStr)
						if writeErr != nil {
							log.Println("Write error:", writeErr)
						}

						buffer.Reset() //清空缓冲区以便接受下一条消息;
					} else {
						log.Println("Failed to unmarshal JSON:", unmarshalErr)
					}
				} else if bytes.Contains(data, []byte("\n")) {
					// 通过换行符切分
					lines := bytes.Split(data, []byte("\n"))
					for _, line := range lines {
						if len(line) > 0 && json.Valid(line) {
							var chatRsp model.StreamChatResp
							if unmarshalErr := json.Unmarshal(line, &chatRsp); unmarshalErr == nil {
								dataStr, err := json.Marshal(chatRsp)
								if err != nil {
									log.Printf("Error marshalling JSON: %v", err)
									continue // 跳过错误的数据，继续处理下一行
								}

								log.Println("get data")

								if chatRsp.StatusCode == 200 {
									log.Println("line in here")
									var content model.StreamChatResp
									if unmarshalErr := json.Unmarshal([]byte(chatRsp.Content), &content); unmarshalErr == nil {
										log.Println("chatOK")
										var chatDoneData model.ChatDone
										if unmarshalErr := json.Unmarshal([]byte(content.Content), &chatDoneData); unmarshalErr == nil {
											log.Println("chatOK")
											chatInfo.Conversation = append(chatInfo.Conversation, chatDoneData.Conversation)
											chatDone = true
										} else {
											log.Printf("Error unmarshalling JSON: %v, Data: %s", unmarshalErr, string(data))
											continue
										}
									} else {
										log.Printf("Error unmarshalling JSON: %v, Data: %s", unmarshalErr, string(data))
										continue
									}
								}

								writeErr := ws.WriteMessage(gorilla.TextMessage, dataStr)
								if writeErr != nil {
									log.Println("Write error:", writeErr)
								}

								buffer.Reset() //清空缓冲区以便接受下一条消息;
							} else {
								log.Println("Failed to unmarshal JSON:", unmarshalErr)
							}
						} else {
							break //跳出内层循环并等待更多数据
						}
					}

				} else {
					log.Println("Invalid data:", string(data))
				}
			}
		}
	}

	return nil, nil
}

func AssistantChatCreate(c *gin.Context, io *bot_io.Io, req *model.AssistantChatCreateReq) (*model.AssistantChatCreateResp, error) {
	assistant, err := io.GetAssistant(c, req.AssistantID)
	if err != nil {
		return nil, err
	}

	chatID := uuid.New().String()
	assistant.Chats = append(assistant.Chats, &model.ChatInfo{
		Id:       chatID,
		Name:     req.ChatName,
		CreateAt: time.Now().Unix(),
	})

	err = io.UpdateAssistant(c, assistant)
	if err != nil {
		return nil, err
	}

	return &model.AssistantChatCreateResp{
		ChatID: chatID,
	}, nil
}

func AssistantChatDelete(c *gin.Context, io *bot_io.Io, req *model.AssistantChatDeleteReq) (string, error) {
	assistant, err := io.GetAssistant(c, req.AssistantID)
	if err != nil {
		return "", err
	}

	chatInfo := &model.ChatInfo{}
	for _, chat := range assistant.Chats {
		if chat.Id == req.ChatID {
			chatInfo = chat
		}
	}

	if chatInfo == nil {
		return "", errors.New("chat not found")
	}

	chatInfo.DeleteAt = time.Now().Unix()

	err = io.UpdateAssistant(c, assistant)
	if err != nil {
		return "", err
	}

	return chatInfo.Id, nil
}

func AssistantChatUpdate(c *gin.Context, io *bot_io.Io, req *model.AssistantChatUpdateReq) error {
	assistant, err := io.GetAssistant(c, req.AssistantID)
	if err != nil {
		return err
	}

	chatInfo := &model.ChatInfo{}
	for _, chat := range assistant.Chats {
		if chat.Id == req.ChatID {
			chatInfo = chat
		}
	}

	if chatInfo == nil {
		return errors.New("chat not found")
	}

	chatInfo.Name = req.ChatName
	err = io.UpdateAssistant(c, assistant)
	if err != nil {
		return err
	}

	return nil
}

func AssistantChatSaveConversation(c *gin.Context, io *bot_io.Io, req *model.AssistantChatSaveConversationReq) error {
	assistant, err := io.GetAssistant(c, req.AssistantID)
	if err != nil {
		return err
	}

	chatInfo := &model.ChatInfo{}
	for _, chat := range assistant.Chats {
		if chat.Id == req.ChatID {
			chatInfo = chat
		}
	}

	if chatInfo == nil {
		return errors.New("chat not found")
	}

	chatInfo.Conversation = append(chatInfo.Conversation, req.Conversation...)
	err = io.UpdateAssistant(c, assistant)
	if err != nil {
		return err
	}

	return nil
}

func CheckAssistantDocResult(c *gin.Context, io *bot_io.Io, aiServer string, req *model.CheckAssistantDocResultReq) (*model.CheckAssistantDocResultResp, error) {
	_, err := io.GetAssistant(c, req.AssistantID)
	if err != nil {
		return nil, err
	}

	result, err := io.GetCheckDocResult(c, aiServer, req.DocID)
	if err != nil {
		log.Println("get check doc result err: ", err)
		return nil, err
	}

	log.Println("result: ", result)

	return result, nil

}

func CheckAssistantDocClear(c *gin.Context, io *bot_io.Io, req *model.CheckAssistantDocClearReq) error {
	assistant, err := io.GetAssistant(c, req.AssistantID)
	if err != nil {
		return err
	}

	assistant.Docs = make([]*model.CacleDocInfo, 0)

	err = io.UpdateAssistant(c, assistant)
	if err != nil {
		return err
	}

	return nil
}
