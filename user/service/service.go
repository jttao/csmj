package service

import (
	"context"
	"fmt"
	gamedb "game/db"
	gameredis "game/redis"
	usermodel "game/user/model"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

const (
	userKey = "user"
)

type Sex int

const (
	SexUnknown Sex = iota
	SexBoy
	SexGirl
)

func (s Sex) Valid() bool {
	switch s {
	case SexUnknown:
	case SexBoy:
	case SexGirl:
	default:
		return false
	}
	return true
}

type UserConfig struct {
	//毫秒
	Key         string `json:"key"`
	ExpiredTime int64  `json:"expiredTime"`
}

type UserService interface {
	WxLogin(wxId string, deviceMac string, wxName string, wxSex Sex, wxImage string) (t string, expiredTime int64, err error)
	VisitLogin(deviceMac string) (t string, expiredTime int64, err error)
	Verify(tokenStr string) (int64, error)
	Logout(id int64) error
	Login(id int64) (t string, expiredTime int64, err error)
	GetCardNum(id int64) (int64, error)
	IsForbid(id int64) (bool, error)
	ChangeCardNum(id int64, num int64, reason usermodel.ReasonType) error
	GetUserById(id int64) (u *usermodel.User, err error)
	UpdateUser(id int64, name string, lastLoginIp string) error
	Online(id int64) error
	Offline(id int64) error
	ChangeUserLocation(id int64,location string) error
}

type userService struct {
	userConfig *UserConfig
	key        []byte
	rs         gameredis.RedisService
	db         gamedb.DBService
}

func (us *userService) WxLogin(wxId string, deviceMac string, wxName string, wxSex Sex, wxImage string) (t string, expiredTime int64, err error) {
	user := &usermodel.User{}
	tdb := us.db.DB().First(user, "weixin=?", wxId)
	if tdb.Error != nil {
		if tdb.Error != gorm.ErrRecordNotFound {
			err = tdb.Error
			return
		}
		user, err = us.register(wxId, deviceMac, wxName, wxSex, wxImage)
		if err != nil {
			return
		}
	}
	user.Name = wxName
	user.Sex = int(wxSex)
	user.Image = wxImage

	tdb = us.db.DB().Model(user).Update(user)
	err = tdb.Error
	if err != nil {
		return
	}
	t, expiredTime, err = us.Login(user.Id)
	return
}

func (us *userService) Logout(id int64) (err error) {
	conn := us.rs.Pool().Get()
	if conn.Err() != nil {
		err = conn.Err()
		return
	}
	defer conn.Close()

	userTokenKey := gameredis.Combine(userKey, fmt.Sprintf("%s", id))
	_, err = conn.Do("del", userTokenKey)
	if err != nil {
		return err
	}
	return nil
}

func (us *userService) Login(id int64) (t string, expiredTime int64, err error) {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	expiredTime += now + us.userConfig.ExpiredTime
	claims := &jwt.StandardClaims{}
	claims.ExpiresAt = expiredTime
	claims.Issuer = fmt.Sprintf("%d", id)
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	t, err = token.SignedString(us.key)
	if err != nil {
		return
	}
	//保存redis
	conn := us.rs.Pool().Get()
	if conn.Err() != nil {
		err = conn.Err()
		return
	}
	defer conn.Close()

	userTokenKey := gameredis.Combine(userKey, fmt.Sprintf("%d", id))

	ok, err := redis.String(conn.Do("setex", userTokenKey, us.userConfig.ExpiredTime/1000, t))
	if err != nil {
		return
	}
	if ok != gameredis.OK {
		err = fmt.Errorf("redis set failed %s", ok)
		return
	}
	return
}

func (us *userService) Verify(tokenStr string) (id int64, err error) {
	claims := &jwt.StandardClaims{}
	_, err = jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) { return us.key, nil })
	if err != nil {
		return 0, err
	}

	idStr := claims.Issuer

	if len(idStr) == 0 {
		return 0, nil
	}
	//TODO 为什么不能类型转换int
	id, err = strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}

	//保存redis
	conn := us.rs.Pool().Get()
	if conn.Err() != nil {
		err = conn.Err()
		return 0, err
	}
	defer conn.Close()

	userTokenKey := gameredis.Combine(userKey, fmt.Sprintf("%d", id))

	cacheToken, err := redis.String(conn.Do("get", userTokenKey))
	if err != nil {
		return
	}
	if cacheToken != tokenStr {
		return 0, nil
	}

	return id, nil
}

func (us *userService) register(wxId string, deviceMac string, name string, sex Sex, image string) (user *usermodel.User, err error) {
	user = &usermodel.User{}
	user.Weixin = wxId
	user.DeviceMac = deviceMac
	user.Name = name
	user.CardNum = 3
	user.Sex = int(sex)
	user.Image = image
	user.CreateTime = time.Now().UnixNano() / int64(time.Millisecond)
	user.UpdateTime = user.CreateTime
	tdb := us.db.DB().Save(user)
	if tdb.Error != nil {
		return nil, tdb.Error
	}
	return
}

func (us *userService) VisitLogin(deviceMac string) (t string, expiredTime int64, err error) {
	user := &usermodel.User{}
	tdb := us.db.DB().First(user, "deviceMac = ?", deviceMac)
	if tdb.Error != nil {
		if tdb.Error != gorm.ErrRecordNotFound {
			err = tdb.Error
			return
		}
		user, err = us.register("", deviceMac, "", SexUnknown, "")
		if err != nil {
			return
		}
	}
	t, expiredTime, err = us.Login(user.Id)
	return
}

