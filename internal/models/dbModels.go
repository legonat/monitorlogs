package models

 type User struct {
	 Login 		string
	 Password 	[]byte
	 Salt 		[]byte
	 Create_at  int
	 Blocked  	bool
	 Try_count 	int
	 Blocked_at int
	 Deleted 	bool
 }

type RefreshSession struct {
	Login 		string
	Token 		string
	Ua 			string
	Fingerprint string
	Ip 			string
	ExpiresIn 	int64
	CreatedAt 	int64
}