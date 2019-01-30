package main

import (
	"context"
	"flag"
	"io"
	"log"
	"math/rand"
	"time"

	pb "../route_guide"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The TLS cert file")
	serverAddr         = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format ip:port")
	serverHostOverride = flag.String("server_host_override", "example.com", "the server name use to verify the hostname returned by TLS handshake")
)

// printFeature get the feature fo the given point
func printFeature(client pb.RouteGuideClient, point *pb.Point) {
	log.Printf("Getting the feature for point(%d, %d)", point.Latitude, point.Longitude)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	feature, err := client.GetFeature(ctx, point)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v", client, err)
	}
	log.Println(feature)
}

// printFeatures list all the features within a given bounding rectangle.
func printFeatures(client pb.RouteGuideClient, rect *pb.Rectangle) {
	log.Printf("Looking for feature within:%v", rect)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ListFeatures(ctx, rect)
	if err != nil {
		log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
	}
	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
		}
		log.Println(feature)
	}

}

// runRecordRoute sends a sequence of ppooints to server and expects to get a RouteSummary from server
func runRecordRoute(client pb.RouteGuideClient) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pointCount := int(r.Int31n(100)) + 2 // Traverse at least 2 points
	var points []*pb.Point

	for i := 0; i < pointCount; i++ {
		points = append(points, randomPoint(r))
	}
	log.Printf("Traversing %d points", len(points))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.RecordRoute(ctx)
	if err != nil {
		log.Fatalf("%v.RecordRoute(_) = _, %v", client, err)
	}
	for _, point := range points {
		if err := stream.Send(point); err != nil {
			log.Fatalf("%v.Send(%v) = _, %v", stream, point, err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v,  want %v", stream, err, nil)
	}
	log.Printf("Route summary: %v", reply)
}

// runRouteChat receives a sequence of route notes, while sending notes for various locations
func runRouteChat(client pb.RouteGuideClient) {
	notes := []*pb.RouteNote{
		{Location: &pb.Point{Latitude: 0, Longitude: 1}, Message: "First Message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 2}, Message: "second Message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 3}, Message: "third Message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 1}, Message: "fourth Message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 2}, Message: "fifth Message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 3}, Message: "sixth Message"},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.RouteChat(ctx)
	if err != nil {
		log.Fatalf("%v.RouteChat(_) = _, %v", client, err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				//read done
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note: %v", err)
			}
			log.Printf("Got message %s at point (%d %d)", in.Message, in.Location.Latitude, in.Location.Longitude)
		}
	}()
	for _, note := range notes {
		if err := stream.Send(note); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}
	}
	stream.CloseSend()
	<-waitc
}

func randomPoint(r *rand.Rand) *pb.Point {
	lat := (r.Int31n(180) - 90) * 1e7
	lon := (r.Int31n(360) - 180) * 1e7
	return &pb.Point{Latitude: lat, Longitude: lon}
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create tls credentials: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("Could not connect to server:%v", err)
	}
	defer conn.Close()
	client := pb.NewRouteGuideClient(conn)

	//Looking for a valid feature
	printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})

	//Looking for a invalid feature
	printFeature(client, &pb.Point{Latitude: 0, Longitude: 0})

	printFeatures(client, &pb.Rectangle{
		Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	})

	// RecordRoute
	runRecordRoute(client)

	// RouteChat
	runRouteChat(client)

}
