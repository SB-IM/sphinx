package livestream

type RTPBroadcastConfigOptions struct {
	MQTTClientConfigOptions
	WebRTCConfigOptions
	RTPSourceConfigOptions
}

type RTSPBroadcastConfigOptions struct {
	MQTTClientConfigOptions
	WebRTCConfigOptions
	RTSPSourceConfigOptions
}

type broadcastConfigOptions struct {
	MQTTClientConfigOptions
	WebRTCConfigOptions
}

type MQTTClientConfigOptions struct {
	OfferTopic        string
	AnswerTopicPrefix string
	Qos               uint
	Retained          bool
}

type WebRTCConfigOptions struct {
	ICEServer  string
	Username   string
	Credential string
}

type RTPSourceConfigOptions struct {
	RTPHost string
	RTPPort int
}

type RTSPSourceConfigOptions struct {
	RTSPAddr string
}
