package sep0007

type s0007 struct {
  header        string
  destination   string
  amount        string
  assetCode     string
  assetIssuer   string
}

func MakePayUrl(recipient string, amount string, assetcode string, assetissuer string) string {
  url := s0007 {
    header:       "web+stellar:pay",
    destination:  "?destination="   + recipient,
    amount:       "&amount="        + amount,
    assetCode:    "&asset_code="    + assetcode,
    assetIssuer:  "&asset_issuer="  + assetissuer,
  }
  return url.header + url.destination + url.amount + url.assetCode + url.assetIssuer
}
