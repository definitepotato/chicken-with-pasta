package main

import (
  "math/rand"
  "net/http"
  "strconv"
  "strings"
  "flag"
  "fmt"
  "os"

  "github.com/definitepotato/chicken-with-pasta/lib/api/anchorusd"
  "github.com/definitepotato/chicken-with-pasta/lib/stellar/sep0007"
  "github.com/keybase/go-keybase-chat-bot/kbchat"
  qrcode "github.com/skip2/go-qrcode"
)

func fail(msg string, args ...interface{}) {
  fmt.Fprintf(os.Stderr, msg+"\n", args...)
  os.Exit(3)
}

type messageBody struct {
  recipient     string
  command       string
  amount        string
  assetCode     string
}

func makeQrCode(payUrl string) string {
  qrName := strconv.Itoa(rand.Int()) + ".png"
  qrcode.WriteFile(payUrl, qrcode.Medium, 256, qrName)
  return qrName
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

    // WE NEED TO SANITIZE USER INPUT !!
    splitMsg := strings.Split(msg.Message.Content.Text.Body, " ")
    switch splitMsg[0] {
    case "!payme":
      // We get an index out of range issue when the user provides
      // a message that doesn't include at least 2 spaces (3 values)
      // since we force push at minimum 3 expected values into a
      // struct (sortedMsg).  This is a cheap way to bypass the crashing
      // issue and avoid the bot going offline.  We need to validate
      // input anywhere we pass values directly from the user.
      if len(splitMsg) < 3 {
        continue
      }

      sortedMsg := messageBody {
        recipient:    msg.Message.Sender.Username + "*keybase.io",
        command:      splitMsg[0],
        amount:       splitMsg[1],
        assetCode:    strings.ToUpper(splitMsg[2]),
      }

      if assetIssuer, ok := assetCode[sortedMsg.assetCode]; ok {

        payUrl := sep0007.MakePayUrl(sortedMsg.recipient, sortedMsg.amount, sortedMsg.assetCode, assetIssuer)
        payMsg := makeQrCode(payUrl)

        if _, err = kbc.SendAttachmentByConvID(msg.Message.ConvID, payMsg, "Scan This QR Code To Pay"); err != nil {
          fail("Error sending attachment: %s", err.Error())
        }
      }
    case "!withdraw":
      // Need to pull email and anchor public address from !register db.
      withdrawUrl := anchorusd.MakeWithdrawUrl("", "")

      response, err := http.Get(withdrawUrl)
      if err != nil {
        fail("Error during HTTP GET: %s", err.Error())
      }
      fmt.Println(response)
      response.Body.Close()

      if _, err = kbc.SendMessage(msg.Message.Channel, "Withdrawl Request Sent!"); err != nil {
        fail("Error echo'ing message: %s", err.Error())
      }
    case "!register":
      if _, err = kbc.SendMessage(msg.Message.Channel, "Coming soon ..."); err != nil {
        fail("Error echo'ing message: %s", err.Error())
      }
    default:
      if _, err = kbc.SendMessage(msg.Message.Channel, "Help!"); err != nil {
        fail("Error echo'ing message: %s", err.Error())
      }
    }

  }
}
