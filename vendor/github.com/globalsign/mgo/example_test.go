package mgo

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"sync"
)

func ExampleCredential_x509Authentication() {
	// MongoDB follows RFC2253 for the ordering of the DN - if the order is
	// incorrect when creating the user in Mongo, the client will not be able to
	// connect.
	//
	// The best way to generate the DN with the correct ordering is with
	// openssl:
	//
	// 		openssl x509 -in client.crt -inform PEM -noout -subject -nameopt RFC2253
	// 		subject= CN=Example App,OU=MongoDB Client Authentication,O=GlobalSign,C=GB
	//
	//
	// And then create the user in MongoDB with the above DN:
	//
	//		db.getSiblingDB("$external").runCommand({
	//			createUser: "CN=Example App,OU=MongoDB Client Authentication,O=GlobalSign,C=GB",
	//			roles: [
	//				{ role: 'readWrite', db: 'bananas' },
	//				{ role: 'userAdminAnyDatabase', db: 'admin' }
	//			],
	//			writeConcern: { w: "majority" , wtimeout: 5000 }
	//		})
	//
	//
	// References:
	// 		- https://docs.mongodb.com/manual/tutorial/configure-x509-client-authentication/
	// 		- https://docs.mongodb.com/manual/core/security-x.509/
	//

	// Read in the PEM encoded X509 certificate.
	//
	// See the client.pem file at the path below.
	clientCertPEM, err := ioutil.ReadFile("harness/certs/client.pem")

	// Read in the PEM encoded private key.
	clientKeyPEM, err := ioutil.ReadFile("harness/certs/client.key")

	// Parse the private key, and the public key contained within the
	// certificate.
	clientCert, err := tls.X509KeyPair(clientCertPEM, clientKeyPEM)

	// Parse the actual certificate data
	clientCert.Leaf, err = x509.ParseCertificate(clientCert.Certificate[0])

	// Use the cert to set up a TLS connection to Mongo
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},

		// This is set to true so the example works within the test
		// environment.
		//
		// DO NOT set InsecureSkipVerify to true in a production
		// environment - if you use an untrusted CA/have your own, load
		// its certificate into the RootCAs value instead.
		//
		// RootCAs: myCAChain,
		InsecureSkipVerify: true,
	}

	// Connect to Mongo using TLS
	host := "localhost:40003"
	session, err := DialWithInfo(&DialInfo{
		Addrs: []string{host},
		DialServer: func(addr *ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", host, tlsConfig)
		},
	})

	// Authenticate using the certificate
	cred := &Credential{Certificate: tlsConfig.Certificates[0].Leaf}
	if err := session.Login(cred); err != nil {
		panic(err)
	}

	// Done! Use mgo as normal from here.
	//
	// You should actually check the error code at each step.
	_ = err
}

func ExampleSession_concurrency() {
	// This example shows the best practise for concurrent use of a mgo session.
	//
	// Internally mgo maintains a connection pool, dialling new connections as
	// required.
	//
	// Some general suggestions:
	// 		- Define a struct holding the original session, database name and
	// 			collection name instead of passing them explicitly.
	// 		- Define an interface abstracting your data access instead of exposing
	// 			mgo to your application code directly.
	// 		- Limit concurrency at the application level, not with SetPoolLimit().

	// This will be our concurrent worker
	var doStuff = func(wg *sync.WaitGroup, session *Session) {
		defer wg.Done()

		// Copy the session - if needed this will dial a new connection which
		// can later be reused.
		//
		// Calling close returns the connection to the pool.
		conn := session.Copy()
		defer conn.Close()

		// Do something(s) with the connection
		_, _ = conn.DB("").C("my_data").Count()
	}

	///////////////////////////////////////////////

	// Dial a connection to Mongo - this creates the connection pool
	session, err := Dial("localhost:40003/my_database")
	if err != nil {
		panic(err)
	}

	// Concurrently do things, passing the session to the worker
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go doStuff(wg, session)
	}
	wg.Wait()

	session.Close()
}

func ExampleDial_usingSSL() {
	// To connect via TLS/SSL (enforced for MongoDB Atlas for example) requires
	// to set the ssl query param to true.
	url := "mongodb://localhost:40003?ssl=true"

	session, err := Dial(url)
	if err != nil {
		panic(err)
	}

	// Use session as normal
	session.Close()
}

func ExampleDial_tlsConfig() {
	// You can define a custom tlsConfig, this one enables TLS, like if you have
	// ssl=true in the connection string.
	url := "mongodb://localhost:40003"

	tlsConfig := &tls.Config{
		// This can be configured to use a private root CA - see the Credential
		// x509 Authentication example.
		//
		// Please don't set InsecureSkipVerify to true - it makes using TLS
		// pointless and is never the right answer!
	}

	dialInfo, err := ParseURL(url)
	dialInfo.DialServer = func(addr *ServerAddr) (net.Conn, error) {
		return tls.Dial("tcp", addr.String(), tlsConfig)
	}

	session, err := DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}

	// Use session as normal
	session.Close()
}
