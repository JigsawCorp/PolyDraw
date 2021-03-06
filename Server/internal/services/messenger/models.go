package messenger

//MessageSent represent a message that was sent by the client
type MessageSent struct {
	Message   string
	ChannelID string
}

//MessageReceived represent a message that was relayed by the server
type MessageReceived struct {
	ChannelID string
	Timestamp int64
	UserID    string
	Username  string
	Message   string
}

//ChannelJoin represent a message that was sent by the server to inform the client that a new user has join
type ChannelJoin struct {
	UserID    string
	Username  string
	ChannelID string
	Timestamp int64
}

//ChannelCreateRequest represent a message request that was sent to create a channel
type ChannelCreateRequest struct {
	ChannelName string
}

//ChannelCreateResponse represent a message response to create a channel
type ChannelCreateResponse struct {
	ChannelName string
	ChannelID   string
	Username    string
	UserID      string
	IsGame      bool
	Timestamp   int64
}

//ChannelLeaveResponse represent a message response to quit a channel
type ChannelLeaveResponse struct {
	UserID    string
	Username  string
	ChannelID string
	Timestamp int64
}

//ChannelDestroyResponse represent a message response to a channel destroy
type ChannelDestroyResponse struct {
	UserID    string
	Username  string
	ChannelID string
	Timestamp int64
}
