package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

const messageFetchSize = 10

type Message struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	RoomID    bson.ObjectId `bson:"room_id" json:"room_id"`
	Content   string        `bson:"content" json:"content"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	User      *User         `bson:"user" json:"user"`
}

func (m *Message) create() error {
	// 몽고 DB 세션 생성
	session := mongoSession.Copy()
	defer session.Close()

	// 몽고 DB ID 생성
	m.ID = bson.NewObjectId()
	// 메세지 생성 시간 기록
	m.CreatedAt = time.Now()
	// message 정보 저장을 위한 몽고DB 컬렉션 객체 생성
	c := session.DB("test").C("messages")

	// messages 컬렉션에 message 정보 저장
	if err := c.Insert(m); err != nil {
		return err
	}
	return nil
}

func retireveMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mongoSession.Copy()

	defer session.Close()

	// 쿼리 매개변수로 전달된 limit 값 확인
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		// 정상적인 limit 값이 전달되지 않으면 limit를 messageFetchSize으로 세팅
		limit = messageFetchSize
	}

	var messages []Message
	// _id 역순으로 정렬하여 limit 수만큼 message 조회
	err = session.DB("test").C("messages").Find(bson.M{"room_id": bson.ObjectIdHex(ps.ByName("id"))}).
		Sort("-_id").Limit(limit).All(&messages)

	if err != nil {
		// 오류 발생 시 500 에러 반환
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}

	// message 조회 결과 반환
	renderer.JSON(w, http.StatusOK, messages)

}
