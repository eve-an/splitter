package feature

type UserID uint64

type User struct {
	id          UserID
	anonymousID UserID
}

func NewUser(id UserID) User {
	return User{id: id}
}

func NewAnonymousUser(id UserID) User {
	return User{anonymousID: id}
}

func (u User) ID() UserID {
	if u.id == 0 && u.anonymousID == 0 {
		panic("user has no ID")
	}

	if u.id == 0 {
		return u.anonymousID
	}

	return u.id
}
