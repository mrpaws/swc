package lib

import (
	"log"
	"os"
)

type SWCVerifyToken struct {
	value string
	/*
		subscriptions []string
		accessed string
		modified stringf
		created string
	*/
}

func (SWCVerifyToken *SWCVerifyToken) String() string {
	return SWCVerifyToken.value
}

func GetSwcVerifyToken() SWCVerifyToken {
	token := os.Getenv("swcVerifyToken")
	if token == "" {
		log.Fatal("swcVerifyToken env var not found")
	} else {
		log.Print("swcVerifyToken set from env var (hidden)")
	}
	return SWCVerifyToken{value: token}
}

func GetSwcSubscriptionId() string {
	id := os.Getenv("swcSubscriptionId")
	if id == "" {
		log.Print("swcSubscriptionId not found in env; none set")
	} else {
		log.Printf("swcSubscriptionId set from env: %s", id)
	}
	return id
}

func LoadEnvVar(env string, required bool) string {
    log.Printf("env var lookup: (%s), required: (%v)", env, required)
	val, found := os.LookupEnv(env)
	if found == false {
		log.Fatalf("required env var does not exist")
	} else if val == ""{
		log.Printf("empty env var detected: %s", env)
	}
	return val
}