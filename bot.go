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

// Function to gracefully fail and provide user with a useful message.
func fail(msg string, args ...interface{}) {
  fmt.Fprintf(os.Stderr, msg+"\n", args...)
  os.Exit(3)
}

func main() {
  var kbLoc string
  var kbc *kbchat.API
  var err error

  flag.StringVar(&kbLoc, "keybase", "keybase", "the location of the Keybase app")
  flag.Parse()

  // Initiate keybase and establish API connection.
  if kbc, err = kbchat.Start(kbchat.RunOptions{KeybaseLocation: kbLoc}); err != nil {
    fail("Error creating API: %s", err.Error())
  }

  sub, err := kbc.ListenForNewTextMessages()
  if err != nil {
      fail("Error listening: %s", err.Error())
  }

  // Hardcoded map of Asset Code to Issuer Address.
  var assetCode map[string]string
  assetCode = make(map[string]string)
  assetCode["USD"] = "GDUKMGUGDZQK6YHYA5Z6AY2G4XDSZPSZ3SW5UN3ARVMO6QSRDWP5YLEX"
  assetCode["BTC"] = "GAUTUYY2THLF7SGITDFMXJVYH3LHDSMGEAKSBU267M2K7A3W543CKUEF"
  assetCode["ETH"] = "GBDEVU63Y6NTHJQQZIKVTC23NWLQVP3WJ2RI2OTSJTNYOIGICST6DUXR"

  // Begin loop to listen for messages in the comm channel.
  for {
    msg, err := sub.Read()
    if err != nil{
      fail("Failed to read message: %s", err.Error())
    }

    // Split received message by spaces into elements in a slice.
    s := strings.Split(msg.Message.Content.Text.Body, " ")

    // Ignore message unless it begins with !pay.
    if s[0] == "!payme" {

      // TODO: Validate Stellar Address.
      // TODO: Validate Asset Code.
      // TODO: Map Asset Code to default list.
      // TODO: User specified Asset Issuer (end of message).
      // TODO: Validate Asset Issuer if provided.

      if msg.Message.Content.TypeName != "text" {
        continue
      }

      if msg.Message.Sender.Username == kbc.GetUsername() {
        continue
      }

      // Check that we have a hardcoded Issuer for the provided Asset Code. This is sloppy shit, we can do better.
      s[2] = strings.ToUpper(s[2])
      if ac, ok := assetCode[s[2]]; ok {

        // Build the user federated address based on the message sender.
        recipient := msg.Message.Sender.Username + "*keybase.io"

        // Build Stellar Pay URL.
        spurl := "web+stellar:pay?destination=" + recipient + "&amount=" + s[1] + "&asset_code=" + s[2] + "&asset_issuer=" + ac

        // Build QR Code. More hardcoded shit, need to scrap this later.
        qrName := strconv.Itoa(rand.Int()) + ".png"
        fileLoc := "/var/www/html/" + qrName
        payMsg := "Pay at: https://chickenwithpasta.com/" + qrName
        err := qrcode.WriteFile(spurl, qrcode.Medium, 256, qrName)
        err = os.Rename(qrName, fileLoc)

        if _, err = kbc.SendMessage(msg.Message.Channel, payMsg); err != nil {
          fail("Error echo'ing message: %s", err.Error())
        }
      }
    }

  }
}
