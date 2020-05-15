A project for New York Blockchain Week

<b>Purpose:</b> Empower merchants to accept USD payments without paying merchant services fees utilizing.

There are two stakeholders we need to build for, merchants and buyers.  Currently, we are participating in the New York Blockchain Week hackathon and our focus is small merchants accepting payments from customers that have Stellar-Blockchain-based wallet.  Later our focus will expand to all merchants accepting payments from Bitcoin holders and customers without cryptocurrency.

<b>Goal for Hackathon:</b>

<b>1.</b>  Ability for a merchant to quickly create a QR code that SEP-7 compliant Stellar wallets can read.  This QR code should include amount, currency, and payment address.

<b>2.</b>  Ability for a merchant to withdraw USD tokens from their Stellar wallet to their traditional bank account.  Only during the KYC process should the user leave the Chicken-with-Pasta UX.

<b>Solution:</b> To rapidly meet these two requirements, we have chosen to build a bot on Keybase.  On Keybase, your username is a Stellar wallet address.  Keybase has also implemented Stellar-based path payments.  This allows the customer to pay in any currency while the merchant receives USD.

<b>Use Case</b>

<b>1.</b> Vendor generates QR code for a bill
<b>2.</b> Customer scans the QR code on their mobile Stellar-based wallet and clicks pay
<b>3.</b> Vendor receives USD tokens
<b>4.</b> Vendor withdraws USD tokens
