package model

type Room struct {
	name       string
	ownerName  string
	usersNames map[string]bool
}

func NewRoom(name string, ownerName string, userNames map[string]bool) Room {
	return Room{name: name, ownerName: ownerName, usersNames: userNames}
}

func (r Room) IsUserInRoom(userName string) bool {
	_, ok := r.usersNames[userName]

	if !ok {
		return false
	}

	return true
}

func (r Room) AppendUserToRoom(userName string) {
	r.usersNames[userName] = true
}

func (r Room) UserNames() map[string]bool {
	return r.usersNames
}

func (r Room) DeleteUserFromRoom(userName string) {
	delete(r.usersNames, userName)
}
