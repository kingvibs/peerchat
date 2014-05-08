package dht

import "fmt"
import "time"
import "os"
import "encoding/gob"

type User struct {
	node *DhtNode
	name string
	// messageHistory map[string]string
}

const PEERCHAT_USERDATA_DIR = "/tmp/"

func usernameToPath(username string) string {
	/*
		Given a username, returns the filepath
		where the backup will be located.
		
		NOTE: Go treats "/" as the path separator on
		all platforms.
	*/
	return PEERCHAT_USERDATA_DIR + "/" + username + ".gob"
}

func (u *User) Serialize() {
	/*
		Serializes this User struct.
	*/
	path := usernameToPath(u.name)
	encodeFile, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	// encode and write to file
	encoder := gob.NewEncoder(encodeFile)
	if err := encoder.Encode(u); err != nil {
		panic(err)
	}
	encodeFile.Close()
}

func Deserialize(username string) *User {
	/*
		Deserializes a DhtNode and loads
		it into a new DhtNode, which is 
		returned. 
	*/
	path := usernameToPath(username)
	decodeFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer decodeFile.Close()

	// create decoder
	decoder := gob.NewDecoder(decodeFile)
	newUser := new(User)
	decoder.Decode(&newUser)
	return newUser
}

func loadTable(username, myIpAddr string) *User {
	/*
		This method loads the User struct for a given 
		username from disk, checking if there needs to 
		be a reconfiguration of the routing table or not
		and acting appropriately. 
	*/
	// first deserialize the old User struct from disk
	user := Deserialize(username)
	
	// check and see if ipaddr is the same as the old one
	// if so, we don't need to change anything
	if user.node.IpAddr != myIpAddr {
		
		// otherwise, create a new nodeId
		user.node.NodeId = Sha1(myIpAddr)
		
		// and rearrange the table based on new nodeId
		// first, get a list of all (nodeId, ipAddr) pairs
		routingEntries := make([]RoutingEntry, 0)
		
		// for each k-buckets row
		for _, row := range user.node.RoutingTable {
		
			// for each RoutingEntry in row
			for _, entry := range row {
				routingEntries = append(routingEntries, entry)
			}
		}
		
		// now delete old routing table and replace 
		// with a new (empty) one
		var routingTable [IDLen][]RoutingEntry
		for i, _ := range routingTable {
			routingTable[i] = make([]RoutingEntry, 0)
		}
		user.node.RoutingTable = routingTable
		
		// then, for each pair, call:
		// updateRoutingTable(nodeId ID, IpAddr string)
		for _, entry := range routingEntries {
			if user.node.Ping(entry){
				user.node.updateRoutingTable(entry)
			}			
		}
	}
	
	return user
}

//SendMessage RPC Handler
func (user *User) SendMessageHandler(args *SendMessageArgs, reply *SendMessageReply) error{
	fmt.Println("My Message:", args.Content)
	return nil
}

//SendMessage API
func (user *User) SendMessage(username string, content string){
	ipAddr := user.node.GetUser(username)
	switch user.CheckStatus(ipAddr){
		case Online:
			args := &SendMessageArgs{Content: content, ToUsername: username, FromUsername: user.name, Timestamp: time.Now()}
			var reply SendMessageReply
			call(ipAddr, "User.SendMessageHandler", args, &reply)
		case Offline:
			// if not, queue?
	}
}

func (user *User) CheckStatus(ipAddr string) Status {
	return Online
}

func Login(username string, userIpAddr string) *User {
	user := loadTable(username, userIpAddr)
	return user
}

