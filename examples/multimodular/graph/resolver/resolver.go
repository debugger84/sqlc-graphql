package resolver

import "multimodular/post/storage"
import commentStorage "multimodular/comment/storage"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PostQueries    *storage.Queries
	CommentQueries *commentStorage.Queries
}
