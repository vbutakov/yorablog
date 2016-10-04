package main

import "yoradb"

var (
	sessions map[string]*yoradb.User
)

func InitSessionsMap() {
	sessions = make(map[string]*yoradb.User)

	sessions["01"] = &yoradb.User{ID: 1,
		Name:             "User1",
		Email:            "a1@b.c",
		Password:         yoradb.GetPasswordHash("a1@b.c", "111"),
		CreatePostPermit: false,
		EditPostPermit:   false}

	sessions["02"] = &yoradb.User{ID: 2,
		Name:             "User2",
		Email:            "a2@b.c",
		Password:         yoradb.GetPasswordHash("a2@b.c", "222"),
		CreatePostPermit: false,
		EditPostPermit:   true}

	sessions["03"] = &yoradb.User{ID: 3,
		Name:             "User3",
		Email:            "a3@b.c",
		Password:         yoradb.GetPasswordHash("a3@b.c", "333"),
		CreatePostPermit: true,
		EditPostPermit:   false}

	sessions["04"] = &yoradb.User{ID: 4,
		Name:             "User4",
		Email:            "a4@b.c",
		Password:         yoradb.GetPasswordHash("a4@b.c", "444"),
		CreatePostPermit: true,
		EditPostPermit:   true}
}
