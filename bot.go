package main

import (
  "math/rand"
  "strconv"
  "strings"
  "flag"
  "fmt"
  "os"

  "github.com/keybase/go-keybase-chat-bot/kbchat"
  qrcode "github.com/skip2/go-qrcode"
)

type SepSeven struct {
  pre           string
  destination   string
  amount        string
  assetCode     string
  assetIssuer   string
}

type MessageBody struct {
  recipient     string
  command       string
  amount        string
  assetCode     string
}

func fail(msg string, args ...interface{}) {
  fmt.Fprintf(os.Stderr, msg+"\n", args...)
  os.Exit(3)
}

func generateQrCode(payUrl string) string {
  qrName := strconv.Itoa(rand.Int()) + ".png"
  fileLoc := "/var/www/html/" + qrName

  qrcode.WriteFile(payUrl, qrcode.Medium, 256, qrName)
  os.Rename(qrName, fileLoc)

  payMsg := "Pay at: https://chickenwithpasta.com/" + qrName
  return payMsg
}

func main() {
  var kbLoc string
  var kbc *kbchat.API
  var err error

  assetCode := map[string]string {
    "USD":  "GDUKMGUGDZQK6YHYA5Z6AY2G4XDSZPSZ3SW5UN3ARVMO6QSRDWP5YLEX",
    "BTC":  "GAUTUYY2THLF7SGITDFMXJVYH3LHDSMGEAKSBU267M2K7A3W543CKUEF",
    "ETH":  "GBDEVU63Y6NTHJQQZIKVTC23NWLQVP3WJ2RI2OTSJTNYOIGICST6DUXR",
  }

  flag.StringVar(&kbLoc, "keybase", "keybase", "the location of the Keybase app")
  flag.Parse()

  if kbc, err = kbchat.Start(kbchat.RunOptions{KeybaseLocation: kbLoc}); err != nil {
    fail("Error creating API: %s", err.Error())
  }

  sub, err := kbc.ListenForNewTextMessages()
  if err != nil {
      fail("Error listening: %s", err.Error())
  }

  for {
    msg, err := sub.Read()
    if err != nil{
      fail("Failed to read message: %s", err.Error())
    }

    if msg.Message.Content.TypeName != "text" {
      continue
    }

    if msg.Message.Sender.Username == kbc.GetUsername() {
      continue
    }

    splitMsg := strings.Split(msg.Message.Content.Text.Body, " ")
    sortedMsg := MessageBody {
      recipient:    msg.Message.Sender.Username + "*keybase.io",
      command:      splitMsg[0],
      amount:       splitMsg[1],
      assetCode:    strings.ToUpper(splitMsg[2]),
    }

    if sortedMsg.command == "!payme" {
      if assetIssuer, ok := assetCode[sortedMsg.assetCode]; ok {

        pay := SepSeven {
          pre:          "web+stellar:pay",
          destination:  "?destination=" + sortedMsg.recipient,
          amount:       "&amount=" + sortedMsg.amount,
          assetCode:    "&asset_code=" + sortedMsg.assetCode,
          assetIssuer:  "&asset_issuer=" + assetIssuer,
        }

        payUrl := pay.pre + pay.destination + pay.amount + pay.assetCode + pay.assetIssuer
        payMsg := generateQrCode(payUrl)

        if _, err = kbc.SendMessage(msg.Message.Channel, payMsg); err != nil {
          fail("Error echo'ing message: %s", err.Error())
        }
      }
    }

  }
}
