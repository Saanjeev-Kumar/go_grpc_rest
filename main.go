package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go_grpc_rest/MongoSchema"
	"go_grpc_rest/protoPackage"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"time"
)

type SingleResult interface {
	Decode(v interface{}) error
}
type ServerDB struct {
	store *mongo.Client
}

func (s *ServerDB) CreateUser(ctx context.Context, request *protoPackage.CreateUserRequest) (*protoPackage.CreateUserResponse, error) {
	//TODO implement me
	//panic("implement me")
	fmt.Println("Enter to Func Create & received Response: ", request.GetUser())
	userInfo := &MongoSchema.UserData{}
	userInfo.ConvertToMongo(request.GetUser())
	x := s.store
	fmt.Println("Conversion to mongoSchema: ", userInfo)
	collection := x.Database("admin111").Collection("Users")
	fmt.Println("collection inserted in mongodb.go:", collection)
	_, err := collection.InsertOne(ctx, userInfo)
	if err != nil {
		return nil, err
	} else {
		fmt.Println("Result of post", &protoPackage.CreateUserResponse{User: request.GetUser()})
		return &protoPackage.CreateUserResponse{User: request.GetUser()}, nil
	}
}

func (s *ServerDB) GetUser(ctx context.Context, request *protoPackage.GetUserRequest) (*protoPackage.GetUserResponse, error) {

	value := request.GetName()
	fmt.Println("Received through url", value)
	res, _ := s.getContact(ctx, value)
	return &protoPackage.GetUserResponse{User: res}, nil

}

func (s *ServerDB) getContact(ctx context.Context, value string) (*protoPackage.User, error) {
	x := s.store
	collection := x.Database("admin111").Collection("Users")
	sinResult := collection.FindOne(ctx, bson.M{"name": value})
	var userInfo *MongoSchema.UserData
	if err := SingleResult(sinResult).(*mongo.SingleResult).Decode(&userInfo); err != nil {
		return nil, err
	}
	resp := userInfo.ConvertToProto()
	return resp, nil
}

func (s *ServerDB) UpdateUser(ctx context.Context, request *protoPackage.UpdateUserRequest) (*protoPackage.UpdateUserResponse, error) {
	//TODO implement me
	//panic("implement me")
	var updateParams map[string]interface{}
	userInfo, _ := json.Marshal(request.GetUser())
	json.Unmarshal(userInfo, &updateParams)
	fmt.Println("updateParams: ", updateParams)
	fmt.Println("userInfo", userInfo)
	collection := s.store.Database("admin111").Collection("Users")
	_, err := collection.UpdateOne(ctx, bson.M{"name": request.GetName()}, bson.M{"$set": updateParams})
	if err != nil {
		fmt.Println(err)
	}
	updatedUser, err := s.getContact(ctx, request.GetUser().Name)
	return &protoPackage.UpdateUserResponse{User: updatedUser}, nil
}

func (s *ServerDB) DeleteUser(ctx context.Context, request *protoPackage.DeleteUserRequest) (*protoPackage.DeleteUserResponse, error) {
	//TODO implement me
	//panic("implement me")
	collection := s.store.Database("admin111").Collection("Users")
	_, err := collection.DeleteOne(ctx, bson.M{"name": request.GetName()})
	if err != nil {
		fmt.Println(err)
	}
	return &protoPackage.DeleteUserResponse{Status: "Deleted record successfully"}, nil
}

func main() {
	mongoURI := "mongodb+srv://saanjeev:go9z2wF62WKuDrcD@cluster0.iqret.mongodb.net/"
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Println("error in mongoDb Connection", err)
	}
	defer client.Disconnect(ctx)
	fmt.Println("Connected to MongoDB")

	_, err = net.Listen("tcp", "localhost:8082")
	if err != nil {
		fmt.Println("Failed to listen")
	}
	fmt.Println("net.Listen")

	s := grpc.NewServer()
	protoPackage.RegisterUserServiceServer(s, &ServerDB{store: client})
	fmt.Println("RegisterUserServiceServer")
	/*
		//This way not working but this the ideal way of mdm std.
		lis, err := net.Listen("tcp", ":9981")
		if err != nil {
			fmt.Println("Failed to listen")
		}
		fmt.Println("net.Listen")
		if err = s.Serve(lis); err != nil {
			fmt.Println("Failed to serve", err)
			return
		}
		fmt.Println("Serve")
	*/

	gwMux := runtime.NewServeMux()
	//The below two is the ideal way of mdm std.But not working
	//opts := []grpc.DialOption{grpc.WithInsecure()}
	//host := "localhost"
	//port := 9981
	//getMethod := fmt.Sprintf("%v:%v", host, port)
	//protoPackage.RegisterUserServiceHandlerFromEndpoint(context.Background(), gwMux, getMethod, opts)

	//this way is deprecated.
	protoPackage.RegisterUserServiceHandlerServer(context.Background(), gwMux, &ServerDB{store: client})

	fmt.Println("RegisterUserServiceHandlerFromEndpoint")
	log.Printf("Server listening on localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", gwMux))
}
