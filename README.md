# HelloSign Go SDK
A Go wrapper for the HelloSign API.

**Not Ready For Release**

The official library for using the HelloSign API for golang.

## Usage

```go
go get github.com/HelloFax/hellosign-go-sdk
```

Create a client:

```go
client := hellosign.Client{APIKey: "ACCOUNT API KEY HERE"}
```

### Signature Request Methods

#### Create Embedded Signature Request Using Files

```go
  fileOne, _ := os.Open("public/offer_letter.pdf")
  defer fileOne.Close()
  fileTwo, _ := os.Open("public/offer_letter.pdf")
  defer fileTwo.Close()

  request := hellosign.EmbeddedRequest{
    TestMode: true,
    ClientId: os.Getenv("HS_CLIENT_ID"),
    File: []*os.File{
      fileOne,
      fileTwo,
    },
    Title:              "cool title",
    Subject:            "awesome",
    Message:            "cool message bro",
    SigningRedirectURL: "https://google.com",
    Signers: []hellosign.Signer{
      hellosign.Signer{
        Email: "freddy@hellosign.com",
        Name:  "Freddy Rangel",
        Pin:   "1234",
        Order: 1,
      },
      hellosign.Signer{
        Email: "frederick.rangel@gmail.com",
        Name:  "Frederick Rangel",
        Pin:   "1234",
        Order: 2,
      },
    },
    CCEmailAddresses: []string{
      "no@cats.com",
      "no@dogs.com",
    },
    UseTextTags:  false,
    HideTextTags: true,
    Metadata: map[string]string{
      "no":   "cats",
      "more": "dogs",
    },
    FormFieldsPerDocument: [][]hellosign.DocumentFormField{
      []hellosign.DocumentFormField{
        hellosign.DocumentFormField{
          APIId:    "api_id",
          Name:     "display name",
          Type:     "text",
          X:        123,
          Y:        456,
          Width:    678,
          Required: true,
          Signer:   0,
        },
      },
      []hellosign.DocumentFormField{
        hellosign.DocumentFormField{
          APIId:    "api_id_2",
          Name:     "display name 2",
          Type:     "text",
          X:        123,
          Y:        456,
          Width:    678,
          Required: true,
          Signer:   1,
        },
      },
    },
  }
  response, err := client.CreateEmbeddedSignatureRequest(request)
  if err != nil {
    log.Fatal(err)
  } else {
    fmt.Println(response)
  }

```


# License
```
The MIT License (MIT)

Copyright (C) 2015 hellosign.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
