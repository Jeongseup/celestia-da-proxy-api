# Celestia DA Proxy API

Descrption...

## Example

1. Submit Image Blob

```bash
# from local file
curl --location 'https://nftinfo.online/submit_formdata' --form image=@$(pwd)/assets/image/third.webp

# from url image
curl -sL https://raw.githubusercontent.com/hackathemy/celestia-da-proxy-api/main/assets/image/second.webp | curl --location 'https://nftinfo.online/submit_formdata' --form 'image=@-'
```

result

```json
{
  "success": true,
  "result": {
    "hash": "65DC84875F7B9BB96FDCA7B465B4C2CCD62CEEA33EC93F4E16304B09333EB2B5",
    "submitted_height": 1847815
  }
}
```

2. Retrieve Submitted Image from Celestia DA

```bash
# only view in browser, just click the url
https://nftinfo.online/65DC84875F7B9BB96FDCA7B465B4C2CCD62CEEA33EC93F4E16304B09333EB2B5
```

3. Submit Metadata Blob With Namespace Key

```bash
# submit metadata for image
curl -sS -X POST https://nftinfo.online/submit_metadata -H "Content-Type: application/json" -d '{
  "namespace_key": "CelestiaDragonsMetaData",
  "metadata": {
    "description": "Celestia DA Based Dragon NFT Collection",
    "image": "https://nftinfo.online/65DC84875F7B9BB96FDCA7B465B4C2CCD62CEEA33EC93F4E16304B09333EB2B5",
    "name": "Celestia Second DA Dragon",
    "attributes": [
      {
        "trait_type": "Color Palette",
        "value": "Emerald greens and bright yellows"
      },
      {
        "trait_type": "Environment",
        "value": "Fantasy landscape with sparkling rivers"
      },
      {
        "trait_type": "Disposition",
        "value": "Mischievous and lively"
      },
      {
        "trait_type": "Special Feature",
        "value": "Wings that glow like fireflies"
      },
      {
        "trait_type": "Magic Power",
        "value": "Emits a glowing aura"
      }
    ]
  }
}' | jq .
```

result

```json
{
  "success": true,
  "result": {
    "namespace_key": "Q2VsZXN0aW",
    "submitted_height": 1847861,
    "submitted_metadata": {
      "description": "Celestia DA Based Dragon NFT Collection",
      "image": "https://nftinfo.online/65DC84875F7B9BB96FDCA7B465B4C2CCD62CEEA33EC93F4E16304B09333EB2B5",
      "name": "Celestia Second DA Dragon",
      "attributes": [
        {
          "trait_type": "Color Palette",
          "value": "Emerald greens and bright yellows"
        },
        {
          "trait_type": "Environment",
          "value": "Fantasy landscape with sparkling rivers"
        },
        {
          "trait_type": "Disposition",
          "value": "Mischievous and lively"
        },
        {
          "trait_type": "Special Feature",
          "value": "Wings that glow like fireflies"
        },
        {
          "trait_type": "Magic Power",
          "value": "Emits a glowing aura"
        }
      ]
    },
    "submitted_metadata_index": 2
  }
}
```

4. Retrieve Metadata By Namespace Key & Index

```bash
# by above data, you should check two value, "submitted_metadata_index" and "namespace_key"
# example namespace_key: Q2VsZXN0aW
# example submitted_metadata_index: 2
# url format: https://nftinfo.online/<namespace_key>/<submitted_metadata_index>
curl -sS -X GET https://nftinfo.online/Q2VsZXN0aW/2
```

### References

https://jaredwinick.github.io/base64-image-viewer/?ref=tiny-helpers

https://www.based64.xyz/

https://docs.celestia.org/developers/golang-client-tutorial

https://docs.opensea.io/docs/metadata-standards