func (us *userService) GetCardNum(id int64) (cards int64, err error) {
	user := &usermodel.User{}
	tdb := us.db.DB().First(user, "id=?", id)
	err = tdb.Error
	if err != nil {
		return 0, err
	}
	return user.CardNum, nil
}

func (us *userService) IsForbid(id int64) (flag bool, err error) {
	user := &usermodel.User{}
	tdb := us.db.DB().First(user, "id=?", id)
	err = tdb.Error
	if err != nil {
		return false, err
	}
	return user.Forbid == 1, nil
}

func (us *userService) ChangeCardNum(id int64, num int64, reason usermodel.ReasonType) error {
	user := &usermodel.User{}
	tdb := us.db.DB().First(user, "id=?", id)
	err := tdb.Error
	if err != nil {
		return err
	}

	newCardNum := user.CardNum + num
	tdb = us.db.DB().Model(user).Update("cardNum", newCardNum)
	err = tdb.Error
	if err != nil {
		return err
	}

	//添加日志
	cardRecordModel := &usermodel.CardRecordModel{}
	cardRecordModel.ChangeNum = num
	cardRecordModel.UserId = id
	cardRecordModel.Reason = int(reason)
	cardRecordModel.CreateTime = time.Now().UnixNano() / int64(time.Millisecond)
	tdb = us.db.DB().Save(cardRecordModel)
	if tdb.Error != nil {
		log.WithFields(
			log.Fields{
				"userId": id,
				"num":    num,
			}).Error("保存房卡纪录失败")
	}
	return nil
}

func (us *userService) ChangeUserLocation(id int64,location string) error {
	user := &usermodel.User{}
	tdb := us.db.DB().First(user, "id=?", id)
	err := tdb.Error
	if err != nil {
		return err
	}
	
	tdb = us.db.DB().Model(user).Update("location", location)
	err = tdb.Error
	if err != nil {
		return err
	} 
	return nil
}

func (us *userService) GetUserById(id int64) (u *usermodel.User, err error) {
	u = &usermodel.User{}
	tdb := us.db.DB().First(u, "id=?", id)
	err = tdb.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func (us *userService) UpdateUser(id int64, name string, lastLoginIp string) error {
	user := &usermodel.User{}
	user.Id = id
	user.Name = name
	user.LastLoginIp = lastLoginIp
	now := time.Now().UnixNano() / int64(time.Millisecond)
	user.LastLoginTime = now
	user.UpdateTime = now
	tdb := us.db.DB().Model(user).Update(user)
	err := tdb.Error
	if err != nil {
		return err
	}
	return nil
}

func (us *userService) Online(id int64) error {
	user := &usermodel.User{}
	user.Id = id
	now := time.Now().UnixNano() / int64(time.Millisecond)
	user.UpdateTime = now
	user.State = 1
	tdb := us.db.DB().Model(user).Update(user)
	err := tdb.Error
	if err != nil {
		return err
	}
	return nil
}

func (us *userService) Offline(id int64) error {
	user := &usermodel.User{}
	user.Id = id

	tdb := us.db.DB().Model(user).Update("state", 0)
	err := tdb.Error
	if err != nil {
		return err
	}
	return nil
}

func NewUserService(uc *UserConfig, db gamedb.DBService, rs gameredis.RedisService) (us UserService, err error) {

	//读取key
	keyFile, err := filepath.Abs(uc.Key)
	if err != nil {
		return
	}
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return
	}
	us = &userService{
		userConfig: uc,
		key:        key,
		rs:         rs,
		db:         db,
	}
	return
}

const (
	key = "user_service"
)

func WithUserService(ctx context.Context, us UserService) context.Context {
	return context.WithValue(ctx, key, us)
}

func UserServiceInContext(ctx context.Context) UserService {
	us, ok := ctx.Value(key).(UserService)
	if !ok {
		return nil
	}
	return us
}

const (
	userContextKey = "user"
)

func WithUser(ctx context.Context, userId int64) context.Context {
	return context.WithValue(ctx, userContextKey, userId)
}

func UserInContext(ctx context.Context) int64 {
	userId, ok := ctx.Value(userContextKey).(int64)
	if !ok {
		return 0
	}
	return userId
}

func AuthHandler() negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, hf http.HandlerFunc) {
		us := UserServiceInContext(req.Context())
		if us == nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		if ah := req.Header.Get("Authorization"); ah != "" {

			// Should be a bearer token
			if len(ah) > 6 && strings.ToUpper(ah[0:6]) == "BEARER" {
				id, err := us.Verify(ah[7:])
				if err != nil {
					rw.WriteHeader(http.StatusUnauthorized)
					return
				}
				if id == 0 {
					rw.WriteHeader(http.StatusUnauthorized)
					return
				}
				ctx := WithUser(req.Context(), id)
				nreq := req.WithContext(ctx)
				hf.ServeHTTP(rw, nreq)
				return
			}
		}
		rw.WriteHeader(http.StatusUnauthorized)
	})
}
func SetupUserServiceHandler(us UserService) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, hf http.HandlerFunc) {
		ctx := req.Context()
		nctx := WithUserService(ctx, us)
		nreq := req.WithContext(nctx)
		hf(rw, nreq)
	})
}
