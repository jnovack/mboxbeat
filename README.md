# mboxbeat

**mboxbeat** cycles through your mboxen and outputs JSON for each
message, including attachments.

**mboxbeat** is meant to centralize your local mail spools on your linux
servers without forwarding the emails to just another mailbox.

The intent is to send them to some form of noSQL database so they can be
aggregated and queried, both ElasticSearch or MongoDB come to mind.

## Caveats

While **mboxbeat** can process attachments, I heavily recommend against
processing them.  This is not an email-forwarding service, it is meant to
serve as a logging vessel for centralizing, aggregating and querying emails.

To that end, you are more interested in the text and metrics of than you are
with downloading and saving attachments.