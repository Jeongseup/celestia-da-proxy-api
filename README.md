# Celestia DA Proxy API

Descrption...

## Example

```bash
# submit blob
curl -X POST http://localhost:3000/submit_blob -H "Content-Type: application/json" -d '{"data":"SGVsbG8gV29ybGQ="}'

# retrieve blobs
curl -X POST http://localhost:3000/retrieve_blob -H "Content-Type: application/json" -d '{"retrieve_height":"1830482"}'
```
