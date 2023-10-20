package MongoSchema

import (
	"fmt"
	pro "go_grpc_rest/protoPackage"
)

type UserData struct {
	Name         string `bson:"name,omitempty"`
	Age          int32  `bson:"age,omitempty"`
	Gender       string `bson:"gender,omitempty"`
	MobileNumber int32  `bson:"mobile_number,omitempty"`
	EmailId      string `bson:"email_id"`
}

// ConvertToSchema converts proto struct to User.
func (w *UserData) ConvertToMongo(user *pro.User) {
	fmt.Println("converting to schema in schema.go")
	w.Name = user.GetName()
	w.Age = user.GetAge()
	w.Gender = user.GetGender()
	w.MobileNumber = user.GetMobileNumber()
	w.EmailId = user.GetEmailId()
}

// ConvertToProto converts User struct into proto
func (w *UserData) ConvertToProto() *pro.User {
	fmt.Println("converting to proto in schema.go")
	return &pro.User{
		Name:         w.Name,
		Age:          w.Age,
		Gender:       w.Gender,
		MobileNumber: w.MobileNumber,
		EmailId:      w.EmailId,
	}
}
