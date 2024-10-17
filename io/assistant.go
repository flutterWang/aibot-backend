package bot_io

import (
	"aibot-backend/model"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	BaseUrl    = "http://58.246.3.170:28001"
	CheckDoc   = "/check/doc"
	AsyncCheck = "/check/doc_async/"
	GetResult  = "/check/result/doc/"
	Chat       = "/chat"
)

func (ai_Io *Io) GetAssistantList(c *gin.Context) ([]*model.AssistantInfo, error) {
	result, err := ai_Io.MongoClient.Collection("assistants").Find(c, bson.D{})
	if err != nil {
		return nil, err
	}

	var assistantInfos []*model.AssistantInfo
	if err = result.All(c, &assistantInfos); err != nil {
		return nil, err
	}

	return assistantInfos, nil
}

func (ai_Io *Io) CreateAssistant(c *gin.Context, req *model.CreateAssistantReq) (string, error) {
	mongoClient := ai_Io.GetMongoClient()
	ID := uuid.New().String()

	item := &model.AssistantInfo{
		Id:         ID,
		Name:       req.Name,
		Avatar:     req.Avatar,
		Desc:       req.Desc,
		CategoryID: req.CategoryID,
		Type:       req.Type,
	}

	_, err := mongoClient.Collection("assistants").InsertOne(c, item)
	if err != nil {
		return "", err
	}

	return ID, nil

}

func (ai_Io *Io) GetAssistant(c *gin.Context, assistantID string) (*model.AssistantInfo, error) {
	mongoClient := ai_Io.GetMongoClient()
	var assistant model.AssistantInfo
	err := mongoClient.Collection("assistants").FindOne(c, bson.M{"id": assistantID}).Decode(&assistant)
	if err != nil {
		return nil, err
	}
	return &assistant, nil
}

func (ai_Io *Io) UpdateAssistant(c *gin.Context, assistant *model.AssistantInfo) error {
	mongoClient := ai_Io.GetMongoClient()
	_, err := mongoClient.Collection("assistants").UpdateOne(c, bson.M{"id": assistant.Id}, bson.M{"$set": assistant})
	if err != nil {
		return err
	}
	return nil
}

func (ai_Io *Io) CheckAssistantDocs(c *gin.Context, config *model.Config, docInfo *model.CacleDocInfo) (string, error) {
	checkReq := &model.ApiCheckDocReq{
		DocID:   docInfo.Id,
		DocPath: docInfo.FilePath,
	}
	reqData, err := json.Marshal(checkReq)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", config.AiServer+CheckDoc, bytes.NewReader(reqData))

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}
	rsp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		log.Println(rsp.StatusCode)
		return "", err
	}

	rspBody, err := io.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	resp := &model.ApiCheckDocResp{}
	if err = json.Unmarshal(rspBody, resp); err != nil {
		return "", err
	}

	return resp.Status, nil
}

func (ai_Io *Io) GetCheckDocResult(c *gin.Context, baseURl, docID string) (*model.CheckAssistantDocResultResp, error) {
	req, err := http.NewRequest("GET", baseURl+GetResult+docID+"/", nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	rsp, err := client.Do(req)
	if err != nil {

		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, errors.New("GetCheckDocResult err: " + rsp.Status)
	}

	rspBody, err := io.ReadAll(rsp.Body)
	if err != nil {
		log.Println("GetCheckDocResult err: ", err)
		return nil, err
	}

	log.Println(string(rspBody))
	resp := &model.CheckAssistantDocResultResp{}
	if err = json.Unmarshal(rspBody, resp); err != nil {
		return nil, err
	}

	log.Println("GetCheckDocResult: ", resp)
	return resp, nil
}
