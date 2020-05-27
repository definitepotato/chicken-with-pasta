package anchorusd

type apiCall struct {
  header        string
  transaction   string
  asset_code    string
  account       string
  amount        string
  email_address string
}

func MakeWithdrawUrl(email string, account string) string {
  url := apiCall {
    header:         "https://sandbox-api.anchorusd.com/transfer/withdraw",
    transaction:    "?type="           + "bank_account",
    asset_code:     "&asset_code="     + "USD",
    account:        "&account="        + account,
    email_address:  "&email_address="  + email,
  }
  return url.header + url.transaction + url.asset_code + url.email_address + url.account
}
