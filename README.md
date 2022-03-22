# alfred-gitlab

an alfred workflow to search gitlab projects

## Setup and Usage

1. Go to the releases page and download/import the workflow into alfred.
2. Edit the workflow and make sure the following variables are set:

* `GITLAB_URL` (ex: `https://gitlab.<your-domain>.com`)
* `GITLAB_TOKEN`: Your gitlab access token.  You can generate an access token in your Gitlab Profile.

3. Search for gitlab with `gitlab <search>`

## FAQ

**Q: I get the error `“alfred-gitlab” cannot be opened because the developer cannot be verified.`.  How can i fix it?**

**A:** Thats an error gatekeeper returns when a binary isn't signed by an apple certificate.  To fix that follow these steps:

  1. Press cancel on the promp you received.
  2. Go to preferences and select `Security & Privacy`.  On the general tab, make sure "Allow apps download from: App Store and identified developers" is selected.  Beneath that you should see the text `alfred-gitlab was blocked from use because it is not from an identified developer`.  Press the `Allow Anyway` button and it try using the workflow again.
