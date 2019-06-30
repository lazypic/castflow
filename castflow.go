package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	now = time.Now()
	// db setting
	flagRegion  = flag.String("region", "ap-northeast-2", "AWS region name")
	flagProfile = flag.String("profile", "lazypic", "AWS Credentials profile name")
	flagTable   = flag.String("table", "castflow", "AWS Dynamodb table name")

	// mode and partition key
	flagAdd = flag.Bool("add", false, "user addition mode")
	flagSet = flag.Bool("set", false, "user update mode")
	flagRm  = flag.Bool("rm", false, "user remove mode")

	// date
	flagHelp = flag.Bool("help", false, "print help")

	// attributes
	flagID              = flag.String("id", "", "character name")
	flagRegnum          = flag.String("regnum", "", "registration number")
	flagManager         = flag.String("manager", "", "project start date")
	flagFieldOfActivity = flag.String("foa", "", "field of activity")
	flagConcept         = flag.String("concept", "", "character concept")
	flagStartDate       = flag.String("start", now.Format(time.RFC3339), "project end date")
	flagEmail           = flag.String("email", "", "character e-mail")
	flagSearchword      = flag.String("search", "", "search word")
)

func main() {
	log.SetPrefix("castflow: ")
	log.SetFlags(0)
	flag.Parse()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String(*flagRegion)},
		Profile:           *flagProfile,
	}))
	db := dynamodb.New(sess)

	// 테이블이 존재하는지 점검하고 없다면 테이블을 생성한다.
	if !validTable(*db, *flagTable) {
		_, err := db.CreateTable(tableStruct(*flagTable))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Println("Created table:", *flagTable)
		fmt.Println("Please try again in one minute.")
		os.Exit(0)
	}
	if *flagHelp {
		flag.Usage()
	}
	if *flagAdd && *flagID != "" {
		err := AddCharacter(*db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	} else if *flagSet && *flagID != "" {
		err := SetCharacter(*db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	} else if *flagRm && *flagID != "" {
		user, err := user.Current()
		if user.Username != "root" {
			log.Fatal(errors.New("사용자를 삭제하기 위해서는 root 권한이 필요합니다"))
		}
		err = RmCharacter(*db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	} else if *flagSearchword != "" {
		err := GetCharacters(*db, *flagSearchword)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	} else {
		flag.PrintDefaults()
	}
}
