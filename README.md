# kbbi-api

[![godoc](https://godoc.org/github.com/raf555/kbbi-api/pkg/kbbi?status.svg)][godoc]

[godoc]: https://godoc.org/github.com/raf555/kbbi-api/pkg/kbbi

Probably the most complete public KBBI API you will ever find. 


## Documentation

- Swagger docs: [https://kbbi.raf555.dev/swagger/index.html](https://kbbi.raf555.dev/swagger/index.html)

- Sample API endpoint: [https://kbbi.raf555.dev/api/v1/entry/bermalas-malasan](https://kbbi.raf555.dev/api/v1/entry/bermalas-malasan)

- Sample response:

```json
{
  "lemma": "bermalas-malasan",
  "entries": [
    {
      "entry": "ber.ma.las-ma.las.an",
      "baseWord": "malas",
      "entryVariants": [],
      "pronunciation": "",
      "definitions": [
        {
          "definition": "bermalas-malas",
          "referencedLemma": "",
          "labels": [
            {
              "code": "v",
              "name": "Verba",
              "kind": "Kelas Kata"
            }
          ],
          "usageExamples": [
            "dia selalu ~ saat hari libur"
          ]
        }
      ],
      "nonStandardWords": [],
      "variants": [],
      "compoundWords": [],
      "derivedWords": [],
      "proverbs": [],
      "metaphors": []
    }
  ]
}
```


## Data Source

Latest edition: `Oktober 2025`

The dictionary is mirrored from [Official KBBI Application][] `v6.1.0` with some hand-edited data and customly decoded for author's requirement.

The dictionary used by the server will be updated as soon as new version of the application is released.

## Issues

For any issues, be it from the API server or from the dictionary, or questions or inquiries or suggestions, feel free to raise a new issue. I'll try to look into it as soon as possible.

If you are willing to open PR to fix any open issues, I'm more than welcome to review it.

## Run the Application

The dictionary assets are encrypted (for now). You'll need the assets to run the server, therefore provided samples in the [assets/sample](assets/sample) directory.

You can rename the [.env.sample](.env.sample) file to be `.env` and run the server directly with this command.

```sh
go run ./cmd/kbbi
```

To regenerate the swagger, run this command.

```sh
go generate ./...
```

Once success, you should be able to open http://localhost:8888 in your browser.

[Official KBBI Website]: https://kbbi.kemdikbud.go.id/

[Official KBBI Application]: https://play.google.com/store/apps/details?id=yuku.kbbi5

## Copyright and Data Ownership

> [!WARNING]  
> All data in this dictionary is fully owned by **Badan Pengembangan dan Pembinaan Bahasa, Kementerian Pendidikan, Kebudayaan, Riset, dan Teknologi Republik Indonesia**.
> 
> The commercial use of this dataset is strictly prohibited and subject to legal provisions under: **Indonesian Law No. 28 of 2014 on Copyright**.
>
> The accompanying code is licensed separately and is not subject to the same restrictions as the dictionary data. The code is provided under MIT License, allowing its use, modification, and distribution according to the terms of that license.
