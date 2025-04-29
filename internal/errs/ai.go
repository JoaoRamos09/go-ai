package errs

import "errors"

var ErrAIInvoke = errors.New("error ai invoke")
var ErrAIInsert = errors.New("not possible to insert in vector store")
var ErrAISearch = errors.New("not possible to search in vector store")



