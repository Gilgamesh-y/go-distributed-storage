package file_model

import (
	"DistributedStorage/fileMeta"
	"DistributedStorage/model"
)

func Insert(fm *fileMeta.FileMeta) (int64, error) {
	state, err := model.GetConn().Prepare("insert into files (`name`, `size`, `hash`, `path`) values (?,?,?,?)")
	if err != nil {
		return 0, err
	}
	defer state.Close()
	res, err := state.Exec(fm.Name, fm.Size, fm.Hash, fm.Path)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetByHash(hash string) (*fileMeta.FileMeta, error) {
	state, err := model.GetConn().Prepare("select id, name, size, path, updated_at from files where hash = ? order by id desc")
	defer state.Close()
	if err != nil {
		return nil, err
	}
	fm := fileMeta.FileMeta{}
	state.QueryRow(hash).Scan(&fm.Id, &fm.Name, &fm.Size, &fm.Path, &fm.UpdatedAt)

	return &fm, nil
}

func Get() (*fileMeta.FileMeta, error) {
	state, err := model.GetConn().Prepare("select id, name, size, path, updated_at from files order by id desc")
	defer state.Close()
	if err != nil {
		return nil, err
	}
	fm := fileMeta.FileMeta{}
	state.QueryRow().Scan(&fm.Id, &fm.Name, &fm.Size, &fm.Path, &fm.UpdatedAt)

	return &fm, nil
}