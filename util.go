package glog

import "errors"

func NewComError(new string, old error) error {
	return errors.New(old.Error() + "\n" + new)
}
