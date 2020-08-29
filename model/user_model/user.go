package user_model

import (
	"DistributedStorage/model"
	"crypto/md5"
	"fmt"
)

type User struct {
	Id int64 `form:"id"`
	Name string `form:"name" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func AddUser(user *User) (error) {
	state, err := model.GetConn().Prepare("insert into users (`name`, `password`) values (?,?)")
	if err != nil {
		return err
	}
	defer state.Close()
	res, err := state.Exec(user.Name, fmt.Sprintf("%x", md5.Sum([]byte(user.Password))))
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.Id = id
	return nil
}

func GetUser(user *User) error {
	state, err := model.GetConn().Prepare("select id, name, password from users where id = ?")
	if err != nil {
		return err
	}
	defer state.Close()
	err = state.QueryRow(user.Id).Scan(&user.Id, &user.Name, &user.Password)
	if err != nil {
		return err
	}
	return nil
}