package ms

import (
	"context"
	"crypto/tls"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
	"time"
)

var log = logf.Log.WithName("mongo")

type Config struct {
	Hosts       []string
	ReplSetName string
	Username    string
	Password    string
	TLSConf     *tls.Config
}

func Dial(conf *Config) (*mongo.Client, error) {
	ctx, connectcancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer connectcancel()

	opts := options.Client().
		SetHosts(conf.Hosts).
		SetReplicaSet(conf.ReplSetName).
		SetAuth(options.Credential{
			Password: conf.Password,
			Username: conf.Username,
		}).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority(), writeconcern.J(true))).
		SetReadPreference(readpref.Primary()).SetTLSConfig(conf.TLSConf)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, errors.Errorf("failed to connect to mongo rs: %v", err)
	}

	defer func() {
		if err != nil {
			derr := client.Disconnect(ctx)
			if derr != nil {
				log.Error(err, "failed to disconnect")
			}
		}
	}()

	ctx, pingcancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer pingcancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, errors.Errorf("failed to ping mongo: %v", err)
	}

	return client, nil
}

func PingMongoClient(urls string) error {
	urls = strings.Replace(urls, "mongodb://", "", -1)
	urlsList := strings.Split(urls, ";")
	for _, url := range urlsList {
		exp := strings.Split(url, "@")
		host := strings.Split(exp[1], ",")
		account := strings.Split(exp[0], ":")
		username := account[0]
		password := account[1]
		_, err := Dial(&Config{
			Hosts:    host,
			Username: username,
			Password: password,
		})
		if err != nil {
			return errors.Wrapf(err, "failed to ping mongodb host:%s", host)
		}
	}
	return nil
}
