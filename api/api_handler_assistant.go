package api

import (
	"aibot-backend/model"
	"aibot-backend/modules/assistants"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
)

const (
	AsyncCheck = "/check/doc_async/"
)

func (s *server) handlerAssistantList(c *gin.Context) {
	resp, err := assistants.GetAssistantList(c, s.io)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, resp)
}

func (s *server) handlerAssistantCreate(c *gin.Context) {
	req := &model.CreateAssistantReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}
	resp, err := assistants.CreateAssistant(c, s.io, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, resp)
}

func (s *server) handlerAssistantCheckDocs(c *gin.Context) {
	req := &model.CheckAssistantDocsReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}
	resp, err := assistants.CheckAssistantDocs(c, s.io, s.conf, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, resp)
}

// 不再使用
func (s *server) handlerAssistantAsyncCheckDocs(c *gin.Context) {
	var (
		checkOK     = false
		assistantID string
		docID       string
	)
	// 定义WebSocket升级器
	upGrader := gorilla.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade err: ", err)
		return
	}

	defer func() {
		if checkOK {
			log.Println("check ok")
			err = assistants.UpdateAssistantDocStatus(c, s.io, s.conf.AiServer, model.AsyncCheckAssistantDocsReq{
				AssistantID: assistantID,
				DocID:       docID,
			}, model.CacleDocInfoStatusOk)
			if err != nil {
				return
			}
		} else {
			err = assistants.UpdateAssistantDocStatus(c, s.io, s.conf.AiServer, model.AsyncCheckAssistantDocsReq{
				AssistantID: assistantID,
				DocID:       docID,
			}, model.CacleDocInfoStatusNo)
			if err != nil {
				return
			}
		}

		var data model.AsyncCheckDoc
		data.StatusCode = 999
		data.Content = "数据更新完毕"
		dataStr, err := json.Marshal(data)
		if err != nil {
			return
		}

		writeErr := ws.WriteMessage(gorilla.TextMessage, dataStr)
		if writeErr != nil {
			log.Println("Write error:", writeErr)
			return
		}
		_ = ws.Close()
	}()

	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		fmt.Printf("Received message of type %d: %s\n", messageType, string(message))
		var req model.AsyncCheckAssistantDocsReq
		if err := json.Unmarshal(message, &req); err != nil {
			log.Println("read:", err)
			break
		}

		assistantID = req.AssistantID
		docID = req.DocID

		checkReq := &model.ApiCheckDocReq{
			DocID:    req.DocID,
			DocPath:  req.DocFilePath,
			Template: req.Template,
		}
		reqData, err := json.Marshal(checkReq)
		if err != nil {
			log.Println("read:", err)
			return
		}

		pythonReq, err := http.NewRequest("POST", s.conf.AiServer+AsyncCheck, bytes.NewReader(reqData))
		client := &http.Client{
			Timeout: 10 * time.Minute,
		}
		rsp, err := client.Do(pythonReq)
		if err != nil {
			log.Println("read:", err)
			return
		}
		defer rsp.Body.Close()

		err = assistants.UpdateAssistantDocStatus(c, s.io, s.conf.AiServer, req, model.CacleDocInfoStatusUn)
		if err != nil {
			return
		}

		buf := make([]byte, 1024) // Buffer size can be adjusted as needed.
		for {
			n, readErr := rsp.Body.Read(buf)
			if n > 0 {
				// 先解析 buf[:n], 中是否存在换行符
				// 如果存在，就将 buf[:n] 写入 ws
				// 如果不存在，就继续等待

				if bytes.Contains(buf[:n], []byte("\n")) {
					// 通过换行符切分
					lines := bytes.Split(buf[:n], []byte("\n"))
					for _, line := range lines {
						if len(line) > 0 {
							var data model.AsyncCheckDoc
							err = json.Unmarshal(line, &data)
							if err != nil {
								log.Printf("Error unmarshalling lines JSON: %v, Data: %s", err, string(line))
								continue // 跳过错误的数据，继续处理下一行
							}

							log.Println("data: ", data)

							if data.StatusCode == 200 {
								checkOK = true
							}

							dataStr, err := json.Marshal(data)
							if err != nil {
								log.Printf("Error marshalling JSON: %v", err)
								continue // 跳过错误的数据，继续处理下一行
							}

							writeErr := ws.WriteMessage(gorilla.TextMessage, dataStr)
							if writeErr != nil {
								log.Println("Write error:", writeErr)
								return
							}
						}
					}
					continue
				} else {
					var data model.AsyncCheckDoc
					err = json.Unmarshal(buf[:n], &data)
					if err != nil {
						log.Printf("Error unmarshalling JSON: %v, Data: %s", err, string(buf[:n]))
						continue // 跳过错误的数据，继续处理下一行
					}

					log.Println("data: ", data)

					if data.StatusCode == 200 {
						checkOK = true
					}

					dataStr, err := json.Marshal(data)
					if err != nil {
						log.Printf("Error marshalling JSON: %v", err)
						continue // 跳过错误的数据，继续处理下一行
					}

					writeErr := ws.WriteMessage(gorilla.TextMessage, dataStr)
					if writeErr != nil {
						log.Println("Write error:", writeErr)
						return
					}
				}
			}

			if readErr == io.EOF {
				log.Println("End of stream")
				break
			} else if readErr != nil {
				return
			}
		}

		break
	}
}

