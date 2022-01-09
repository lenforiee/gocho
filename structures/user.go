package structures

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/lenforiee/gocho/globals"
	"github.com/lenforiee/gocho/logger"
	"github.com/savsgio/gotils/uuid"
)

var usersLock = &sync.Mutex{}
var USERS []User

type Action struct {
	ID      uint8
	Text    string
	BmapMD5 string
	BmapID  int32
	Mods    uint32
}

type User struct {
	ID         int32
	Name       string `db:"username"`
	SafeName   string `db:"username_safe"`
	UUID       string
	Country    uint8
	Privileges uint64 `db:"privileges"`
	Location   [2]float32
	IP         string
	IsBot      bool

	TotalScore  int64
	RankedScore int64
	PP          int16
	PlayCount   int32
	PlayTime    uint64
	Accuracy    float32
	MaxCombo    uint16
	Rank        int32

	Action       Action
	CurrentMode  uint8
	CurrentCmode uint8
	TimeOffset   uint8

	PacketQueue bytes.Buffer
	PacketLock  *sync.Mutex
}

func Broadcast(data []byte) {
	for _, user := range USERS {
		user.Queue(data)
	}
}

func NewUser(userID int) (u User, err error) {
	u.ID = int32(userID)
	u.UUID = uuid.V4()
	u.IsBot = false
	//err = globals.DbConn.Get(&u, "SELECT username, username_safe, privileges FROM users WHERE id = ?", userID)
	//if err != nil {
	//	return u, err
	//}
	u.CurrentMode = uint8(0)
	u.CurrentCmode = uint8(0)
	//u.UpdateStats()
	u.PacketLock = &sync.Mutex{}
	usersLock.Lock()
	USERS = append(USERS, u)
	usersLock.Unlock()
	return u, nil
}

func (u *User) UpdateStats() {
	// TODO: rewrite it.
	table := "users_stats"
	if u.CurrentCmode == 1 {
		table = "rx_stats"
	} else if u.CurrentCmode == 2 {
		table = "ap_stats"
	}
	suffix := "std"
	if u.CurrentMode == 1 {
		suffix = "taiko"
	} else if u.CurrentMode == 2 {
		suffix = "ctb"
	} else if u.CurrentMode == 3 {
		suffix = "mania"
	}
	err := globals.DbConn.QueryRow(fmt.Sprintf(`
	SELECT
		total_score_%[1]s, ranked_score_%[1]s, pp_%[1]s, 
		playcount_%[1]s, playtime_%[1]s, accuracy_%[1]s, 
		max_combo_%[1]s
	FROM %[2]s
	WHERE id = ? LIMIT 1
	`, suffix, table), u.ID).Scan(
		&u.TotalScore, &u.RankedScore, &u.PP,
		&u.PlayCount, &u.PlayTime, &u.Accuracy,
		&u.MaxCombo,
	)
	if err != nil {
		logger.Error(fmt.Sprintf("There was an issue updating %s stats!", u.Name))
		fmt.Println(err)
		return
	}
}

func (u *User) Queue(data []byte) {
	u.PacketLock.Lock()
	u.PacketQueue.Write(data)
	u.PacketLock.Unlock()
}

func (u *User) Dequeue() (resp []byte) {
	u.PacketLock.Lock()
	resp = u.PacketQueue.Bytes()
	u.PacketQueue.Reset()
	u.PacketLock.Unlock()
	return resp
}
