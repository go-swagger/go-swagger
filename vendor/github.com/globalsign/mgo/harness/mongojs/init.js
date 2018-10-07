//var settings = {heartbeatSleep: 0.05, heartbeatTimeout: 0.5}
var settings = {};

// We know the master of the first set (pri=1), but not of the second.
var rs1cfg = {
    _id: "rs1",
    members: [{ _id: 1, host: "127.0.0.1:40011", priority: 1, tags: { rs1: "a" } },
    { _id: 2, host: "127.0.0.1:40012", priority: 0, tags: { rs1: "b" } },
    { _id: 3, host: "127.0.0.1:40013", priority: 0, tags: { rs1: "c" } }],
    settings: settings
}
var rs2cfg = {
    _id: "rs2",
    members: [{ _id: 1, host: "127.0.0.1:40021", priority: 1, tags: { rs2: "a" } },
    { _id: 2, host: "127.0.0.1:40022", priority: 1, tags: { rs2: "b" } },
    { _id: 3, host: "127.0.0.1:40023", priority: 1, tags: { rs2: "c" } }],
    settings: settings
}
var rs3cfg = {
    _id: "rs3",
    members: [{ _id: 1, host: "127.0.0.1:40031", priority: 1, tags: { rs3: "a" } },
    { _id: 2, host: "127.0.0.1:40032", priority: 0, tags: { rs3: "b" } },
    { _id: 3, host: "127.0.0.1:40033", priority: 0, tags: { rs3: "c" } }],
    settings: settings
}

for (var i = 0; i != 60; i++) {
    try {
        db1 = new Mongo("127.0.0.1:40001").getDB("admin")
        rs1a = new Mongo("127.0.0.1:40011").getDB("admin")
        rs2a = new Mongo("127.0.0.1:40021").getDB("admin")
        rs3a = new Mongo("127.0.0.1:40031").getDB("admin")
        cfg1 = new Mongo("127.0.0.1:40101").getDB("admin")
        cfg2 = new Mongo("127.0.0.1:40102").getDB("admin")
        cfg3 = new Mongo("127.0.0.1:40103").getDB("admin")
        break
    } catch (err) {
        print("Can't connect yet...")
    }
    sleep(1000)
}

function hasSSL() {
    return Boolean(db1.serverBuildInfo().OpenSSLVersion)
}

function versionAtLeast() {
    var version = db1.version().split(".")
    for (var i = 0; i < arguments.length; i++) {
        if (i == arguments.length) {
            return false
        }
        if (arguments[i] != version[i]) {
            return version[i] >= arguments[i]
        }
    }
    return true
}


if (versionAtLeast(3, 4)) {
    print("configuring config server for mongodb 3.4+")
    cfg1.runCommand({ replSetInitiate: { _id: "conf1", configsvr: true, members: [{ "_id": 1, "host": "localhost:40101" }] } })
    cfg2.runCommand({ replSetInitiate: { _id: "conf2", configsvr: true, members: [{ "_id": 1, "host": "localhost:40102" }] } })
    cfg3.runCommand({ replSetInitiate: { _id: "conf3", configsvr: true, members: [{ "_id": 1, "host": "localhost:40103" }] } })
}

sleep(3000)

rs1a.runCommand({ replSetInitiate: rs1cfg })
rs2a.runCommand({ replSetInitiate: rs2cfg })
rs3a.runCommand({ replSetInitiate: rs3cfg })

function configAuth() {
    var addrs = ["127.0.0.1:40002", "127.0.0.1:40203", "127.0.0.1:40031"]
    if (hasSSL()) {
        addrs.push("127.0.0.1:40003")
    }
    for (var i in addrs) {
        print("Configuring auth for", addrs[i])
        var db = new Mongo(addrs[i]).getDB("admin")
        var timedOut = false
        createUser:
        for (var i = 0; i < 60; i++) {
            try {
                db.createUser({ user: "root", pwd: "rapadura", roles: ["root"] })
            } catch (err) {
                // 3.2 consistently fails replication of creds on 40031 (config server) 
                print("createUser command returned an error: " + err)
                if (String(err).indexOf("timed out") >= 0) {
                    timedOut = true;
                }
                // on 3.6 cluster with keyFile, we sometimes get this error 
                if (String(err).indexOf("Cache Reader No keys found for HMAC that is valid for time")) {
                    sleep(500)
                    continue createUser;
                }
            }
            break;
        }

        for (var i = 0; i < 60; i++) {
            var ok = db.auth("root", "rapadura")
            if (ok || !timedOut) {
                break
            }
            sleep(1000);
        }
        sleep(500)
        db.createUser({ user: "reader", pwd: "rapadura", roles: ["readAnyDatabase"] })
        sleep(3000)
    }
}

function addShard(adminDb, shardList) {
    for (var index = 0; index < shardList.length; index++) {
        for (var i = 0; i < 10; i++) {
            var result = adminDb.runCommand({ addshard: shardList[index] })
            if (result.ok == 1) {
                print("shard " + shardList[index] + " sucessfully added")
                break
            } else {
                print("fail to add shard: " + shardList[index] + " error: " + JSON.stringify(result) + ", retrying in 1s")
                sleep(1000)
            }
        }
    }
}

function configShards() {
    s1 = new Mongo("127.0.0.1:40201").getDB("admin")
    addShard(s1, ["127.0.0.1:40001", "rs1/127.0.0.1:40011"])

    s2 = new Mongo("127.0.0.1:40202").getDB("admin")
    addShard(s2, ["rs2/127.0.0.1:40021"])

    s3 = new Mongo("127.0.0.1:40203").getDB("admin")
    for (var i = 0; i < 10; i++) {
        var ok = s3.auth("root", "rapadura")
        if (ok) {
            break
        }
        sleep(1000)
    }
    addShard(s3, ["rs3/127.0.0.1:40031"])
}

function countHealthy(rs) {
    var status = rs.runCommand({ replSetGetStatus: 1 })
    var count = 0
    var primary = 0
    if (typeof status.members != "undefined") {
        for (var i = 0; i != status.members.length; i++) {
            var m = status.members[i]
            if (m.health == 1 && (m.state == 1 || m.state == 2)) {
                count += 1
                if (m.state == 1) {
                    primary = 1
                }
            }
        }
    }
    if (primary == 0) {
        count = 0
    }
    return count
}

var totalRSMembers = rs1cfg.members.length + rs2cfg.members.length + rs3cfg.members.length

for (var i = 0; i != 60; i++) {
    var count = countHealthy(rs1a) + countHealthy(rs2a) + countHealthy(rs3a)
    print("Replica sets have", count, "healthy nodes.")
    if (count == totalRSMembers) {
        configAuth()
        sleep(2000)
        configShards()
        quit(0)
    }
    sleep(1000)
}

print("Replica sets didn't sync up properly.")
quit(12)

// vim:ts=4:sw=4:et
