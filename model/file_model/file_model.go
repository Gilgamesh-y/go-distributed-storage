package file_model

import (
	"DistributedStorage/fileMeta"
	"DistributedStorage/model"
)

func Insert(fm fileMeta.FileMeta) (int64, error) {
	state, err := model.GetConn().Prepare("insert into files (`name`, `size`, `path`) values (?,?,?)")
	println(1)
	if err != nil {
		return 0, err
	}
	defer state.Close()
	res, err := state.Exec(fm.Name, fm.Size, fm.Path)
	println(2)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func Get() (*fileMeta.FileMeta, error) {
	state, err := model.GetConn().Prepare("select name, size, path, updated_at from files order by id desc ")
	defer state.Close()
	if err != nil {
		return nil, err
	}
	fm := fileMeta.FileMeta{}
	err = state.QueryRow().Scan(&fm.Name, &fm.Size, &fm.Path, &fm.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &fm, nil
}