func (s *server) handlerAssistantDocs(c *gin.Context) {
	docID := c.Query("doc_id")

	resp, err := assistants.GetAssistantDocsList(c, s.io, docID)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, resp)
}

func (s *server) handlerAssistantUploadDoc(c *gin.Context) {
	req := &model.UploadAssistantDocsReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}
	resp, err := assistants.UploadAssistantDocs(c, s.io, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, resp)
}

func (s *server) handlerAssistantUpdateDocsOK(c *gin.Context) {
	req := &model.AssistantUpdateDocsStatusReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}
	err := assistants.AssistantUpdateDocsStatus(c, s.io, req, model.CacleDocInfoStatusOk)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, nil)
}

func (s *server) handlerAssistantChatList(c *gin.Context) {
	req := &model.AssistantChatListReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}
	resp, err := assistants.GetAssistantChatList(c, s.io, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, resp)
}

// 不再使用 websocket
func (s *server) handlerAssistantChatSend(c *gin.Context) {
	resp, err := assistants.AssistantChatSend(c, s.io, s.conf)
	if err != nil {
		log.Println("AssistantChatSend err: ", err)
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, resp)
}

func (s *server) handlerAssistantChatCreate(c *gin.Context) {
	req := &model.AssistantChatCreateReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}
	resp, err := assistants.AssistantChatCreate(c, s.io, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, resp)
}

func (s *server) handlerAssistantChatDelete(c *gin.Context) {
	req := &model.AssistantChatDeleteReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}
	resp, err := assistants.AssistantChatDelete(c, s.io, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, resp)
}

func (s *server) handlerAssistantChatUpdate(c *gin.Context) {
	req := &model.AssistantChatUpdateReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}
	err := assistants.AssistantChatUpdate(c, s.io, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, nil)
}

func (s *server) handlerAssistantChatSaveConversation(c *gin.Context) {
	req := &model.AssistantChatSaveConversationReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}
	err := assistants.AssistantChatSaveConversation(c, s.io, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, nil)
}

func (s *server) handlerAssistantCheckDocResult(c *gin.Context) {
	req := &model.CheckAssistantDocResultReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}
	resp, err := assistants.CheckAssistantDocResult(c, s.io, s.conf.AiServer, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, resp)
}

func (s *server) handlerAssistantCheckDocClear(c *gin.Context) {
	req := &model.CheckAssistantDocClearReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}
	err := assistants.CheckAssistantDocClear(c, s.io, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, nil)
}

func (s *server) handlerAssistantChatStart(c *gin.Context) {
	chatID := c.Query("chat_id")
	if chatID == "" {
		// 新建一个 Chat 数据
	}
}
